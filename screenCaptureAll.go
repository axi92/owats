package main

import (
	"fmt"
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
