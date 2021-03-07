Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := -s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)

run: build
	./build/debug/adanos-scheduler

build:
	go build -race -ldflags "$(LDFLAGS)" -o build/debug/adanos-scheduler cmd/scheduler/main.go

.PHONY: build run