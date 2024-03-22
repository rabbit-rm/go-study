Maker := "Yozosoft.com"
#BuildTime := $(shell date +"%Y-%m-%d %H:%M:%S")
BuildTime := $(shell powershell date -Format "yyyy-MM-dd")T$(shell powershell date -Format "HH:mm:ss")

dir := $(subst /,\,$(CURDIR))
prefix := $(dir)\manager
source := $(wildcard $(dir)\*.go)

ldflags := -ldflags "-w -s -X 'main.Maker=$(Maker)' -X 'main.BuildTime=$(BuildTime)'"
tags := -tags="jsoniter"
amdTags := -tags="sonic avx"
#delete an existing executable file,arg1:file path
exist_del = @$(if $(wildcard $1), (del $1 & echo $1,"del..."),)
# set the go environment variable,arg1:os,arg2:arch
env = @echo "env,GOOS=$(1) GOARCH=$(2)" & go env -w CGO_ENABLED=0 GOOS=$1 GOARCH=$2
# go build ,arg1:file name of the build,arg2:is need amd tags
build = @echo "build,$(1)" & go build $(ldflags) $(if $2,$(amdTags),$(tags)) -o $1 $(source)
# upx compressing executables,arg1:file name,arg2:force-macos
upx = @echo "compression,$(1)" & "C:\developmentEnv\BuildTools\upx-win64-4.2.2\upx.exe" $(if $2,--force-macos,) -9q $1


all: prepare windows Linux-amd64 Linux-amd64-noavx Linux-arm64 darwin-arm64 Linux-loong64 Linux-mips64 Linux-mips64le

windows:
	$(call exist_del,$(prefix)-windows-amd64.exe)
	$(call env,windows,amd64)
	$(call build,$(prefix)-windows-amd64.exe,true)
	$(call upx,$(prefix)-windows-amd64.exe)
	$(call env,windows,amd64)
darwin-arm64:
	$(call exist_del,$(prefix)-darwin-arm64)
	$(call env,darwin,arm64)
	$(call build,$(prefix)-darwin-arm64)
	$(call upx,$(prefix)-darwin-arm64,true)
	$(call env,windows,amd64)
Linux-amd64:
	$(call exist_del,$(prefix)-linux-amd64)
	$(call env,linux,amd64)
	$(call build,$(prefix)-linux-amd64,true)
	$(call upx,$(prefix)-linux-amd64)
	$(call env,windows,amd64)
Linux-amd64-noavx:
	$(call exist_del,$(prefix)-linux-amd64-noavx)
	$(call env,linux,amd64)
	$(call build,$(prefix)-linux-amd64-noavx)
	$(call upx,$(prefix)-linux-amd64-noavx)
	$(call env,windows,amd64)
Linux-arm64:
	$(call exist_del,$(prefix)-linux-arm64)
	$(call env,linux,arm64)
	$(call build,$(prefix)-linux-arm64)
	$(call upx,$(prefix)-linux-arm64)
	$(call env,windows,amd64)
Linux-loong64:
	$(call exist_del,$(prefix)-linux-loong64)
	$(call env,linux,loong64)
	$(call build,$(prefix)-linux-loong64)
	$(call env,windows,amd64)
Linux-mips64:
	$(call exist_del,$(prefix)-linux-mips64)
	$(call env,linux,mips64)
	$(call build,$(prefix)-linux-mips64)
	$(call env,windows,amd64)
Linux-mips64le:
	$(call exist_del,$(prefix)-linux-mips64le)
	$(call env,linux,mips64le)
	$(call build,$(prefix)-linux-mips64le)
	$(call env,windows,amd64)

prepare:
	@go mod tidy
	@go-bindata -o=internal/resource/bindata.go -pkg=resource _resource/...
	@echo "Maker: [${Maker}]"
	@echo "BUILD_TIME: [${BuildTime}]"

test:
	echo "current time:$(BuildTime)"
