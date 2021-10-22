GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_TAG = $(shell git describe --tags 2>/dev/null || echo "v.0.0.1" )

PKG        := ./...
LDFLAGS    := -w -s
LDFLAGS += -X "github.com/rancher-sandbox/luet-cosign/internal/version.version=${GIT_TAG}"
LDFLAGS += -X "github.com/rancher-sandbox/luet-cosign/internal/version.gitCommit=${GIT_COMMIT}"

LUET?=/usr/bin/luet
BACKEND?=docker
CONCURRENCY?=1
export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
COMPRESSION?=gzip
export LUET_BIN?=$(LUET)


build:
	go build -ldflags '$(LDFLAGS)' -o bin/

vet:
	go vet ${PKG}

fmt:
	go fmt ${PKG}

test:
	go test ${PKG} -race -coverprofile=coverage.txt -covermode=atomic

lint: fmt vet

all: lint test build