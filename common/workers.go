package common

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// Worker Worker
type Worker func(imagePaths <-chan string, waitGroup *sync.WaitGroup)

// MakeLocalWorker MakeLocalWorker
func MakeLocalWorker() Worker {
	return localWorker
}

func localWorker(paths <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for imagePath := range paths {
		fmt.Printf("[Processing] %s\n", imagePath)
		convert, err := exec.LookPath("convert")
		if nil != err {
			panic("[Fatal] " + err.Error())
		}
		command := exec.Command(
			convert, "-quality", "50", imagePath, imagePath+".HEIF",
		)
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

// MakeRemoteWorker MakeRemoteWorker
func MakeRemoteWorker(host string) Worker {
	return func(paths <-chan string, wg *sync.WaitGroup) {
		remoteWorker(host, paths, wg)
	}
}

func remoteWorker(
	host string, paths <-chan string, wg *sync.WaitGroup,
) {
	defer wg.Done()
	for imagePath := range paths {
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
		digest, err := ComputeMD5(imageFile)
		digestHex := strings.ToUpper(fmt.Sprintf("%x", digest))
		imageFile.Seek(0, 0)
		url := fmt.Sprintf("http://%v/convert", host)
		request, err := http.NewRequest(http.MethodPost, url, imageFile)
		request.Header.Add("Content-Type", "image/jpeg")
		request.Header.Add("Content-MD5", digestHex)
		resp, err := http.DefaultClient.Do(request)
		if nil != err {
			fmt.Printf(
				"[Error] Failed to transfer image: %s, error: %s\n",
				imagePath,
				err.Error(),
			)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[Error] Failed to process image: %v\n", resp.Status)
			continue
		}
		imageFileNew, err := os.OpenFile(
			imagePath+".HEIF", os.O_CREATE|os.O_WRONLY, 0644,
		)
		if nil != err {
			fmt.Printf(
				"[Error] Failed to create image: %s, error: %s\n",
				imagePath+".HEIF",
				err,
			)
			continue
		}
		defer imageFileNew.Close()
		digest, _, err = CopyAndComputeMD5(imageFileNew, resp.Body)
		if err != nil {
			fmt.Printf("[Error] Unable to copy image: %v\n", imagePath)
			continue
		}
		digestHex = strings.ToUpper(fmt.Sprintf("%x", digest))
		md5Header := resp.Header.Get("Content-MD5")
		if len(md5Header) == 0 {
			fmt.Printf("[Error] Missing Content-MD5 for image: %v\n", imagePath)
			continue
		}
		if md5Header != digestHex {
			fmt.Printf("[Error] Content-MD5 not matched for image: %v\n", imagePath)
			continue
		}
	}
}

// DoImageShrink DoImageShrink
func DoImageShrink(scanPath string, workers []Worker) {
	var waitGroup sync.WaitGroup
	imagePaths := make(chan string, 128)
	numWorkers := len(workers)
	waitGroup.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go workers[i](imagePaths, &waitGroup)
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
