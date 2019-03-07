package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	_ = filepath.Walk("/home/huahang/Desktop", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
}
