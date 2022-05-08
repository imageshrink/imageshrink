#!/bin/sh

GOARCH="arm64" GOOS="linux" go build server/imageshrink_server.go
mv imageshrink_server docker/imageshrink_server
docker build docker -t imageshrink/imageshrink_server
docker push imageshrink/imageshrink_server
