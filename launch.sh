#!/bin/sh

docker kill watchtower
docker kill imageshrink_server

docker system prune

docker run -d \
  --name watchtower \
  --restart always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  containrrr/watchtower imageshrink_server --interval 60

docker run -d \
  --name imageshrink_server \
  --restart always \
  -p 58080:58080 \
  imageshrink/imageshrink_server
