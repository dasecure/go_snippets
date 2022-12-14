package main

import (
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

	// Modify the XMP data.
	xmpData.SetProperty("http://ns.adobe.com/xap/1.0/", "Creator", "John Doe")

	// Write the modified XMP data to the raw image file.
	err = xmp.Write(f, xmpData)
	if err != nil {
		log.Fatal(err)
	}
}
