#!/usr/bin/env bash

mkdir -p target

cd target

cmake \
  -DCMAKE_BUILD_TYPE=MinSizeRel \
  -DCMAKE_INSTALL_PREFIX=$HOME/opt/imageshrink ..

cmake --build .

make install

cd ..
