package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"

	"github.com/deluan/lookup"
)

// Call this from main
func CaptureAndLookup() {
	captureWindowToFile("all.png", "Untitled - Notepad")
	runLookup("all.png", "ana_blue.png")
}

func captureWindowToFile(outputFile, windowTitle string) {
	hwnd := findWindow(windowTitle)
	if hwnd == 0 {
		panic("Window not found: " + windowTitle)
	}

	var rect win.RECT
	win.GetClientRect(hwnd, &rect)

	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top

	hwndDC := win.GetDC(hwnd)
	memDC := win.CreateCompatibleDC(hwndDC)
	bitmap := win.CreateCompatibleBitmap(hwndDC, width, height)
	old := win.SelectObject(memDC, win.HGDIOBJ(bitmap))

	win.BitBlt(memDC, 0, 0, width, height, hwndDC, 0, 0, win.SRCCOPY)
	win.SelectObject(memDC, old)

	saveBitmapAsPNG(outputFile, bitmap, memDC, width, height)

	win.ReleaseDC(hwnd, hwndDC)
	win.DeleteDC(memDC)
	win.DeleteObject(win.HGDIOBJ(bitmap))
}

func findWindow(title string) win.HWND {
	var target win.HWND
	win.EnumWindows(syscall.NewCallback(func(hwnd win.HWND, lparam uintptr) uintptr {
		buf := make([]uint16, 200)
		win.GetWindowText(hwnd, &buf[0], int32(len(buf)))
		windowText := syscall.UTF16ToString(buf)
		if strings.Contains(windowText, title) {
			target = hwnd
			return 0
		}
		return 1
	}), 0)
	return target
}

func saveBitmapAsPNG(filename string, hBitmap win.HBITMAP, hdc win.HDC, width, height int32) {
	var bmp win.BITMAP
	win.GetObject(win.HGDIOBJ(hBitmap), unsafe.Sizeof(bmp), unsafe.Pointer(&bmp))

	bmpInfo := win.BITMAPINFO{
		BmiHeader: win.BITMAPINFOHEADER{
			BiSize:        uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})),
			BiWidth:       width,
			BiHeight:      -height, // top-down
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: win.BI_RGB,
		},
	}

	bufSize := int(width * height * 4)
	buf := make([]byte, bufSize)
	win.GetDIBits(hdc, hBitmap, 0, uint32(height), unsafe.Pointer(&buf[0]), &bmpInfo, win.DIB_RGB_COLORS)

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	copy(img.Pix, buf)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func runLookup(fullImagePath, templateImagePath string) {
	img := loadImage(fullImagePath)
	l := lookup.NewLookup(img)
	template := loadImage(templateImagePath)

	pp, _ := l.FindAll(template, 0.9)

	if len(pp) > 0 {
		fmt.Printf("Found %d matches:\n", len(pp))
		for _, p := range pp {
			fmt.Printf("- (%d, %d) with %f accuracy\n", p.X, p.Y, p.G)
		}
	} else {
		fmt.Println("No matches found")
	}
}

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
