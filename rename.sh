#!/bin/bash

cd `dirname $0`
find \( -type d -name .git -prune \) -o -type f -print0 | \
	xargs -0 sed -i  's|github.com/danbrough/mobile|github.com/danbrough/mobile|g'



