PKGS := $(shell go list ./...)

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

blinux:
	mkdir -p build
	GOOS=linux go build -ldflags \
		"-X main.buildHash=$$(git rev-parse HEAD) \
		-X main.buildVersion=$$(git describe --tags) \
		-X main.buildDate=$$(date -u -Iseconds)" \
			-o build/megantory.linux

bmacos:
	mkdir -p build
	GOOS=darwin go build -ldflags \
		"-X main.buildHash=$$(git rev-parse HEAD) \
		-X main.buildVersion=$$(git describe --tags) \
		-X main.buildDate=$$(date -u -Iseconds)" \
			-o build/megantory.macos

bwindows:
	mkdir -p build
	GOOS=windows go build -ldflags \
		"-X main.buildHash=$$(git rev-parse HEAD) \
		-X main.buildVersion=$$(git describe --tags) \
		-X main.buildDate=$$(date -u -Iseconds)" \
			-o build/megantory.windows

build: blinux bmacos bwindows



all: goclean clean test lint build

.PHONY: all clean goclean test lint all build blinux bwindows bmacos
