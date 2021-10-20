#!/bin/bash

#go get github.com/ahmetb/govvv

go build \
  -ldflags="$(govvv -flags -pkg $(go list ./version)) -s -w" \
  -o mj

GOOS=windows GOARCH=amd64 go build \
  -ldflags="$(govvv -flags -pkg $(go list ./version)) -s -w" \
  -o mj.exe

upx mj

