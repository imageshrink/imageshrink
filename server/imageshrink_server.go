package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/imageshrink/imageshrink/common"
)

func main() {
	port := flag.Int("port", 58080, "server port")
	dir := flag.String("dir", "/tmp/imageshrink", "working dir")
	flag.Parse()
	err := os.MkdirAll(*dir, 0755)
	if err != nil {
		panic(fmt.Sprintf("[Fatal] Unable to create working dir: %v\n", err))
	}
	http.HandleFunc("/ruok", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "imok")
	})
	http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fileNameOld := *dir + string(os.PathSeparator) + uuid.NewString() + ".jpg"
		fileNameNew := fileNameOld + ".output.heif"
		imageFileOld, err := os.OpenFile(fileNameOld, os.O_CREATE|os.O_WRONLY, 0644)
		if nil != err {
			fmt.Printf("[Error] Unable to open file: %v, err: %v\n", fileNameOld, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer os.Remove(fileNameOld)
		digest, _, err := common.CopyAndComputeMD5(imageFileOld, r.Body)
		if nil != err {
			fmt.Printf("[Error] Unable to read request err: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		imageFileOld.Close()
		digestHex := strings.ToUpper(fmt.Sprintf("%x", digest))
		md5Header := r.Header.Get("Content-MD5")
		if len(md5Header) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Content-MD5 missing")
			return
		}
		if md5Header != digestHex {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Content-MD5 not matched")
			return
		}
		convert, err := exec.LookPath("convert")
		if nil != err {
			panic("[Fatal] " + err.Error())
		}
		command := exec.Command(
			convert,
			"-auto-orient",
			"-quality", "50",
			"-resize", "8192>",
			fileNameOld, fileNameNew,
		)
		err = command.Run()
		defer os.Remove(fileNameNew)
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
			fmt.Printf("[Error] Unable to open file: %v, err: %v\n", fileNameNew, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer imageFileNew.Close()
		digest, err = common.ComputeMD5(imageFileNew)
		if err != nil {
			fmt.Printf("[Error] Hit an error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		digestHex = strings.ToUpper(fmt.Sprintf("%x", digest))
		imageFileNew.Seek(0, 0)
		header := w.Header()
		header.Add("Content-Type", "image/heif")
		header.Add("Content-MD5", digestHex)
		w.WriteHeader(http.StatusOK)
		_, err = io.CopyBuffer(w, imageFileNew, make([]byte, 32*1024))
		if err != nil {
			fmt.Printf("[Error] Hit an error: %v\n", err)
		}
	})
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		panic(fmt.Sprintf("[Fatal] Unable to start http server: %v\n", err))
	}
	return
}
