package main

import (
	"flag"
	"fmt"
	"os/exec"
	"sync"

	imageshrink "github.com/imageshrink/imageshrink/go"
)

func workerImpl(workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup) {
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
			fmt.Printf("[Error] Failed to process image: %s, error: %s\n", imagePath, err.Error())
			continue
		}
	}
}

func buildWorker() imageshrink.Worker {
	return func(workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup) {
		workerImpl(workerID, imagePaths, waitGroup)
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
