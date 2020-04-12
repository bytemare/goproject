#
# Some copyright
#

# Target main files
TARGETS  := "./cmd/goproject"

# Project path and name
PROJECT_REPO := $(shell go list -m)
PROJECT_NAME := $(shell basename $(PROJECT_REPO))
BINARY  :=  $(PROJECT_NAME)

# Get the list of packages in the project
PACKAGES    := $(shell go list ./...)

# Use git tags to set version, and commit hash to identify the build
VERSION             := $(shell git describe --tags --always --dirty)
COMMIT              := $(shell git rev-parse HEAD)
DATE                := $(shell date +'%F_%T:%N')

# Go compiler flags
# these are the global variables' names inside the program code source
SYMBOL_NAME_VERSION := "version"
SYMBOL_NAME_COMMIT  := "commit"
SYMBOL_NAME_DATE    := "date"

SYMBOLS             =$(shell (go tool nm "$(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY)" | grep "$(BINARY)" | rev | grep " D " | cut -d ' ' -f1 | rev)) # rev is a hack to get the last field
SYMBOL_PATH_VERSION =$(shell (echo $(SYMBOLS) | tr " " "\n" | grep "\.$(SYMBOL_NAME_VERSION)$$"))
SYMBOL_PATH_COMMIT  =$(shell (echo $(SYMBOLS) | tr " " "\n" | grep "\.$(SYMBOL_NAME_COMMIT)$$"))
SYMBOL_PATH_DATE    =$(shell (echo $(SYMBOLS) | tr " " "\n" | grep "\.$(SYMBOL_NAME_DATE)$$"))
LD_BIN_INFO         = -X "'$(SYMBOL_PATH_VERSION)=$(VERSION)'" -X "'$(SYMBOL_PATH_COMMIT)=$(COMMIT)'" -X "'$(SYMBOL_PATH_DATE)=$(DATE)'"
LD_STATIC           := -extldflags "-static" -installsuffix "-static"
LD_LIGHT            := -s -w
#LD_NO_UNSAFE_PKG    := -u
LD_ALL_FLAGS        = $(LD_BIN_INFO) $(LD_STATIC) $(LD_LIGHT)

# Go env vars
OS          := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH        := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
PLATFORMS   := linux_amd64 linux_arm64

# Directories for builds and tests
BUILD_DIR   := bin
COVERAGE    := coverage

# Security profiles
SECCOMP     := $(BUILD_DIR)/$(BINARY).seccomp

# Make flags
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

#
#   Commands
#

all:

# Install tools and check environment

.PHONY: prepare-lint
prepare-lint:
	@echo "Installing golangci-lint ..."
	path := $(which golangci-lint)
	ifeq (echo $?,1)
		@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin latest
	endif


.PHONY: prepare-python3
prepare-python3:
	@echo "Installing python3  ..."
	@sudo apt-get -y install python3.5 python3-pip python3-dev python3-setuptools

.PHONY: prepare-pre-commit
prepare-pre-commit: prepare-lint
	@echo "Installing pre-commit ..."
	path := $(which pre-commit)
    ifeq (echo $?,1)
    	@pip3 install --upgrade pip
        @pip3 install pre-commit
        @pre-commit install
    endif

.PHONY: prepare-tests
prepare-tests:
	@echo "Installing libs for tests ..."
	@go get -u github.com/onsi/ginkgo/ginkgo
	@go get -u github.com/onsi/gomega/...

GINKGO ?= $(GOBIN)/ginkgo

.PHONY: setup-ci-linux
setup-ci-linux:
	prepare-python3
	prepare-pre-commit

.PHONY: setup-ci-mac
setup-ci-mac:
	brew install pre-commit golangci-lint

.PHONY: setup-ci
setup-ci:
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		setup-ci-linux
	endif
	ifeq ($(UNAME_S),Darwin)
		setup-ci-mac
	endif


# Create directories
.PHONY: dirs
dirs:
	@echo "Creating dirs ..."
	@mkdir -v -p $(BUILD_DIR) $(COVERAGE)

