PROJECTNAME ?= dellhw_exporter
DESCRIPTION ?= dellhw_exporter - Prometheus exporter for Dell Hardware components using OMSA.
MAINTAINER  ?= Alexander Trost <galexrt@googlemail.com>
HOMEPAGE    ?= https://github.com/galexrt/dellhw_exporter

GO111MODULE  ?= on
GO           ?= go
FPM          ?= fpm
PREFIX       ?= $(shell pwd)
BIN_DIR      ?= $(PREFIX)/.bin
TARBALL_DIR  ?= $(PREFIX)/.tarball
PACKAGE_DIR  ?= $(PREFIX)/.package
ARCH         ?= amd64
PACKAGE_ARCH ?= linux-amd64

VERSION      := $(shell cat VERSION)
TOPDIR       := $(shell pwd)

RPMBUILD     := $(TOPDIR)/RPMBUILD

# The GOHOSTARM and PROMU parts have been taken from the prometheus/promu repository
# which is licensed under Apache License 2.0 Copyright 2018 The Prometheus Authors
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))

GOHOSTOS     ?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH   ?= $(shell $(GO) env GOHOSTARCH)

ifeq (arm, $(GOHOSTARCH))
	GOHOSTARM ?= $(shell GOARM= $(GO) env GOARM)
	GO_BUILD_PLATFORM ?= $(GOHOSTOS)-$(GOHOSTARCH)v$(GOHOSTARM)
else
	GO_BUILD_PLATFORM ?= $(GOHOSTOS)-$(GOHOSTARCH)
endif

PROMU_VERSION ?= 0.20.0
PROMU_URL     := https://github.com/prometheus/promu/releases/download/v$(PROMU_VERSION)/promu-$(PROMU_VERSION).$(GO_BUILD_PLATFORM).tar.gz

PROMU := $(FIRST_GOPATH)/bin/promu

WAYBILL_VERSION ?= v0.9.0
WAYBILL_TARGET  ?= x86_64-unknown-linux-gnu
WAYBILL_URL     := https://github.com/kusari-oss/waybill/releases/download/$(WAYBILL_VERSION)/waybill-$(WAYBILL_VERSION)-$(WAYBILL_TARGET).tar.gz
WAYBILL         := $(BIN_DIR)/waybill
SBOM_FILE       ?= dellhw_exporter.sbom.spdx.json
IMAGE_SBOM_FILE ?= dellhw_exporter.image.sbom.spdx.json
# END copied code

pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

DOCKER_IMAGE_NAME ?= dellhw_exporter
DOCKER_IMAGE_TAG  ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

FILELIST := .promu.yml CHANGELOG.md cmd collector contrib docs go.sum Makefile NOTICE README.md VERSION \
            charts CODE_OF_CONDUCT.md container Dockerfile go.mod LICENSE mkdocs.yml pkg systemd


all: format style vet test build

build: promu
	@echo ">> building binaries"
	GO111MODULE=$(GO111MODULE) $(PROMU) build --prefix $(PREFIX)

check_license:
	@OUTPUT="$$($(PROMU) check licenses)"; \
	if [[ $$OUTPUT ]]; then \
		echo "Found go files without license header:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	else \
		echo "All files with license header"; \
	fi

docker:
	@echo ">> building docker image"
	docker build \
		--build-arg BUILD_DATE="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		--build-arg VCS_REF="$(shell git rev-parse HEAD)" \
		-t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" \
		.

format:
	go fmt $(pkgs)

	
tree: build
	mkdir -p -m0755 $(PACKAGE_DIR)/lib/systemd/system $(PACKAGE_DIR)/usr/sbin
	mkdir -p $(PACKAGE_DIR)/etc/sysconfig
	cp dellhw_exporter $(PACKAGE_DIR)/usr/sbin
	cp systemd/dellhw_exporter.service $(PACKAGE_DIR)/lib/systemd/system
	cp systemd/sysconfig.dellhw_exporter $(PACKAGE_DIR)/etc/sysconfig/dellhw_exporter

package-%: tree
	cd $(PACKAGE_DIR) && $(FPM) -s dir -t $(patsubst package-%, %, $@) \
	--deb-user root --deb-group root \
	--name $(PROJECTNAME) \
	--version $(shell cat VERSION) \
	--architecture $(PACKAGE_ARCH) \
	--description "$(DESCRIPTION)" \
	--maintainer "$(MAINTAINER)" \
	--url $(HOMEPAGE) \
	usr/ etc/
	
	
