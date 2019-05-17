package main

import (
  "flag"
  "fmt"
  "github.com/imageshrink/imageshrink/go"
  "os/exec"
  "strconv"
  "sync"
)

func workerImpl(workerID int, imagePaths <-chan string, size int, quality int, waitGroup *sync.WaitGroup) {
  defer waitGroup.Done()
  sizeString := strconv.Itoa(size)
  qualityString := strconv.Itoa(quality)
  for imagePath := range imagePaths {
    fmt.Printf("[Processing] %s\n", imagePath)
    convert, err := exec.LookPath("convert")
    if nil != err {
      panic("[Fatal] " + err.Error())
      return
    }
    command := exec.Command(
      convert,
      "-resize" , sizeString +"x" + sizeString,
      "-quality", qualityString,
      "-interlace", "JPEG",
      imagePath, imagePath)
    err = command.Run()
    if nil != err {
      fmt.Printf("[Error] Failed to process image: %s, error: %s\n", imagePath, err.Error())
      continue
    }
  }
}

func buildWorker(size int, quality int) imageshrink.Worker {
  return func (workerID int, imagePaths <-chan string, waitGroup *sync.WaitGroup) {
    workerImpl(workerID, imagePaths, size, quality, waitGroup)
  }
}

func main()  {
  size := flag.Int("size", 4096, "image size")
  quality := flag.Int("quality", 90, "image quality (1~100)")
  flag.Parse()
  args := flag.Args()
  if len(args) != 1 {
    fmt.Printf("Usage: imageshrink [path to scan]\n")
    return
  }
  imageshrink.ImageShrink(args[0], buildWorker(*size, *quality))
  return
}
