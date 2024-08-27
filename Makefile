# Copyright 2024 The huhouhua Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build all by default, even if it's not first
GO := go
OS = linux darwin
architecture = amd64 arm64
NAME = glctl
ROOT_PACKAGE=github.com/huhouhua/glctl
VERSION_PACKAGE=github.com/huhouhua/glctl/util/version
GIT_TREE_STATE:="dirty"
COVERAGE := 60
SHELL := /bin/bash
# Linux command settings
FIND := find . ! -path './vendor/*'
XARGS := xargs -r

ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(shell git describe --tags --always --match='v*') \
	-X $(VERSION_PACKAGE).GitCommit=$(shell git rev-parse HEAD) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR) && pwd -P))
endif

ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif

include Makefile.tools.mk

.DEFAULT_GOAL := help

.PHONY: copyright.verify
copyright.verify: tools.verify.licctl
	@echo "===========> Verifying the boilerplate headers for all files"
	@licctl --check -f $(ROOT_DIR)/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea

.PHONY: copyright.add
copyright.add: tools.verify.licctl
	@echo $(ROOT_DIR)
	@licctl -v -f $(ROOT_DIR)/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

## test: Run unit test.
.PHONY: test
test: tools.verify.go-junit-report
	@echo "===========> Run unit test"
	@set -o pipefail;$(GO) test ./cmd/... ./util/...  -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short \

	@sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out # remove mock_.*.go files from test coverage
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html
## cover: Run unit test and get test coverage.
.PHONY: cover
cover: test
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		awk -v target=$(COVERAGE) -f $(ROOT_DIR)/coverage.awk

.PHONY: build
build: clean tidy ## Generate releases for unix systems
	@for arch in $(architecture);\
	do \
		for os in ${OS};\
		do \
			echo "Building $$os-$$arch"; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build  $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/$(NAME)-$$os-$$arch; \
		done \
	done

.PHONY: clean
clean: ## Remove building artifacts
	@echo "===========> Cleaning all build output"
	rm -rf $(OUTPUT_DIR)/*

## tools: install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

.PHONY: go.updates
go.updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'