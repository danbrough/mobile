#!/bin/bash

cd $(dirname $0)

# set JAVA_HOME to a w32 jdk
export JAVA_HOME=/mnt/files2/windows/jdk

gomobile bind -target=windows/amd64 -x -v -work -o build \
  github.com/danbrough/mobile/demos/basic/hello || exit 1

