#!/bin/bash

cd `dirname $0`


go install github.com/danbrough/mobile/cmd/gomobile@latest
go install github.com/danbrough/mobile/cmd/gobind@latest
go run github.com/danbrough/mobile/cmd/gomobile  bind -target=linux/amd64 \
  -v -o build github.com/danbrough/mobile/demos/basic/hello




