package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dsoprea/go-xmp"
)

func main() {
	// Open the raw image file.
	f, err := os.Open("image.raw")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read the XMP data from the raw image file.
	xmpData, err := xmp.Read(f)
	if err != nil {
		log.Fatal(err)
	}

	// Print the XMP data.
	fmt.Println(xmpData)
}
