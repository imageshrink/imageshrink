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
	var err error
	for imagePath := range imagePaths {
		fmt.Printf("[Processing] %s\n", imagePath)
		convert, _ := exec.LookPath("convert")
		command := exec.Command(
			convert,
			"-resize" , "4096x4096", "-quality", "90", imagePath, imagePath)
		err = command.Run()
		if nil != err {
			out, _ := command.CombinedOutput()
			fmt.Printf(string(out))
			fmt.Printf("[Error] Failed to process image: %s, error: %s\n", imagePath, err.Error())
			continue
		}
	}
}

func main()  {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: imageshrink-go [path to scan]\n")
		return
	}
	imageshrink.ImageShrink(args[1], convertWorker)
	return
}