.PHONY: values
values:
	@echo "Values"
	@echo "PROJECT_REPO :" $(PROJECT_REPO)
	@echo "PROJECT_NAME :" $(PROJECT_NAME)
	@echo "BINARY :" $(BINARY)
	@echo "PACKAGES :" $(PACKAGES)
	@echo "VERSION :" $(VERSION)
	@echo "COMMIT :" $(COMMIT)
	@echo "DATE :" $(DATE)
	@echo "SYMBOL_NAME_VERSION :" $(SYMBOL_NAME_VERSION)
	@echo "SYMBOL_NAME_COMMIT :" $(SYMBOL_NAME_COMMIT)
	@echo "SYMBOL_NAME_DATE :" $(SYMBOL_NAME_DATE)
	@echo "SYMBOLS :" $(SYMBOLS)
	@echo "SYMBOL_PATH_VERSION :" $(SYMBOL_PATH_VERSION)
	@echo "SYMBOL_PATH_COMMIT :" $(SYMBOL_PATH_COMMIT)
	@echo "SYMBOL_PATH_DATE :" $(SYMBOL_PATH_DATE)
	@echo "LD_BIN_INFO :" $(LD_BIN_INFO)
	@echo "LD_STATIC :" $(LD_STATIC)
	@echo "LD_LIGHT :" $(LD_LIGHT)
	@echo "LD_ALL_FLAGS :" '$(LD_ALL_FLAGS)'
	@echo "OS :" $(OS)
	@echo "ARCH :" $(ARCH)
	@echo "PLATFORMS :" $(PLATFORMS)
	@echo "BUILD_DIR :" $(BUILD_DIR)
	@echo "COVERAGE :" $(COVERAGE)
	@echo "SECCOMP :" $(SECCOMP)

.PHONY: fmt
fmt:
	@echo "Formatting ..."
	@go fmt ./...

.PHONY: lint
lint: fmt
	@echo "Linting and security ..."
	@go vet ./...
	@golangci-lint run --fix ./...

.PHONY: pre-commit
pre-commit:
	@echo "Extensive checking with pre-commit ..."
	@pre-commit run --all-files

# Build a first time to read symbol location, then a 2nd time with loading LD FLAGS
.PHONY: pre-build
pre-build: dirs
	@echo "Vanilla build of $(BINARY) in $(BUILD_DIR)/$(OS)_$(ARCH)"
	@go build -v \
	    -o $(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY) \
	    $(TARGETS)

.PHONY: build
build: | lint pre-build
	@echo "Fetching symbols and Building $(BINARY) in $(BUILD_DIR)/$(OS)_$(ARCH) with flags"
	@go build -v \
	    -ldflags '$(LD_ALL_FLAGS)' \
	    -o $(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY) \
	    $(TARGETS)

.PHONY: install
install: lint pre-build
	@echo "Installing $(BINARY) ..."
	@go install -v \
	    -ldflags '$(LD_ALL_FLAGS)' \
	    $(TARGETS)
	@echo "Installed at $(shell which $(BINARY))"

.PHONY: uninstall
uninstall:
	@echo "Uninstall $(TARGETS) ... TODO."

.PHONY: cover
cover:
	@echo "Coverage ..."
	@for PACK in $(PACKAGES); do \
		echo "Testing $(PACK)" \
		go test -v -i -race -covermode=atomic \
		    -coverpkg=$(PACKAGES) \
		    -coverprofile=$(COVERAGE)/unit-`echo $$PACK | tr "/" "_"`.out; done

.PHONY: test
test:
	@echo "Testing ... TODO."
#@$(GINKGO) -r -v

.PHONY: release
release: lint
	@echo "Releasing"
	goreleaser release

.PHONY: version
version:
	@echo $(VERSION) - $(COMMIT)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR) $(COVERAGE)

#
#   Docker
#

# Docker run arguments
DOCKER_GO_COMPILER_ENV  := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
DOCKER_NO_NETWORK       := --net=none      # disable network access
DOCKER_DROP_CAP         := --cap-drop=all  # drop all process capabilities
DOCKER_ATTACH_STD       := --attach=STDIN --attach=STDOUT --attach=STDERR
DOCKER_RO               := --read-only     # make the root file system read-only
DOCKER_NO_PRIV          := --security-opt="no-new-privileges"
DOCKER_SECC_PROFILE     := --security-opt="seccomp=$seccomp"
DOCKER_LOCKDOWN         := $(DOCKER_NO_NETWORK) $(DOCKER_DROP_CAP) $(DOCKER_RO)
DOCKER_SHARE_FS         := --mount type=bind,source="$file",target=/"$base",readonly

.PHONY: build-docker
build-docker:
    # Build seccomp profile
    #@go2seccomp $(BINARY) $(SECCOMP)
    # Build binary
	@$(DOCKER_GO_COMPILER_ENV) \
	 go build -a -msan -u \
	    --tags netgo \
	    -ldflags '$(LD_ALL_FLAGS)' \
	    -o $(BINARY) \
	    $(TARGETS)

.PHONY: docker-run
docker-run:
	@docker run \
	    $(DOCKER_LOCKDOWN) \
	    --rm -i
