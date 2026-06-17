#! /usr/bin/make

.DEFAULT_GOAL := build

.PHONY: depend run clean build test release install

DATE8601=$(shell date --utc --iso-8601=seconds)
VERSION=$(shell git describe --tags --always)

# Path to the built binary
BINARY_PATH=mj
# Path to the source code
SOURCE=.
# Linker flags to strip the debugging information
LD_FLAGS_STRIP=-s -w
# Inject some info (could not figure out how to not use the fully qualified package path)
LG_FLAGS_VERSION=-X github.com/MieuxVoter/majority-judgment-cli/version.GitSummary=$(VERSION)
LG_FLAGS_DATE=-X github.com/MieuxVoter/majority-judgment-cli/version.BuildDate=$(DATE8601)

depend:
	go get
	sudo apt install --yes upx

run:
	@echo "(running from source code, at version $(VERSION))"
	@echo "(you won't be able to pass parameters via make though)"
	@echo "(best directly use go like so:    go run . example/example.csv --sort  )\n"
	@go run "$(SOURCE)"

clean:
	rm --force "$(BINARY_PATH)"
	rm --force "$(BINARY_PATH).exe"

build: build-linux-amd64

build-linux-amd64: $(shell find . -name \*.go)
	GOOS=linux GOARCH=amd64 go build \
		-v \
		-ldflags="$(LD_FLAGS_STRIP) $(LG_FLAGS_VERSION) $(LG_FLAGS_DATE)" \
		-o "$(BINARY_PATH)" \
		$(SOURCE)
	@echo "Done building $(BINARY_PATH) at $(shell pwd):"
	@ls -lAhF "$(BINARY_PATH)"

build-windows-amd64: $(shell find . -name \*.go)
	GOOS=windows GOARCH=amd64 go build \
		-ldflags="$(LD_FLAGS_STRIP) $(LG_FLAGS_VERSION) $(LG_FLAGS_DATE)" \
		-o "$(BINARY_PATH).exe" \
		$(SOURCE)
	@echo "Done building $(BINARY_PATH).exe at $(shell pwd):"
	@ls -lAhF "$(BINARY_PATH).exe"

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
