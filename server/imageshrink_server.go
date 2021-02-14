package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func main() {
	port := flag.Int("port", 8080, "server port")
	dir := flag.String("dir", "/tmp/imageshrink", "working dir")
	flag.Parse()
	err := os.MkdirAll(*dir, 0755)
	if err != nil {
		log.Fatalf("Unable to create working dir: %v\n", err)
	}
	http.HandleFunc("/ruok", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "imok")
	})
	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		buffer := make([]byte, 8192)
		fileNameOld := *dir + string(os.PathSeparator) + uuid.NewString() + ".jpg"
		fileNameNew := fileNameOld + ".output.heif"
		imageFileOld, err := os.OpenFile(fileNameOld, os.O_CREATE|os.O_WRONLY, 0644)
		if nil != err {
			log.Printf("Unable to open file: %v, err: %v\n", fileNameOld, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		length, err := io.CopyBuffer(imageFileOld, r.Body, buffer)
		if nil != err {
			log.Printf("Unable to read request err: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Bytes read: %v", length)
		r.Body.Close()
		imageFileOld.Close()
		convert, err := exec.LookPath("convert")
		if nil != err {
			panic("[Fatal] " + err.Error())
		}
		command := exec.Command(convert, "-quality", "50", fileNameOld, fileNameNew)
		err = command.Run()
		if nil != err {
			fmt.Printf(
				"[Error] Failed to process image: %s, error: %s\n",
				fileNameOld,
				err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		imageFileNew, err := os.OpenFile(fileNameNew, os.O_RDONLY, 0644)
		if nil != err {
			log.Printf("Unable to open file: %v, err: %v\n", fileNameNew, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		header := w.Header()
		header.Add("Content-Type", "image/heif")
		w.WriteHeader(http.StatusOK)
		io.CopyBuffer(w, imageFileNew, buffer)
		imageFileNew.Close()
		os.Remove(fileNameOld)
		os.Remove(fileNameNew)
	})
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		log.Fatalf("Unable to start http server: %v\n", err)
	}
	return
}
