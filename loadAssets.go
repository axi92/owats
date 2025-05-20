package main

import (
	"image"
	"math"
	"os"
	"path/filepath"

	"github.com/deluan/lookup"
)

// Global map to store loaded images
var Assets = map[string]struct {
	Image        image.Image
	SearchBounds image.Rectangle
	Lookup       *lookup.Lookup
}{}

// LoadImages loads all images from the assets directory into the global map
func LoadImages() {
	assetsDir := "assets/"
	// image host to upload testing images to use in codepen example
	// https://imgbox.com/
	// Get percentage of parts of an full image:
	// https://codepen.io/mattiacalicchia/full/zMOyrY
	// Inputs are:
	// width: 1920
	// height: 1080
	// image url: https://images2.imgbox.com/64/29/jdJjtVzg_o.jpg
	files := []struct {
		Path         string
		SearchBounds image.Rectangle
	}{
		{
			Path: "scoreboard/ana_big.jpg",
			SearchBounds: image.Rect(
				int(math.Round(0.64*1920)), // 64% of width
				int(math.Round(0.14*1080)), // 14% of height
				int(math.Round(0.74*1920)), // 74% of width
				int(math.Round(0.31*1080)), // 31% of height
			),
		},
		{
			Path: "scoreboard/ana_small.jpg",
			SearchBounds: image.Rect(
				// x points from left to right
				// y points from top down
				int(math.Round(0.14*1920)), // x min
				int(math.Round(0.17*1080)), // y min
				int(math.Round(0.18*1920)), // x max
				int(math.Round(0.47*1080)), // y max
			),
		},
	}

	for _, file := range files {
		imgPath := filepath.Join(assetsDir, file.Path)
		img := loadImage(imgPath)
		Assets[file.Path] = struct {
			Image        image.Image
			SearchBounds image.Rectangle
			Lookup       *lookup.Lookup
		}{
			Image:        img,
			SearchBounds: file.SearchBounds,
		}
	}
}

// Helper function to load an image from the filesystem
func loadImage(imgPath string) image.Image {
	imageFile, err := os.Open(imgPath)
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		panic(err)
	}
	return img
}
