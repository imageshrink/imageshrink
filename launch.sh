#!/bin/sh

docker kill imageshrink_server
docker rm imageshrink_server
docker pull docker.io/imageshrink/imageshrink_server:latest
docker run -d                    \
  --name imageshrink_server      \
  --restart always               \
  --tmpfs /tmp                   \
  --publish 58080:58080          \
  docker.io/imageshrink/imageshrink_server
