#!/bin/bash

# 切换到项目根目录

# 为Linux平台编译
env GOOS=linux GOARCH=amd64 go build -o dist/ppp-linux-amd64 main.go

# 为Windows平台编译
env GOOS=windows GOARCH=amd64 go build -o dist/ppp-windows-amd64.exe main.go
