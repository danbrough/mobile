#!/bin/bash

cd `dirname $0`   

go install github.com/danbrough/mobile/cmd/gomobile
go install github.com/danbrough/mobile/cmd/gobind
