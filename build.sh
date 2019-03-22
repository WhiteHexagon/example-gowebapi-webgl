#!/bin/bash
GOOS=js GOARCH=wasm go build -o website/app.wasm
#go build /website/webserver.go -o website/webserver        TODO whats wrong with this?

