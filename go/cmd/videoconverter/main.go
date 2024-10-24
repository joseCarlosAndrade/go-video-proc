package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	mergeChunks("mediatest/media/uploads/1", "merged.mp4" )
}

func extractNumber(fileName string) int {
	re := regexp.MustCompile(`\d+`) // finds numbers on the string
	numStr := re.FindString(filepath.Base(fileName)) // gets path/to/file2.chunk -> file2.chunk

	num , err := strconv.Atoi(numStr)

	if err != nil {
		return -1
	}
	
	return num
}

// search for all "*.chunk" alike files and merge them into the output file
func mergeChunks(inputDir, outputFile string) error {

	// search for all files *.chunk in this inputDir
	chunks, err := filepath.Glob(filepath.Join(inputDir, "*.chunk"))

	if err != nil {
		return fmt.Errorf("failed to find all chunks: %v", err)
	}

	// sorting the slice 
	sort.Slice(chunks, func(i int, j int) bool { // we have to explicit the way we want to order
		return extractNumber(chunks[i]) < extractNumber(chunks[j])
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