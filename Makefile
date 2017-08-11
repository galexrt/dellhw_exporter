PROMU := $(GOPATH)/bin/promu
PREFIX ?= $(shell pwd)

pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

all: build

all: check-license format build test

build: promu
	@$(PROMU) build --prefix $(PREFIX)

crossbuild: promu
	@$(PROMU) crossbuild

test:
	@go test -short $(pkgs)

format:
	go fmt $(pkgs)

promu:
	@go get -u github.com/prometheus/promu

.PHONY: all build crossbuild test format promu
