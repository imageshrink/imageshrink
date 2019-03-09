package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	imageWriter, _ := bodyWriter.CreateFormFile("image", "image")
	imageFile, _ := os.Open("/home/huahang/Desktop/DSCF6064.jpg")
	_, _ = io.Copy(imageWriter, imageFile)
	_ = bodyWriter.WriteField("digest", "15d8b5958a3c2d34d05ad65e5cebfa19")
	_ = bodyWriter.Close()
	response, err := http.Post("http://127.0.0.1:50000", bodyWriter.FormDataContentType(), bodyBuffer)
	if err != nil {
		fmt.Println("Hit an error: " + err.Error())
		return
	}
	if response.StatusCode != 200 {
		fmt.Println("Failed to upload: " + response.Status)
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Failed to upload: " + string(bodyBytes))
		return
	}
	digest := response.Header.Get("DIGEST")
	mimeType := response.Header.Get("Content-Type")
	fmt.Println("digest: " + digest)
	fmt.Println("mimeType: " + mimeType)
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	_ = ioutil.WriteFile("/home/huahang/Desktop/out.jpg", bodyBytes, 0644)
}
