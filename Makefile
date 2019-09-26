# Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

SHELL=/bin/bash

VERSION ?= 'v0.1.0'
BUILD_TIME=$(shell date -u +%Y-%m-%dT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

LD_FLAGS= '-X "main.buildTime=$(BUILD_TIME)" -X main.gitCommit=$(GIT_COMMIT) -X main.version=$(VERSION)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

ifdef XDG_CONFIG_HOME
	KUBEON_PLUGINSTUB_DIR ?= ${XDG_CONFIG_HOME}/kubeon/plugins
else ifeq ($(OS),Windows_NT)
	KUBEON_PLUGINSTUB_DIR ?= ${LOCALAPPDATA}/kubeon/plugins
endif

.PHONY: version
version:
	@echo "===> $@"
	@echo $(VERSION)

.PHONY: goenv
goenv:
	@env GO111MODULE=on $(GOINSTALL) github.com/GeertJohan/go.rice
	@env GO111MODULE=on $(GOINSTALL) github.com/GeertJohan/go.rice/rice
	@env GO111MODULE=on $(GOINSTALL) github.com/golang/mock/gomock
	@env GO111MODULE=on $(GOINSTALL) github.com/golang/mock/mockgen
	@env GO111MODULE=on $(GOINSTALL) github.com/golang/protobuf/protoc-gen-go

.PHONY: generate
generate:
	@echo "===> $@"
	@find pkg internal -name fake -type d | xargs rm -rf
	@go generate -v ./pkg/... ./internal/...

.PHONY: vet
vet:
	@echo "===> $@"
	@env go vet ./internal/... ./pkg/...

.PHONY: test
test:
	@echo "===> $@"
	@env go test -cover -v ./internal/... ./pkg/...

.PHONY: build
build:
	@echo "===> $@"
	@mkdir -p ./build
	@env $(GOBUILD) -o build/kubeon $(GO_FLAGS) -v ./cmd/main

.PHONY: release
release:
	@echo "===> $@"
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --follow-tags

.PHONY: changelogs
changelogs:
	hacks/changelogs.sh

.PHONY: clean
clean:
	@echo "===> $@"
	@find pkg internal -name fake -type d | xargs rm -rf
	@rm ./pkg/icon/rice-box.go

.PHONY: ci
ci: generate vet test build



