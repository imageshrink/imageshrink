package main

import (
	"fmt"
	"github.com/gographics/gmagick"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)


func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: imageshrink-go [path to scan]\n")
		return
	}
	gmagick.Initialize()
	defer gmagick.Terminate()
	var waitGroup sync.WaitGroup
	imagePaths := make(chan string, 128)
	worker := func(workerID int) {
		defer waitGroup.Done()
		var err error
		for imagePath := range imagePaths {
			fmt.Printf("[Processing] %s\n", imagePath)
			wand := gmagick.NewMagickWand()
			err = wand.SetResourceLimit(gmagick.RESOURCE_MEMORY, 512*1024*1024)
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
	numCPU := runtime.NumCPU()
	waitGroup.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go worker(i)
	}
	_ = filepath.Walk(args[1], func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Hit an error! " + err.Error() + "\n")
			return err
		}
		if !fileInfo.Mode().IsRegular() {
			return nil
		}
		ext := path.Ext(filePath)
		if !strings.EqualFold(ext, ".jpeg") && !strings.EqualFold(ext, ".jpg") {
			return nil
		}
		imagePaths <- filePath
		return nil
	})
	close(imagePaths)
	waitGroup.Wait()
	return
}
