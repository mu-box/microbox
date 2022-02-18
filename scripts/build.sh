#!/usr/bin/env bash
set -e

# disable cgo for true static binaries (will work on alpine linux)
export CGO_ENABLED=0;

# vfor versioning
getCurrCommit() {
  echo `git rev-parse --short HEAD | tr -d "[ \r\n\']"`
}

getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

getTagDir() {
  echo `getCurrTag | sed -E 's/([0-9]+)\..+/v\1/'`
}

BUILD_DATE=`date -u +%Y%m%dT%H%M%S`
# ^for versioning

# remove any previous builds that may have failed
[ -e "./.build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./.build"

printf "\nBuilding microbox...\n"

# build microbox
gox -ldflags "-s -w -X github.com/mu-box/microbox/util/odin.apiKey=$API_KEY \
			  -X github.com/mu-box/microbox/models.microVersion=$(getCurrTag) \
			  -X github.com/mu-box/microbox/models.microCommit=$(getCurrCommit) \
			  -X github.com/mu-box/microbox/models.microBuild=$BUILD_DATE" \
			  -osarch "darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 linux/s390x windows/amd64" \
        -output="./.build/$(getTagDir)/{{.OS}}/{{.Arch}}/microbox"

printf "\nWriting version file...\n"
echo -n "Microbox Version $(getCurrTag)-$BUILD_DATE ($(getCurrCommit))" > ./.build/$(getTagDir)/version

printf "\nBuilding microbox updater...\n"

# change into updater directory and build microbox updater
cd ./updater && \
  gox -ldflags="-s" \
    -osarch "darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 linux/s390x windows/amd64" \
    -output="../.build/$(getTagDir)/{{.OS}}/{{.Arch}}/microbox-update"

#cd ..

#printf "\nCompacting binaries...\n"
#upx ./.build/$(getTagDir)/*/*/microbox*
