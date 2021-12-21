#!/bin/bash

cd `dirname $0`

grep -rl golang.org/x/mobile --exclude=.git --exclude=rename.sh | xargs sed -i  's|golang.org/x/mobile|github.com/danbrough/mobile|g'




