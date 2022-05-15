#!/bin/sh

GOARCH="amd64" GOOS="linux" go build server/imageshrink_server.go
mv imageshrink_server docker/imageshrink_server_x86_64

GOARCH="arm64" GOOS="linux" go build server/imageshrink_server.go
mv imageshrink_server docker/imageshrink_server_aarch64

docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx rm builder
docker buildx create --name builder --driver docker-container --use
docker buildx build -t docker.io/imageshrink/imageshrink_server --push --platform="linux/arm64,linux/amd64" docker
