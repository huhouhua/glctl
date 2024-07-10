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
NAME = glctl
OS = linux darwin
architecture = amd64 arm6
VERSION_PACKAGE=github.com/huhouhua/glctl/util/version
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
# include the common make file
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


.DEFAULT_GOAL := help

.PHONY: copyright.verify
copyright.verify: tools.verify.licctl
	@echo "===========> Verifying the boilerplate headers for all files"
	@licctl --check -f $(ROOT_DIR)/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea

.PHONY: copyright.add
copyright.add: tools.verify.licctl
	@echo $(ROOT_DIR)
	@licctl -v -f $(ROOT_DIR)/boilerplate.txt $(ROOT_DIR) --skip-dirs=.idea


.PHONY: install.licctl
install.addlicense:
	@$(GO) install github.com/seacraft/licctl@latest

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: tidy
tidy:
	@$(GO) mod tidy

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
	rm -rf $(OUTPUT_DIR)/*

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"