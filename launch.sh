#!/bin/sh

docker kill imageshrink_server
docker rm imageshrink_server
docker run -d                    \
  --name imageshrink_server      \
  --restart always               \
  --tmpfs /tmp                   \
  --publish 58080:58080          \
  imageshrink/imageshrink_server
