#!/bin/bash

cd $(dirname $0)


unset ANDROID_HOME
export CGO_CFLAGS="-I${JAVA_HOME}/include -I${JAVA_HOME}/include/linux"


go run github.com/danbrough/mobile/cmd/gomobile bind -target=linux/amd64 -x -v -o . \
  github.com/danbrough/mobile/demos/basic/hello

CLASSPATH=.:hello.jar

javac -cp $CLASSPATH Main.java
java -cp $CLASSPATH -Djava.library.path=libs/amd64 Main