GO     := go
PROMU  := $(GOPATH)/bin/promu
PREFIX ?= $(shell pwd)

pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

DOCKER_IMAGE_NAME ?= dellhw_exporter
DOCKER_IMAGE_TAG  ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

all: format style vet test build

build: promu
	@$(PROMU) build --prefix $(PREFIX)

crossbuild: promu
	@$(PROMU) crossbuild

docker:
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

format:
	go fmt $(pkgs)

promu:
	@echo ">> fetching promu"
	@GOOS="$(shell uname -s | tr A-Z a-z)" \
	GOARCH="$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m)))" \
	$(GO) get -u github.com/prometheus/promu

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

tarball: promu
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

test:
	@$(GO) test $(pkgs)

test-short:
	@echo ">> running short tests"
	@$(GO) test -short $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

.PHONY: all build crossbuild docker format promu style tarball test vet
