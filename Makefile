PKGS := $(shell go list ./...)

all: goclean clean test lint build

clean:
	rm -rf vendor
	rm -rf build

goclean:
	go clean -cache -testcache -i -x -modcache $(PKGS)

test:
	go test

lint:
	go fmt $(PKGS)
	go vet $(PKGS)


prep:
	go mod vendor

blinux:
	mkdir -p build
	GOOS=linux go build -mod vendor -ldflags \
		"-X main.buildHash=$$(git rev-parse HEAD) \
		-X main.buildVersion=$$(git describe --tags) \
		-X main.buildDate=$$(date -u -Iseconds)" \
		-X main.buildGoInfo=$$(go version)" \
			-o build/megantory.linux

bmacos:
	mkdir -p build
	GOOS=darwin go build -mod vendor -ldflags \
		"-X main.buildHash=$$(git rev-parse HEAD) \
		-X main.buildVersion=$$(git describe --tags) \
		-X main.buildDate=$$(date -u -Iseconds)" \
		-X main.buildGoInfo=$$(go version)" \
			-o build/megantory.macos

bwindows:
	mkdir -p build
	GOOS=windows go build -mod vendor -ldflags \
		"-X main.buildHash=$$(git rev-parse HEAD) \
		-X main.buildVersion=$$(git describe --tags) \
		-X main.buildDate=$$(date -u -Iseconds)" \
		-X main.buildGoInfo=$$(go version)" \
			-o build/megantory.windows

build: blinux bmacos bwindows


.PHONY: all clean goclean test lint all build blinux bwindows bmacos
