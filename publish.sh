#!/bin/bash

cd `dirname $0`

if [ -f .version ]; then 
	VERSION=$(cat .version)
else 
	VERSION=0
fi

VERSION=$(($VERSION + 1))

VERSION_NAME=`printf "v0.0.1" $VERSION`

echo Creating release: '<'$VERSION_NAME'>'

#while true; do
#    read -p "Do you wish to proceed?: " yn
#    case $yn in
#        [Yy]* ) break;;
#        [Nn]* ) exit;;
#        * ) echo "Please answer yes or no.";;
#    esac
#done

if git tag | grep "$VERSION_NAME" > /dev/null; then
  while true; do
    read -p "Existing Tag $VERSION_NAME found. Shall I delete it?: " yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
  done
  echo removing existing tag "$VERSION_NAME"
  git tag -d "$VERSION_NAME"
  git push origin --delete "$VERSION_NAME"
fi

printf $VERSION > .version
git add .
git commit -am "$VERSION_NAME"
git tag "$VERSION_NAME" && git push && git push origin "$VERSION_NAME"

echo running GOPROXY=proxy.golang.org go list -m  github.com/danbrough/mobile@$VERSION_NAME
GOPROXY=proxy.golang.org go list -m  github.com/danbrough/mobile@$VERSION_NAME




