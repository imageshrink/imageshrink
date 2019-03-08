package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gographics/gmagick"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func checkError(err error, writer http.ResponseWriter) bool {
	if nil == err {
		return false
	}
	writer.WriteHeader(500)
	_, _ = io.WriteString(writer, err.Error() + "\n")
	return true
}

func handlePost(writer http.ResponseWriter, request *http.Request)  {
	var err error
	file, _, err := request.FormFile("image")
	if checkError(err, writer) {
		return
	}
	digest := request.FormValue("digest")
	size, _ := strconv.Atoi(request.FormValue("size"))
	if size == 0 {
		size = 2048
	}
	quality, _ := strconv.Atoi(request.FormValue("quality"))
	if quality == 0 {
		quality = 95
	}
	data, err := ioutil.ReadAll(file)
	if checkError(err, writer) {
		return
	}
	if digest != "" {
		md5Bytes := md5.Sum(data)
		md5Digest := hex.EncodeToString(md5Bytes[:])
		if !strings.EqualFold(md5Digest, digest) {
			writer.WriteHeader(500)
			_, _ = io.WriteString(writer, "Digest no match!\n")
			return
		}
	}
	wand := gmagick.NewMagickWand()
	err = wand.ReadImageBlob(data)
	if checkError(err, writer) {
		return
	}
	width := float64(wand.GetImageWidth())
	height := float64(wand.GetImageHeight())
	var scale = 1.0
	if width > height {
		scale = float64(size) / width
	} else {
		scale = float64(size) / height
	}
	newWidth := uint(width * scale)
	newHeight := uint(height * scale)
	err = wand.ResizeImage(newWidth, newHeight, gmagick.FILTER_LANCZOS, 1)
	if checkError(err, writer) {
		return
	}
	err = wand.SetCompressionQuality(uint(quality))
	if checkError(err, writer) {
		return
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 16)
	random := strconv.FormatUint(rand.Uint64(), 16)
	filename := "/tmp/imageshrink2_" + timestamp + "_" + random + ".jpg"
	err = wand.WriteImage(filename)
	if checkError(err, writer) {
		return
	}
	data, err = ioutil.ReadFile(filename)
	if checkError(err, writer) {
		return
	}
	err = os.Remove(filename)
	if checkError(err, writer) {
		return
	}
	digestBytes := md5.Sum(data)
	digest = hex.EncodeToString(digestBytes[:])
	writer.Header().Add("Content-Type", "image/jpeg")
	writer.Header().Add("Content-Disposition", "inline")
	writer.Header().Add("DIGEST", digest)
	_, err = writer.Write(data)
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	const bodyHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Image Shrink</title>
</head>
<body>


<form method="post" action="/" enctype="multipart/form-data">
    <p>Choose file to upload</p>
    <input type="file" name="image">
    <p>Digest</p>
    <input type="text" name="digest">
    <p>Size</p>
    <input type="text" name="size">
    <p>Quality</p>
    <input type="text" name="quality">
    <p></p>
    <input type="submit" value="Upload Image" name="submit">
</form>

</body>
</html>
`
	writer.WriteHeader(200)
	_, _ = io.WriteString(writer, bodyHTML)
}

func main() {
	gmagick.Initialize()
	defer gmagick.Terminate()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			handlePost(writer, request)
			return
		} else if request.Method == "GET" {
			handleGet(writer, request)
			return
		}
		writer.WriteHeader(405)
		_, _ = io.WriteString(writer, "405 Method Not Allowed")
	})
	_ = http.ListenAndServe(":50000", nil)
}
