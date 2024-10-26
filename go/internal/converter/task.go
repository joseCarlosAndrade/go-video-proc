package converter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type VideoConverter struct {
}

func NewVideoConverter()*VideoConverter {
	return &VideoConverter{}
}

// json that we'll receive
type VideoTask struct {
	VideoId int    `json:"video_id"`
	Path    string `json:"path"`
}

func (vc *VideoConverter) Handle(msg []byte) {
	var task VideoTask

	if err := json.Unmarshal(msg, &task); err != nil {
		vc.logError(task, "failed to unmarshal task", err)
		return
	}

	err := vc.processVideo(&task)
	if err != nil {
		vc.logError(task, "failed to process video", err)
		return 
	}
	
}

func (vc *VideoConverter) processVideo(task *VideoTask) error {
	mergedFile := filepath.Join(task.Path, "merged.mp4")
	mpegDashPath := filepath.Join(task.Path, "mpeg-dash")

	slog.Info("merging chunks", slog.String("path", task.Path))

	err := vc.mergeChunks(task.Path, mergedFile) // will generate a .mp4 on the specified path
	if err != nil {
		vc.logError(*task, "failed to merge chunk", err)
		return err
	}

	slog.Info("Creating mpeg-dash directory", slog.String("path", task.Path))
	err = os.MkdirAll(mpegDashPath, os.ModePerm)
	if err != nil {
		vc.logError(*task, "failed to create mpeg-dash directory", err)
		return err
	}

	slog.Info("Converting to mpeg-dash", slog.String("path", task.Path))
	ffmpegCmd := exec.Command(
		"ffmpeg", "-i", mergedFile,
		"-f", "dash",
		filepath.Join(mpegDashPath, "output.mpd"),
	)

	output, err := ffmpegCmd.CombinedOutput()

	if err != nil {
		vc.logError(*task, "failed to convert video to mpeg-dash, outpuy: " + string(output), err)
		return err
	}
	
	slog.Info("video converted to mpeg-dash", slog.String("path", mpegDashPath))

	// removingo the previous merged mp4 file
	err = os.Remove(mergedFile)

	if err != nil {
		vc.logError(*task, "failed to remove merged file", err)
		return err
	}
	
	return nil
}

// creating a structured default log to facilitate log handling
func (vc *VideoConverter) logError(task VideoTask, message string, err error) {
	errorData := map[string]any{
		"video_id": task.VideoId,
		"error":    message,
		"details":  err.Error(),
		"time":     time.Now(),
	} // any, interface{} = qualquer tipo = any

	serializedError, _ := json.Marshal(errorData)
	slog.Error("processing error", slog.String("error_details", string(serializedError)))

	// todo: save error on database
}

func (vc *VideoConverter) extractNumber(fileName string) int {
	re := regexp.MustCompile(`\d+`)                  // finds numbers on the string
	numStr := re.FindString(filepath.Base(fileName)) // gets path/to/file2.chunk -> file2.chunk

	num, err := strconv.Atoi(numStr)

	if err != nil {
		return -1
	}

	return num
}

// search for all "*.chunk" alike files and merge them into the output file
func (vc *VideoConverter) mergeChunks(inputDir, outputFile string) error {

	// search for all files *.chunk in this inputDir
	chunks, err := filepath.Glob(filepath.Join(inputDir, "*.chunk"))

	if err != nil {
		return fmt.Errorf("failed to find all chunks: %v", err)
	}

	// sorting the slice
	sort.Slice(chunks, func(i int, j int) bool { // we have to explicit the way we want to order
		return vc.extractNumber(chunks[i]) < vc.extractNumber(chunks[j])
	})

	// creating the output file
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("output file not created: %v", err)
	}
	defer output.Close()

	for _, chunk := range chunks {
		input, err := os.Open(chunk)
		if err != nil {
			return fmt.Errorf("could not open chunk: %v", err)
		}

		_, err = output.ReadFrom(input) // reading chunk into output file, to merge

		if err != nil {
			return fmt.Errorf("could not write chunk %s to merged file: %v", chunk, err)
		}
		input.Close()
	}

	return nil
}
