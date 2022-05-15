#!/bin/sh

# GOARCH="amd64" GOOS="linux" go build server/imageshrink_server.go
# mv imageshrink_server docker/imageshrink_server
# docker build docker -t docker.io/imageshrink/imageshrink_server:manifest-amd64 --build-arg ARCH=amd64
# docker push docker.io/imageshrink/imageshrink_server:manifest-amd64

# GOARCH="arm64" GOOS="linux" go build server/imageshrink_server.go
# mv imageshrink_server docker/imageshrink_server
# docker build docker -t :manifest-arm64v8 --build-arg ARCH=arm64v8
# docker push docker.io/imageshrink/imageshrink_server:manifest-arm64v8

docker buildx build -t docker.io/imageshrink/imageshrink_server --push --platform="linux/arm64,linux/amd64" docker