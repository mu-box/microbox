TAG=$(shell git describe --always --tags --abbrev=0 | tr -d "[v\r\n]")
TAGDIR=$(shell echo $(TAG) | sed -E 's/([0-9]+)\..+/v\1/')
COMMIT=$(shell git rev-parse --short HEAD | tr -d "[ \r\n\']")
BUILD_DATE=$(shell date -u +%Y%m%dT%H%M%S)
GITSTATUS=$(shell git status 2> /dev/null | tail -n1)
DIRTY=$(shell [ "$(GITSTATUS)" = 'nothing to commit, working tree clean' ] || echo -n "+")
GO_LDFLAGS="-s -X github.com/mu-box/microbox/util/odin.apiKey=$(API_KEY) -X github.com/mu-box/microbox/models.microVersion=$(TAG) -X github.com/mu-box/microbox/models.microCommit=$(COMMIT)$(DIRTY) -X github.com/mu-box/microbox/models.microBuild=$(BUILD_DATE)"

default: all

local: linux windows darwin

clean:
	@echo "Cleaning old builds"
	@rm -rf "./.build"

all: build

# go get github.com/mitchellh/gox
build: clean
	@echo "Building microbox"
	@gox -ldflags=$(GO_LDFLAGS) -osarch "darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 linux/s390x windows/amd64" -output="./.build/$(TAGDIR)/{{.OS}}/{{.Arch}}/microbox"
	@echo -n "Microbox Version $(TAG)-$(BUILD_DATE) ($(COMMIT)$(DIRTY))" > ./.build/$(TAGDIR)/version
	@echo "Building microbox-update"
	@cd ./updater && gox -osarch "darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 linux/s390x windows/amd64" -ldflags="-s" -output="../.build/$(TAGDIR)/{{.OS}}/{{.Arch}}/microbox-update"


native:
	@echo "Building microbox (native os/arch)"
	@go build -ldflags=$(GO_LDFLAGS) -o microbox


linux: linux-amd64 linux-arm32 linux-arm64 linux-s390x

linux-amd64:
	@echo "Building microbox-linux-amd64"
	@GOOS=linux GOARCH=amd64 go build -ldflags=$(GO_LDFLAGS) -o microbox-linux-amd64

linux-arm32:
	@echo "Building microbox-linux-arm32"
	@GOOS=linux GOARCH=arm go build -ldflags=$(GO_LDFLAGS) -o microbox-linux-arm32

linux-arm64:
	@echo "Building microbox-linux-arm64"
	@GOOS=linux GOARCH=arm64 go build -ldflags=$(GO_LDFLAGS) -o microbox-linux-arm64

linux-s390x:
	@echo "Building microbox-linux-s390x"
	@GOOS=linux GOARCH=s390x go build -ldflags=$(GO_LDFLAGS) -o microbox-linux-s390x


windows: windows-amd64

windows-amd64:
	@echo "Building microbox-windows-amd64"
	@GOOS=windows GOARCH=amd64 go build -ldflags=$(GO_LDFLAGS) -o microbox-windows-amd64


darwin: darwin-amd64 darwin-arm64

darwin-amd64:
	@echo "Building microbox-darwin-amd64"
	@GOOS=darwin GOARCH=amd64 go build -ldflags=$(GO_LDFLAGS) -o microbox-darwin-amd64

darwin-arm64:
	@echo "Building microbox-darwin-arm64"
	@GOOS=darwin GOARCH=arm64 go build -ldflags=$(GO_LDFLAGS) -o microbox-darwin-arm64


# go get github.com/kardianos/govendor
test: 
	@go test -v ./...


.PHONY: test clean all build native linux windows windows-amd64 darwin darwin-amd64 darwin-arm64
