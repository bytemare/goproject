#
# Some copyright
#

# Project path and name
PROJECT_REPO := $(shell go list -m)
PROJECT_NAME := $(shell basename $(PROJECT_REPO))
BINARY  :=  $(PROJECT_NAME)

# Get the list of packages in the project
PACKAGES    := $(shell go list ./...)

# Use git tags to set version, and commit hash to identify the build
VERSION             := $(shell git describe --tags --always --dirty)
COMMIT              := $(shell git rev-parse HEAD)

# Go compiler flags
SYMBOL_NAME_VERSION := "version" # this is the global variable name inside the program code source
SYMBOL_NAME_COMMIT  := "commit" # this is the global variable name inside the program code source

SYMBOLS             =$(shell (go tool nm "$(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY)" | grep "$(BINARY)" | grep " D " | cut -d ' ' -f5))
SYMBOL_PATH_VERSION =$(shell (echo $(SYMBOLS) | tr " " "\n" | grep $(SYMBOL_NAME_VERSION)))
SYMBOL_PATH_COMMIT  =$(shell (echo $(SYMBOLS) | tr " " "\n" | grep $(SYMBOL_NAME_COMMIT)))
LD_BIN_INFO         =-X "'$(SYMBOL_PATH_VERSION)=$(VERSION)'" -X "'$(SYMBOL_PATH_COMMIT)=$(COMMIT)'"
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

#
#   Commands
#

all: build
.PHONY: build install lint test cover version clean

# Create directories and build project
$(BINARY):
	@echo "Creating dirs"
	@mkdir -p $@

values:
	@echo "Values"
	@echo "PROJECT_REPO :" $(PROJECT_REPO)
	@echo "PROJECT_NAME :" $(PROJECT_NAME)
	@echo "BINARY :" $(BINARY)
	@echo "PACKAGES :" $(PACKAGES)
	@echo "VERSION :" $(VERSION)
	@echo "COMMIT :" $(COMMIT)
	@echo "SYMBOL_NAME_VERSION :" $(SYMBOL_NAME_VERSION)
	@echo "SYMBOL_NAME_COMMIT :" $(SYMBOL_NAME_COMMIT)
	@echo "SYMBOLS :" $(SYMBOLS)
	@echo "SYMBOL_PATH_VERSION :" $(SYMBOL_PATH_VERSION)
	@echo "SYMBOL_PATH_COMMIT :" $(SYMBOL_PATH_COMMIT)
	@echo "LD_BIN_INFO :" $(LD_BIN_INFO)
	@echo "LD_STATIC :" $(LD_STATIC)
	@echo "LD_LIGHT :" $(LD_LIGHT)
	@echo "LD_ALL_FLAGS :" $(LD_ALL_FLAGS)
	@echo "OS :" $(OS)
	@echo "ARCH :" $(ARCH)
	@echo "PLATFORMS :" $(PLATFORMS)
	@echo "BUILD_DIR :" $(BUILD_DIR)
	@echo "COVERAGE :" $(COVERAGE)
	@echo "SECCOMP :" $(SECCOMP)

fmt:
	@echo "Formatting"
	@go fmt ./...

lint: fmt
	@echo "Linting ..."
	@golangci-lint run ./...
	@go vet ./...

sec:
	@echo "Checking security"
	@gosec -exclude=G204 ./...

# Build a first time to read symbol location, then a 2nd time with loading LD FLAGS
pre-build: lint sec
	@go build -v \
	    -o $(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY)

build: pre-build
	@echo "Building $(BINARY) in $(BUILD_DIR)/$(OS)_$(ARCH)"

	@go build -v \
	    -ldflags '$(LD_ALL_FLAGS)' \
	    -o $(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY)

install: lint sec
	@echo "Installing $(BINARY) in $(BUILD_DIR)/$(OS)_$(ARCH)""
	@go install -v \
	    -ldflags '$(LD_ALL_FLAGS)' \
	    -o $(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY)

cover: lint sec
	@echo "Coverage"
	@for PACK in $(PACKAGES); do \
		echo "Testing $(PACK)" \
		go test -v -i -race -covermode=atomic \
                	-coverpkg=$(PACKAGES) \
	                -coverprofile=$(COVERAGE)/unit-`echo $$PACK | tr "/" "_"`.out; done

GINKGO ?= $(GOBIN)/ginkgo
test: lint sec build
	@echo "Testing"
	@go get -u github.com/onsi/ginkgo/ginkgo
	@$(GINKGO) -r -v

release: lint sec
	@echo "Releasing"
	goreleaser release

version:
	@echo $(VERSION) - $(COMMIT)

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

build-docker:
    # Build seccomp profile
    #@go2seccomp $(BINARY) $(SECCOMP)
    # Build binary
	@$(DOCKER_GO_COMPILER_ENV) \
	 go build -a -msan -u \
	 --tags netgo \
	 -ldflags $(LD_ALL_FLAGS) \
	 -o $(BINARY)

#docker-run:
#	@docker run \
#        $(DOCKER_LOCKDOWN) \
        --rm \  # remove container after shutdown
#        -i \    # allow sending and receiving on STDIN

