package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"

	imageshrink "github.com/imageshrink/imageshrink"
)

func localWorker(
	workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()
	for imagePath := range imagePaths {
		fmt.Printf("[Processing] %s\n", imagePath)
		convert, err := exec.LookPath("convert")
		if nil != err {
			panic("[Fatal] " + err.Error())
		}
		command := exec.Command(convert, imagePath, imagePath+".HEIF")
		err = command.Run()
		if nil != err {
			fmt.Printf(
				"[Error] Failed to process image: %s, error: %s\n",
				imagePath,
				err.Error(),
			)
			continue
		}
	}
}

func remoteWorker(
	workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()
	for imagePath := range imagePaths {
		fmt.Printf("[Processing] %s\n", imagePath)
		imageFile, err := os.OpenFile(imagePath, os.O_RDONLY, 0)
		if nil != err {
			fmt.Printf(
				"[Error] Failed to read image: %s, error: %s\n",
				imagePath,
				err.Error(),
			)
			continue
		}
		defer imageFile.Close()
		resp, err := http.Post("http://localhost:8080/convert", "image/jpeg", imageFile)
		if nil != err {
			fmt.Printf(
				"[Error] Failed to transfer image: %s, error: %s\n",
				imagePath,
				err.Error(),
			)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[Error] Failed to process image: %v\n", resp.Status)
			continue
		}
		buffer := make([]byte, 8192)
		imageFileNew, err := os.OpenFile(imagePath+".HEIF", os.O_CREATE|os.O_WRONLY, 0644)
		if nil != err {
			fmt.Printf(
				"[Error] Failed to create image: %s, error: %s\n",
				imagePath+".HEIF",
				err,
			)
			continue
		}
		io.CopyBuffer(imageFileNew, resp.Body, buffer)
		imageFileNew.Close()
		resp.Body.Close()
	}
}
func buildWorker() imageshrink.Worker {
	return func(workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup) {
		remoteWorker(workerID, imagePaths, waitGroup)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("Usage: imageshrink [path to scan]\n")
		return
	}
	imageshrink.ImageShrink(args[0], buildWorker())
	return
}
