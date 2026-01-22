#!/bin/bash

# Install things first:
#apt install upx
#go get github.com/ahmetb/govvv

GOOS=linux GOARCH=amd64 go build \
  -v \
  -ldflags="$(govvv -flags -pkg $(go list ./version)) -s -w" \
  -o mj

GOOS=windows GOARCH=amd64 go build \
  -v \
  -ldflags="$(govvv -flags -pkg $(go list ./version)) -s -w" \
  -o mj.exe

# Compress the linux binary
upx mj
