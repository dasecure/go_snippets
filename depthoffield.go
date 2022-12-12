package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	"github.com/unixpickle/essentials"
)

const (
	Threshold = 0.1
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: dof_analysis <image_file>")
	}

	// Load the image and convert it to grayscale.
	file, err := os.Open(os.Args[1])
	if err != nil {
		essentials.Die(err)
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		essentials.Die(err)
	}
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y,
				color.GrayModel.Convert(img.At(x, y)))
		}
	}

	// Apply a Gaussian blur to the grayscale image.
	blurredImg := gaussianBlur(grayImg)

	// Use the Sobel operator to detect edges in the blurred image.
	edgeImg := sobelEdges(blurredImg)

	// Calculate the DOF score for each pixel and classify it as in
	//focus or out of focus.
	dofImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			score := dofScore(edgeImg, x, y)
			if score > Threshold {
				dofImg.Set(x, y, color.Gray{Y: 255})
			}
		}
	}

	// Save the DOF image to a file.
	outFile, err := os.Create("dof_result.jpg")
	if err != nil {
		essentials.Die(err)
	}
	defer outFile.Close()
	if err := jpeg.Encode(outFile, dofImg, &jpeg.Options{Quality: 100}); err != nil {
		essentials.Die(err)
	}
}

// gaussianBlur applies a Gaussian blur to the input image.
// func gaussianBlur(img *image.Gray) *image.Gray {
// 	// TODO: implement Gaussian blur
// }

// gaussianBlur applies a Gaussian blur to the input image.
func gaussianBlur(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	blurredImg := image.NewGray(bounds)

	// Create a Gaussian kernel with sigma=1 and size=5.
	// You can adjust these values to change the strength of the blur.
	var kernel [5][5]float64
	sigma := 1.0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			x := float64(i - 2)
			y := float64(j - 2)
			kernel[i][j] = math.Exp(-(x*x + y*y) / (2 * sigma * sigma))
		}
	}

	// Convolve the image with the kernel.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sum float64
			for i := 0; i < 5; i++ {
				for j := 0; j < 5; j++ {
					if x+i-2 < bounds.Min.X || x+i-2 >= bounds.Max.X || y+j-2 < bounds.Min.Y || y+j-2 >= bounds.Max.Y {
						continue
					}
					p := img.GrayAt(x+i-2, y+j-2).Y
					sum += float64(p) * kernel[i][j]
				}
			}
			blurredImg.Set(x, y, color.Gray{Y: uint8(sum)})
		}
	}

	return blurredImg
}

// sobelEdges applies the Sobel operator to the input image and returns
// the resulting edge image.
//
//	func sobelEdges(img *image.Gray) *image.Gray {
//		// TODO: implement Sobel operator
//	}
//
// sobelEdges applies the Sobel operator to the input image and returns the resulting edge image.
func sobelEdges(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	edgeImg := image.NewGray(bounds)

	// Create the Sobel kernels.
	var xKernel [3][3]float64
	xKernel[0][0] = -1
	xKernel[0][1] = 0
	xKernel[0][2] = 1
	xKernel[1][0] = -2
	xKernel[1][1] = 0
	xKernel[1][2] = 2
	xKernel[2][0] = -1
	xKernel[2][1] = 0
	xKernel[2][2] = 1
	var yKernel [3][3]float64
	yKernel[0][0] = -1
	yKernel[0][1] = -2
	yKernel[0][2] = -1
	yKernel[1][0] = 0
	yKernel[1][1] = 0
	yKernel[1][2] = 0
	yKernel[2][0] = 1
	yKernel[2][1] = 2
	yKernel[2][2] = 1

	// Convolve the image with the x and y kernels to compute the gradient in each direction.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var xGradient float64
			var yGradient float64
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if x+i-1 < bounds.Min.X || x+i-1 >= bounds.Max.X || y+j-1 < bounds.Min.Y || y+j-1 >= bounds.Max.Y {
						continue
					}
					p := img.GrayAt(x+i-1, y+j-1).Y
					xGradient += float64(p) * xKernel[i][j]
					yGradient += float64(p) * yKernel[i][j]
				}
			}
			// Set the pixel intensity in the edge image based on the gradient magnitude.
			// You can adjust the scaling factor to change the sensitivity of the edge detection.
			mag := math.Sqrt(xGradient*xGradient + yGradient*yGradient)
			edgeImg.Set(x, y, color.Gray{Y: uint8(mag / 100)})
		}
	}

	return edgeImg
}

// dofScore calculates the DOF score for the pixel at (x, y) in the input image.
func dofScore(img *image.Gray, x, y int) float64 {
	bounds := img.Bounds()
	var centerContrast float64
	var surroundingContrast float64
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if x+i < bounds.Min.X || x+i >= bounds.Max.X || y+j < bounds.Min.Y || y+j >= bounds.Max.Y {
				continue
			}
			p1 := img.GrayAt(x, y).Y
			p2 := img.GrayAt(x+i, y+j).Y
			if i == 0 && j == 0 {
				centerContrast += math.Abs(float64(p1 - p2))
			} else {
				surroundingContrast += math.Abs(float64(p1 - p2))
			}
		}
	}
	return centerContrast / surroundingContrast
}
