# Copyright 2024 The Kevin Berger <huhouhuam@outlook.com> Authors
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
.DEFAULT_GOAL := all

.PHONY: all
all: tidy add-copyright format lint testdata cover build

# ==============================================================================
# Build options

GO := go
OS = linux darwin
ARCHITECTURE = amd64 arm64
NAME = glctl
ROOT_PACKAGE=github.com/huhouhua/glctl
COVERAGE := 30
SHELL := /bin/bash
DOCKER := docker
GOLANG_CI_LINT_VERSION ?= 2.9.0

# docker command settings
REGISTRY_PREFIX ?= "ghcr.io"
IMAGE ?= "huhouhua/glctl"
VERSION ?= $(shell git describe --tags)
TAG := $(REGISTRY_PREFIX)/$(IMAGE):$(VERSION)
DOCKER_FILE ?= "Dockerfile.dev"
DOCKER_BUILD_ARG_RELEASE ?= $(VERSION)
DOCKER_MULTI_ARCH ?= linux/amd64,linux/arm64

# Linux command settings
FIND := find . ! -path './vendor/*'
XARGS := xargs -r
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR) && pwd -P))
endif

ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif

ifeq ($(origin BIN_DIR),undefined)
BIN_DIR := $(ROOT_DIR)/bin
$(shell mkdir -p $(BIN_DIR))
endif


GO_LDFLAGS := $(shell $(GO) run $(ROOT_DIR)/scripts/gen-ldflags.go)
GO_BUILD_FLAGS = --ldflags "$(GO_LDFLAGS)"

# Copy githook scripts when execute makefile
COPY_GITHOOK:=$(shell cp -f $(ROOT_DIR)/githooks/* $(ROOT_DIR)/.git/hooks/)

# ==============================================================================
# Includes

include scripts/Makefile.tools.mk

# ==============================================================================
# Targets

## verify-copyright: Verify the boilerplate headers for all files.
.PHONY: verify-copyright
verify-copyright: tools.verify.licctl
	@echo "===========> Verifying the boilerplate headers for all files"
	@licctl --check -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea,_output,bin,.github

## add-copyright: Ensures source code files have copyright license headers.
.PHONY: add-copyright
add-copyright: tools.verify.licctl
	@echo $(ROOT_DIR)
	@licctl -v -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea,_output,.github,bin

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
lint: tools.verify.local.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@$(BIN_DIR)/golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

## test: Run unit test.
.PHONY: test
test: tools.verify.go-junit-report run-gitlab
	@echo "===========> Run unit test"
	@set -o pipefail;$(GO) test ./cmd/... ./pkg/...  -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short \

	@sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out # remove mock_.*.go files from test coverage
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

## cover: Run unit test and get test coverage.
.PHONY: cover
cover: test
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		awk -v target=$(COVERAGE) -f $(ROOT_DIR)/scripts/coverage.awk

## build: Build the Go binary for all OS/ARCHITECTURE combinations
.PHONY: build
build: clean tidy
	@for arch in $(ARCHITECTURE);\
	do \
		for os in ${OS};\
		do \
			echo "Building $$os-$$arch"; \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch $(GO) build  $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/$(NAME)-$$os-$$arch; \
		done \
	done

## image.build.%: Build and push a multi-arch Docker image for the specified platform (e.g., image.build.linux_amd64)
.PHONY: image.build.%
image.build.%: build
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building and Pushing $(TAG) for $(OS) $(ARCH) $(ROOT_DIR)/$(DOCKER_FILE)"
	@${DOCKER} buildx build --push -t $(TAG) --build-arg TARGETARCH=${OS}-${ARCH} --build-arg RELEASE=${DOCKER_BUILD_ARG_RELEASE} \
 		--platform $(DOCKER_MULTI_ARCH) -f $(ROOT_DIR)/$(DOCKER_FILE)  $(ROOT_DIR)

## image.push: Push the Docker image to the registry
.PHONY: image.push
image.push:
	@echo "===========> Pushing image $(TAG)"
	@${DOCKER} push $(TAG)

## clean: Remove building artifacts
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	rm -rf $(OUTPUT_DIR)/*

## tools: Install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: Check for outdated direct Go module dependencies
.PHONY: check-updates
check-updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct

## tidy: Clean up go.mod and go.sum by removing unused dependencies and adding missing ones
.PHONY: tidy
tidy:
	@$(GO) mod tidy

## testdata: Run gitlab service and test data for e2e test
.PHONY: testdata
testdata: run-gitlab
	@echo -e "\n\033[36mAdding test data for gitlab conformance tests...\033[0m"
	$(ROOT_DIR)/testdata/seeder.sh

## run-gitlab-e2e: Run gitlab service
.PHONY: run-gitlab
run-gitlab:
	@echo -e "\n\033[36mRunning gitlab conformance tests...\033[0m"
	@$(MAKE) gitlab.install

## kill-gitlab: Kill gitlab service
.PHONY: kill-gitlab
kill-gitlab:
	@echo -e "\n\033[36mKill gitlab conformance tests...\033[0m"
	@$(MAKE) gitlab.uninstall

## start-gitlab: Start gitlab service
.PHONY: start-gitlab
start-gitlab:
	@echo -e "\n\033[36mStart gitlab conformance tests...\033[0m"
	@$(MAKE) gitlab.start

## stop-gitlab: Stop gitlab service
.PHONY: stop-gitlab
stop-gitlab:
	@echo -e "\n\033[36mStop gitlab conformance tests...\033[0m"
	@$(MAKE) gitlab.stop

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
