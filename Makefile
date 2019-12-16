#
# Some copyright
#

BIN := goproject
BIN_WIN := goproject

PLATFORMS := linux/amd64 linux/arm linux/arm64 linux/ppc64le linux/s390x

# Use git tags to set version
VERSION ?= $(shell git describe --tags --always --dirty)

# Source code directories
SRC := cmd pkg

# Go wrappers
GO := go
GO_BUILD :=$(GOCMD) build
GO_CLEAN :=$(GOCMD) clean
GO_TEST :=$(GOCMD) test
OS=$GOOS
ARCH=$GOARCH

# Directories for builds and tests
BUILD_DIR := bin/$(OS)_$(ARCH)

# For Docker
BUILD_IMAGE ?= golang:1.13.5-alpine
BASE_IMAGE ?= gcr.io/distroless/static

all: test build

build:
    $(GOBUILD) -o $(BUILD_DIR)/$(BIN) -v -ldflags "-X $(go list -m)/pkg/version.Version=${VERSION}"

build-docker:
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/$BIN ./cmd/$BIN.go

test:
    $(GOTEST) test -v -race -cover -covermode=atomic ./...

run:
    $(GOBUILD) -o $(BUILD_DIR)/$(BIN) -v ./...
    ./$(BIN)

$(BUILD_DIRS):
	@mkdir -p $@

version:
	@echo $(VERSION)

clean:
    $(GOCLEAN)
    rm -rf .go bin