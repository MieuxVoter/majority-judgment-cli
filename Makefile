.DEFAULT_GOAL := build

#.PHONY: build test

VERSION=$(shell git describe --tags --always)

# Path to the built binary
BINARY_PATH=mj
# Path to the source code
SOURCE=.
# Linker flags to strip the debugging information
LD_FLAGS_STRIP=-s -w


depend:
	go get
	go get github.com/ahmetb/govvv
	sudo apt install -y upx

run:
	@echo "(running from source code, at version $(VERSION))"
	@echo "(you won't be able to pass parameters via make though)"
	@echo "(best directly use go like so:    go run . example/example.csv --sort  )\n"
	@go run .

clean:
	rm -f "$(BINARY_PATH)"
	rm -f "$(BINARY_PATH).exe"

build: build-linux-amd64

build-linux-amd64: $(shell find . -name \*.go)
	GOOS=linux GOARCH=amd64 go build \
		-v \
		-ldflags="$(govvv -flags -pkg $(go list ./version)) $(LD_FLAGS_STRIP)" \
		-o "$(BINARY_PATH)" \
		$(SOURCE)
	@echo "Done building $(BINARY_PATH) at $(shell pwd):"
	@ls -lahF "$(BINARY_PATH)"

build-windows-amd64: $(shell find . -name \*.go)
	GOOS=windows GOARCH=amd64 go build \
		-ldflags="$(govvv -flags -pkg $(go list ./version)) $(LD_FLAGS_STRIP)" \
		-o "$(BINARY_PATH).exe" \
		$(SOURCE)
	@echo "Done building $(BINARY_PATH).exe at $(shell pwd):"
	@ls -lahF "$(BINARY_PATH).exe"

release: clean build-linux-amd64 build-windows-amd64
	upx --ultra-brute "$(BINARY_PATH)"
	upx --ultra-brute "$(BINARY_PATH).exe"

test: test-unit

test-unit:
	go test -v `go list ./...`

install: install-release

install-debug: build
	sudo install "$(BINARY_PATH)" /usr/local/bin/

install-release: release
	sudo install "$(BINARY_PATH)" /usr/local/bin/
