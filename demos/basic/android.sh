#!/bin/bash

cd $(dirname $0)


[ ! -d libs ] && mkdir libs

go run github.com/danbrough/mobile/cmd/gomobile bind -target=android/386 -x -v -work -o libs/demo.aar -tags=openssl \
  github.com/danbrough/mobile/demos/basic/hello || exit 1


