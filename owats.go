package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"time"

	"github.com/deluan/lookup"
	"github.com/kbinani/screenshot"
)

var displayNumber int = 2

func main() {
	// Load assets before starting the loop
	LoadImages()

	// Create a ticker that triggers every 500ms
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	// Run the loop until the program is terminated
	for range ticker.C {
		Example_lookup()
	}
}

func ScreenCapture(displayNumber int) *image.RGBA {
	bounds := screenshot.GetDisplayBounds(displayNumber)
	fmt.Printf("Display Bounds: %v\n", bounds)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	return img
}

func Example_lookup() {
	start := time.Now()
	screenCapture := ScreenCapture(displayNumber)
	elapsed := time.Since(start)
	fmt.Printf("Execution time for Screenshot: %s\n", elapsed)

	// Extract smaller images for each asset
	for key, assetImg := range Assets {
		// Get the SearchBounds from the asset image
		rect := assetImg.SearchBounds

		// Create a SubImage using the rectangle from SearchBounds
		subImg := screenCapture.SubImage(rect).(*image.RGBA)

		// SaveImageToFile(subImg, "subimage.jpeg")

		// Create a lookup for that image
		imageLookup := lookup.NewLookup(subImg)

		// Save this image back into the assets as TempImage
		assetImg.Lookup = imageLookup
		assetImg.Screenshot = subImg
		Assets[key] = assetImg
	}

	// Loop over the Assets map
	for name, img := range Assets {
		fmt.Printf("Image Name: %s\n", name)

		// Find all occurrences of the template in the image
		start := time.Now()
		pp, _ := img.Lookup.FindAll(img.Image, 0.98) // Dont go to low, a match was during developing at least 0.99
		elapsed := time.Since(start)
		fmt.Printf("Execution time for Lookup: %s\n", elapsed)
		// Print the results
		if len(pp) > 0 {
			fmt.Printf("Found %d matches:\n", len(pp))
			for _, p := range pp {
				fmt.Printf("- (%d, %d) with %f accuracy\n", p.X, p.Y, p.G)
				drawingRect := image.Rect(p.X, p.Y, p.X+1, p.Y+1)

				drawing := drawRectangle(img.Screenshot, drawingRect, color.RGBA{255, 255, 255, 255})
				// isChanged := isImageChanged(originalImage, painted)
				// fmt.Printf("Image has changed: %v\n", isChanged)
				SaveImageToFile(drawing, "drawing.jpeg")
			}
		} else {
			println("No matches found")
		}
		// Output:
		// Found 1 matches:
		// - (21, 7) with 0.997942 accuracy
	}
}
