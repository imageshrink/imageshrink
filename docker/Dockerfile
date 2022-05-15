FROM ubuntu:22.04

RUN apt -y update
RUN apt -y install git
RUN apt -y install golang
RUN apt -y install imagemagick
RUN apt -y upgrade

ADD imageshrink_server_aarch64 /usr/bin/imageshrink_server_aarch64
ADD imageshrink_server_x86_64 /usr/bin/imageshrink_server_x86_64
ADD imageshrink_server /usr/bin/imageshrink_server

ENTRYPOINT ["/usr/bin/imageshrink_server"]
