#!/bin/bash

cd $(dirname $0)

go run github.com/danbrough/mobile/cmd/gomobile bind -target=linux/amd64 -v -o build demos/basic/hello
