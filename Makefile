buildHash := $(shell git rev-parse HEAD)
buildVersion := $(shell git describe --tags)
buildDate := $(shell date -u -Iseconds)

all: clean prep lint test build

clean:
	rm -rf vendor
	rm -rf build

goclean: clean
	go clean -cache -testcache -i -x -modcache

prep:
	go mod vendor

test:
	go test -mod=vendor


lint:
	go fmt
	go vet -mod=vendor

blinux:
	mkdir -p build
	GOOS=linux go build -mod vendor \
		-ldflags "-X main.buildHash=$(buildHash)" \
		-ldflags "-X main.buildVersion=$(buildVersion)" \
		-ldflags "-X main.buildDate=$(buildDate)" \
			-o build/megantory.linux

bmacos:
	mkdir -p build
	GOOS=darwin go build -mod vendor \
		-ldflags "-X main.buildHash=$(buildHash)" \
		-ldflags "-X main.buildVersion=$(buildVersion)" \
		-ldflags "-X main.buildDate=$(buildDate)" \
			-o build/megantory.macos

bwindows:
	mkdir -p build
	GOOS=windows go build -mod vendor \
		-ldflags "-X main.buildHash=$(buildHash)" \
		-ldflags "-X main.buildVersion=$(buildVersion)" \
		-ldflags "-X main.buildDate=$(buildDate)" \
			-o build/megantory.windows

build: blinux bmacos bwindows

.PHONY: all clean goclean test lint all build blinux bwindows bmacos
