#!/bin/bash
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -v -o ping-bot.exe

