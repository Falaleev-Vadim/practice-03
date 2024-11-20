package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <directory> <watermark>")
		return
	}

	directory := os.Args[1]
	watermarkText := os.Args[2]

	var wg sync.WaitGroup

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isImageFile(path) {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				processImage(filePath, watermarkText)
			}(path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
		return
	}

	wg.Wait()
	fmt.Println("Processing complete.")
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}

func processImage(imagePath, watermarkText string) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening image:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	imgWithWatermark := applyWatermark(img, watermarkText)

	outputPath := filepath.Join("output", filepath.Base(imagePath))
	err = saveImage(outputPath, imgWithWatermark)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Processed:", imagePath)
}

func applyWatermark(img image.Image, watermarkText string) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(dst, bounds, img, image.Point{0, 0}, draw.Over)

	textColor := color.RGBA{255, 255, 255, 255}

	addWatermarkText(dst, watermarkText, width-150, height-50, textColor)

	return dst
}

func addWatermarkText(img *image.RGBA, text string, x, y int, c color.Color) {
}

func saveImage(filePath string, img image.Image) error {     
	ext := strings.ToLower(filepath.Ext(filePath))
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if ext == ".png" {
		return png.Encode(outFile, img)
	} else if ext == ".jpg" || ext == ".jpeg" {
		return jpeg.Encode(outFile, img, nil)
	}
	return fmt.Errorf("unsupported file format")
}
