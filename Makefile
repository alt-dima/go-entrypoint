ifeq ($(strip $(VERSION_STRING)),)
VERSION_STRING := $(shell git rev-parse --short HEAD)
endif

BINDIR    := $(CURDIR)/bin
PLATFORMS := linux/amd64/entrypoint/osusergo*netgo*static_build windows/amd64/entrypoint.exe/osusergo*static_build
BUILDCOMMAND := go build -trimpath -ldflags "-s -w -X main.version=${VERSION_STRING}"
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
label = $(word 3, $(temp))
tags = $(subst *, ,$(word 4, $(temp)))

UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
SHACOMMAND := shasum -a 256
else
SHACOMMAND := sha256sum
endif

.DEFAULT_GOAL := build

.PHONY: release
release: $(PLATFORMS)
$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 $(BUILDCOMMAND) -tags "$(tags)" -o "bin/$(label)"
	$(SHACOMMAND) "bin/$(label)" > "bin/$(label).sha256"

.PHONY: build
build:
	$(BUILDCOMMAND) -o bin/entrypoint

.PHONY: dep
dep:
	go mod download

.PHONY: dev
dev:
	# binary will be $(go env GOPATH)/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.2


.PHONY: lint
lint:
	$(shell go env GOPATH)/bin/golangci-lint run --timeout 5m