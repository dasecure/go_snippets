package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	// Create a slice of JPG files to combine
	jpgFiles := []string{"image1.jpg", "image2.jpg", "image3.jpg"}

	// Use the "ffmpeg" command to combine the JPG files into an MP4 video
	cmd := exec.Command("ffmpeg", "-i", "concat:"+strings.Join(jpgFiles, "|"), "-c", "copy", "output.mp4")

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("MP4 video created successfully")
}
