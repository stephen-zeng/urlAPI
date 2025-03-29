#!/bin/bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build -o urlAPI

