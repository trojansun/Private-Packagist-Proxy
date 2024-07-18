@echo off

REM 为Linux平台编译
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o dist/ppp-linux-amd64 main.go

REM 为Windows平台编译
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o dist/ppp-windows-amd64.exe main.go
