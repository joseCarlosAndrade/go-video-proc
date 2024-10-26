package main

import (
	// "fmt"
	// "os"
	// "path/filepath"
	// "regexp"
	// "sort"
	// "strconv"
	// "github.com/joseCarlosAndrade/videoconverter/converter"
	"videoproc/internal/converter"
)

func main() {
	// mergeChunks("mediatest/media/uploads/1", "merged.mp4" )
	// vc := NewVideo
	vc := converter.NewVideoConverter()
	vc.Handle([]byte(`{ "video_id" : 1 , "path" : "/media/uploads/1" }`))
}

