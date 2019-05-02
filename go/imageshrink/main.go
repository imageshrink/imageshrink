package main

import (
	"fmt"
	"github.com/imageshrink/imageshrink/go"
	"os"
	"os/exec"
	"sync"
)

func convertWorker(workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for imagePath := range imagePaths {
		fmt.Printf("[Processing] %s\n", imagePath)
		convert, err := exec.LookPath("convert")
		if nil != err {
			panic("[Fatal] " + err.Error())
			return
		}
		command := exec.Command(
			convert,
			"-resize" , "4096x4096", "-quality", "90", imagePath, imagePath)
		err = command.Run()
		if nil != err {
			fmt.Printf("[Error] Failed to process image: %s, error: %s\n", imagePath, err.Error())
			continue
		}
	}
}

func main()  {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: imageshrink [path to scan]\n")
		return
	}
	imageshrink.ImageShrink(args[1], convertWorker)
	return
}
