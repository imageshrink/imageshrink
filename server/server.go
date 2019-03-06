package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(404)
	})
	http.ListenAndServe(":50000", nil)
}
