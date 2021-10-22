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
TREE?=./luet-packages
export LUET_BIN?=$(LUET)


build:
	go build -ldflags '$(LDFLAGS)' -o bin/

vet:
	go vet ${PKG}

fmt:
	go fmt ${PKG}

test:
	go test -v ${PKG} -race -coverprofile=coverage.txt -covermode=atomic


clean-repo:
ifneq ($(shell id -u), 0)
	@echo "Clean needs to run under root user."
	@exit 1
else
	rm -rf build/ *.tar *.metadata.yaml
endif

build-repo: clean-repo
ifneq ($(shell id -u), 0)
	@echo "Build repo needs to run under root user."
	@exit 1
else
	mkdir -p $(ROOT_DIR)/build
	$(LUET) build --no-spinner --all --tree=$(TREE) --destination $(ROOT_DIR)/build --backend $(BACKEND) --concurrency $(CONCURRENCY) --compression $(COMPRESSION)
endif

create-repo:
ifneq ($(shell id -u), 0)
	@echo "Create repo need to run under root user."
	@exit 1
else
	$(LUET) create-repo --no-spinner --tree "$(TREE)" \
    --output $(ROOT_DIR)/build \
    --packages $(ROOT_DIR)/build \
    --name "luet-cosign" \
    --descr "Luet cosign official repository" \
    --urls "http://localhost:8000" \
    --tree-compression gzip \
    --type http
endif

lint: fmt vet

all: lint test build