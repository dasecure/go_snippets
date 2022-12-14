package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	// Open the raw image file.
	f, err := os.Open("image.raw")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read the EXIF data from the raw image file.
	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Print the EXIF data.
	fmt.Println(x)
}
