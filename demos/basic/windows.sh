#!/bin/bash

cd $(dirname $0)

# set JAVA_HOME to a w32 jdk
export JAVA_HOME=/mnt/files2/windows/jdk


go run github.com/danbrough/mobile/cmd/gomobile bind -target=windows/amd64 -x -v -work -o build \
  github.com/danbrough/mobile/demos/basic/hello || exit 1

CLASSPATH=build:build/hello.jar
javac -cp $CLASSPATH Main.java -d build
#java -cp $CLASSPATH -Djava.library.path=build/libs/amd64 Main
