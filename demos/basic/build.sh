#!/bin/bash

cd `dirname $0`


go get -d  github.com/danbrough/mobil@latest
#go get -d  github.com/danbrough/mobile/cmd/gobind@latest
#go get -d github.com/danbrough/mobile/cmd/gomobile@latest
go install github.com/danbrough/mobile/cmd/gomobile@latest
go install github.com/danbrough/mobile/cmd/gobind@latest
go run github.com/danbrough/mobile/cmd/gomobile@latest bind -target=linux/amd64  -v -o build demos/basic/hello

