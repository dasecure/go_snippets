package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/unixpickle/essentials"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: time_lapse_video <image_folder> <output_file>")
		os.Exit(1)
	}

	// Open the output file.
	outputFile, err := os.Create(os.Args[2])
	if err != nil {
		essentials.Die(err)
	}
	defer outputFile.Close()

	// Encode the video with ffmpeg.
	videoEncoder := ffmpeg.NewVideoEncoder(outputFile)
	videoEncoder.FPS = 30
	videoEncoder.Width = 640
	videoEncoder.Height = 480

	// Iterate over the images in the folder and add them to the video.
	if err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		img, err := jpeg.Decode(file)
		if err != nil {
			return err
		}
		return videoEncoder.AddFrame(img)
	}); err != nil {
		log.Fatal(err)
	}

	// Finish the video and check for errors.
	if err := videoEncoder.Close(); err != nil {
		log.Fatal(err)
	}
}
