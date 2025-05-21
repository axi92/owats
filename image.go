package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/kbinani/screenshot"
)

func ScreenCaptureAll() {
	n := screenshot.NumActiveDisplays()

	for i := range n {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%d_%dx%d.jpg", i, bounds.Dx(), bounds.Dy())
		file, _ := os.Create(fileName)
		defer file.Close()
		jpeg.Encode(file, img, nil)

		fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}
}

func SaveImageToFile(image image.Image, filename string) {
	// Save the subimage to a file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	err = jpeg.Encode(file, image, nil)
	if err != nil {
		fmt.Printf("Failed to encode subimage to file %s: %v\n", filename, err)
		return
	}
	fmt.Printf("Subimage saved to %s\n", filename)
	// END SAVE
}

// Function to draw a rectangle on an image
func drawRectangle(img image.Image, rect image.Rectangle, col color.Color) image.Image {
	// Create a new RGBA image with the same bounds as the input image
	rgbaImg := image.NewRGBA(img.Bounds())

	// Copy the input image to the RGBA image
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgbaImg.Set(x, y, img.At(x, y))
		}
	}

	// Draw the top and bottom edges
	for x := rect.Min.X; x < rect.Max.X; x++ {
		rgbaImg.Set(x, rect.Min.Y, col)   // Top edge
		rgbaImg.Set(x, rect.Max.Y-1, col) // Bottom edge
	}

	// Draw the left and right edges
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		rgbaImg.Set(rect.Min.X, y, col)   // Left edge
		rgbaImg.Set(rect.Max.X-1, y, col) // Right edge
	}

	return rgbaImg.SubImage(img.Bounds())
}

func isImageChanged(img1, img2 *image.RGBA) bool {
	// Ensure the bounds are the same
	if !img1.Bounds().Eq(img2.Bounds()) {
		return true
	}

	// Compare pixel data
	for y := img1.Bounds().Min.Y; y < img1.Bounds().Max.Y; y++ {
		for x := img1.Bounds().Min.X; x < img1.Bounds().Max.X; x++ {
			if img1.RGBAAt(x, y) != img2.RGBAAt(x, y) {
				return true
			}
		}
	}
	return false
}
