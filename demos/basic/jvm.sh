#!/bin/bash

cd $(dirname $0)



export JAVA_HOME=/opt/jdk/openjdk11

GOMOBILE="go run github.com/danbrough/mobile/cmd/gomobile"
GOMOBILE=gomobile

$GOMOBILE bind -target=linux/amd64 -x -v -o build -tags=openssl \
  github.com/danbrough/mobile/demos/basic/hello || exit 1

CLASSPATH=build:build/hello.jar

javac -cp $CLASSPATH Main.java -d build
java -cp $CLASSPATH -Djava.library.path=build/libs/amd64 Main
