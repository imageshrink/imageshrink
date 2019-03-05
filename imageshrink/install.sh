#!/usr/bin/env bash

mkdir -p target

cd target

cmake -G"Ninja" \
  -DCMAKE_BUILD_TYPE=MinSizeRel \
  -DCMAKE_INSTALL_PREFIX=$HOME/opt/imageshrink ..

cmake --build .

ninja install

cd ..
