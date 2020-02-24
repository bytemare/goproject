package internal

type makefile struct {
}

func newMakefile() *makefile {
	return &makefile{}
}

func NewMakefile() (*ProjectFile, error) {
	return NewProjectFile("makefile", "Makefile", newMakefile())
}

func (makefile) getTemplate() string {
	return makefileTemplate
}

const makefileTemplate = `#
# Some copyright
#

BINARY  :=  $(PROJECT_NAME)

# Project path and name
PROJECT_REPO := $(shell go list -m)
PROJECT_NAME := $(shell basename $(PROJECT_REPO))

# Get the list of packages in the project
PACKAGES    := $(shell go list ./...)

# Use git tags to set version, and commit hash to identify the build
VERSION             := $(shell git describe --tags --always --dirty)
COMMIT              := $(shell git rev-parse HEAD)

# Go compiler flags
LD_BIN_INFO         := -X "main.version=$(VERSION)" -X "main.commit=$(COMMIT)"
LD_STATIC           := -extldflags "-static" -installsuffix "-static"
LD_LIGHT            := -s -w
LD_NO_UNSAFE_PKG    := -u
LD_ALL_FLAGS        := '$(LD_BIN_INFO) $(LD_STATIC) $(LD_LIGHT) $(LD_NO_UNSAFE_PKG)'

# Go env vars
OS          := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH        := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
PLATFORMS   ?= linux_amd64 linux_arm64

# Directories for builds and tests
BUILD_DIR   := bin
COVERAGE    := coverage

# Security profiles
SECCOMP     := $BUILD_DIR/$(BINARY).seccomp

#
#   Commands
#

all: build
.PHONY: build install lint test cover version clean

# Create directories and build project
$(BINARY):
@mkdir -p $@

build:
@go build -v \
-ldflags '$(LD_BIN_INFO)' \
-o $(BUILD_DIR)/$(OS)_$(ARCH)/$(BINARY)

install:
@go install $(LD_FLAGS) -v -o $(BUILD_DIR)/$(BINARY)

lint:
@golangci-lint run ./...

test:
@go test -v -i -race $(PACKAGES)

cover:
@for PACK in $(PACKAGES); do \
echo "Testing $(PACK)"
go test -v -i -race -covermode=atomic \
-coverpkg=$(PACKAGES) \
-coverprofile=$(COVERAGE)/unit-` + "echo $$PACK | tr \"/\" \"_\"`" + `.out
done

version:
@echo $(VERSION)

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
@go2seccomp $(BINARY) $(SECCOMP) \
# Build binary
@$(DOCKER_GO_COMPILER_ENV) \
go build -a -msan -u \
--tags netgo \
-ldflags $(LD_ALL_FLAGS)
-o $(BINARY)

docker-run:
@docker run \
$(DOCKER_LOCKDOWN) \
--rm \  # remove container after shutdown
-i \    # allow sending and receiving on STDIN
`
