package main

import (
	"image"
	"image/color"
	"math"
)

// compareImages returns the percentage difference between two images
func compareImages(img1, img2 image.Image) float64 {
	// get dimensions of both images
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	// make sure the images have the same dimensions
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return 0
	}

	// calculate the total number of pixels in the image
	totalPixels := bounds1.Dx() * bounds1.Dy()

	// initialize counters for different pixel values
	var differentPixels, similarPixels int

	// iterate over all pixels in the image
	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			// get the color of the current pixel in both images
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)

			// check if the colors are the same
			if !colorsEqual(c1, c2) {
				differentPixels++
			} else {
				similarPixels++
			}
		}
	}

	// calculate the percentage difference
	return math.Round(float64(differentPixels) / float64(
