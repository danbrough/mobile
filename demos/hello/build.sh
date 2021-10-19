#!/bin/bash

cd `dirname $0`


#go install golang.org/x/mobile/cmd/gomobile
#go install golang.org/x/mobile/cmd/gobind
#go mod tidy

go run golang.org/x/mobile/cmd/gomobile bind -target=linux/amd64  -v -o build_linux demos/hello/hello