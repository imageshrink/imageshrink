package imageshrink

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Worker Worker
type Worker func(workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup)

// ImageShrink ImageShrink
func ImageShrink(scanPath string, worker Worker) {
	var waitGroup sync.WaitGroup
	imagePaths := make(chan string, 128)
	numCPU := runtime.NumCPU()
	waitGroup.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go worker(i, imagePaths, &waitGroup)
	}
	_ = filepath.Walk(scanPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("[Error] Hit an error! " + err.Error() + "\n")
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
}
