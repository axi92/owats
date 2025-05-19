package main

import (
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"os"

	"github.com/deluan/lookup"
	"github.com/kbinani/screenshot"
)

var displayNumber int = 2

func main() {
	Example_lookup()
}

// Helper function to load an image from the filesystem
func loadImage(imgPath string) image.Image {
	imageFile, _ := os.Open(imgPath)
	defer imageFile.Close()
	img, _, _ := image.Decode(imageFile)
	return img
}

func ScreenCaptureAll() {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)

		fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}
}

func Example_lookup() {
	// Load full image
	img := loadImage("all.png")

	// Create a lookup for that image
	l := lookup.NewLookup(img)

	// Load a template to search inside the image
	template := loadImage("ana_blue.png")

	// Find all occurrences of the template in the image
	pp, _ := l.FindAll(template, 0.9)

	// Print the results
	if len(pp) > 0 {
		fmt.Printf("Found %d matches:\n", len(pp))
		for _, p := range pp {
			fmt.Printf("- (%d, %d) with %f accuracy\n", p.X, p.Y, p.G)
		}
	} else {
		println("No matches found")
	}

	// Output:
	// Found 1 matches:
	// - (21, 7) with 0.997942 accuracy
}
