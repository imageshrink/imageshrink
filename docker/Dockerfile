FROM ubuntu:22.04

RUN apt-get -y update
RUN apt-get -y install build-essential
RUN apt-get -y install curl
RUN apt-get -y install emacs
RUN apt-get -y install git
RUN apt-get -y install imagemagick
RUN apt-get -y upgrade

ADD imageshrink_server /usr/bin/imageshrink_server
ENTRYPOINT ["/usr/bin/imageshrink_server"]
