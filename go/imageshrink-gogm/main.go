package main

import (
	"fmt"
	"github.com/gographics/gmagick"
	"github.com/imageshrink/imageshrink/go"
	"os"
	"sync"
)

func gmWorker(workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	var err error
	for imagePath := range imagePaths {
		fmt.Printf("[Processing] %s\n", imagePath)
		wand := gmagick.NewMagickWand()
		err = wand.SetResourceLimit(gmagick.RESOURCE_MEMORY, 128*1024*1024)
		if err != nil {
			fmt.Printf("[Error] Failed to set resource, error: %s\n", err.Error())
			return
		}
		err = wand.ReadImage(imagePath)
		if err != nil {
			fmt.Printf("[Error] Failed to read image: %s, error: %s\n", imagePath, err.Error())
			continue
		}
		width := float64(wand.GetImageWidth())
		height := float64(wand.GetImageHeight())
		var scale = 1.0
		if width > height {
			scale = float64(4096) / width
			scale = float64(4096) / height
		}
		newWidth := uint(width * scale)
		newHeight := uint(height * scale)
		err = wand.ResizeImage(newWidth, newHeight, gmagick.FILTER_LANCZOS, 1)
		if err != nil {
			fmt.Printf("[Error] Failed to resize image: %s, error: %s\n", imagePath, err.Error())
			continue
		}
		err = wand.SetCompressionQuality(90)
		if err != nil {
			fmt.Printf("[Error] Failed to set quality: %s, error: %s\n", imagePath, err.Error())
			continue
		}
		err = wand.WriteImage(imagePath)
		if err != nil {
			fmt.Printf("[Error] Failed to save file: %s, error: %s\n", imagePath, err.Error())
			continue
		}
	}
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: imageshrink-gogm [path to scan]\n")
		return
	}
	gmagick.Initialize()
	defer gmagick.Terminate()
	imageshrink.ImageShrink(args[1], gmWorker)
	return
}