install: build
	mkdir -p -m0755 $(DESTDIR)/usr/lib/systemd/system $(DESTDIR)/usr/sbin
	mkdir -p $(DESTDIR)/etc/sysconfig
	cp dellhw_exporter $(DESTDIR)/usr/sbin
	cp systemd/dellhw_exporter.service $(DESTDIR)/usr/lib/systemd/system
	cp systemd/sysconfig.dellhw_exporter $(DESTDIR)/etc/sysconfig/dellhw_exporter
	

$(PROJECTNAME).spec: $(PROJECTNAME).spec.in
	sed -e's#@VERSION@#$(VERSION)#g' $< > $@
	
rpm: dist $(PROJECTNAME).spec
	mkdir -p $(RPMBUILD)/SOURCES \
	$(RPMBUILD)/SPECS \
	$(RPMBUILD)/BUILD \
	$(RPMBUILD)/RPMS $(RPMBUILD)/SRPMS
	cp $(PROJECTNAME)-$(VERSION).tar.gz $(RPMBUILD)/SOURCES/
	cp $(PROJECTNAME).spec $(RPMBUILD)/SPECS/
	rpmbuild -vv --define "_topdir $(RPMBUILD)" -ba $(PROJECTNAME).spec


promu:
	$(eval PROMU_TMP := $(shell mktemp -d))
	curl -s -L $(PROMU_URL) | tar -xvzf - -C $(PROMU_TMP)
	mkdir -p $(FIRST_GOPATH)/bin
	cp $(PROMU_TMP)/promu-$(PROMU_VERSION).$(GO_BUILD_PLATFORM)/promu $(FIRST_GOPATH)/bin/promu
	rm -r $(PROMU_TMP)

$(WAYBILL):
	@mkdir -p $(BIN_DIR)
	@WAYBILL_TMP=$$(mktemp -d); \
	trap 'rm -rf "$$WAYBILL_TMP"' EXIT; \
	curl --fail --location --retry 3 "$(WAYBILL_URL)" --output "$$WAYBILL_TMP/waybill.tar.gz"; \
	GH_TOKEN="$${GH_TOKEN:-$${GITHUB_TOKEN}}" gh attestation verify "$$WAYBILL_TMP/waybill.tar.gz" --repo kusari-oss/waybill; \
	mkdir "$$WAYBILL_TMP/extracted"; \
	tar -xzf "$$WAYBILL_TMP/waybill.tar.gz" -C "$$WAYBILL_TMP/extracted"; \
	install "$$(find "$$WAYBILL_TMP/extracted" -type f -name waybill -perm -u+x -print -quit)" "$@"

waybill: $(WAYBILL)

sbom: $(WAYBILL)
	@$(WAYBILL) sbom scan --path . --format spdx-2.3-json --output $(SBOM_FILE)

sbom-image: $(WAYBILL)
	@test -n "$(IMAGE)" || (echo "IMAGE must be set, for example IMAGE=ghcr.io/galexrt/dellhw_exporter:v2.0.0" >&2; exit 1)
	@$(WAYBILL) sbom scan --image $(IMAGE) --format spdx-2.3-json --output $(IMAGE_SBOM_FILE)

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

tarball: tree                                                                                                                                       
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(TARBALL_DIR) $(BIN_DIR)


dist: 
	@rm -rf .srcpackage
	@mkdir .srcpackage
	cp -r $(FILELIST) .srcpackage/
	tar --transform "s/\.srcpackage/$(PROJECTNAME)-$(VERSION)/" -zcvf $(PROJECTNAME)-$(VERSION).tar.gz .srcpackage
	
clean:
	rm -rf .srcpackage RPMBUILD $(PROJECTNAME) $(PROJECTNAME).spec $(PROJECTNAME)-$(VERSION).tar.gz 
	
	
test:
	@$(GO) test $(pkgs)

test-short:
	@echo ">> running short tests"
	@$(GO) test -short $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

docs-serve:
	docker run --net=host --volume "$$(pwd)":"$$(pwd)" --workdir "$$(pwd)" -it squidfunk/mkdocs-material

docs-build:
	docker run --net=host --volume "$$(pwd)":"$$(pwd)" --workdir "$$(pwd)" -it squidfunk/mkdocs-material build --clean

helm-docs-binary:
	GO111MODULE=on $(GO) install github.com/norwoodj/helm-docs/cmd/helm-docs@v1.11.0

helm-docs: helm-docs-binary
	helm-docs

.PHONY: all build docker format package promu waybill sbom sbom-image style tarball test vet docs-serve docs-build helm-docs docs
