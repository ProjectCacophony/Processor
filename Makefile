.PHONY: all clean build lint

NAME ?= processor
SOURCE ?= ./cmd/processor

GOOS ?= linux
GOARCH ?= amd64
BUILD_DIR ?= bin/$(GOOS).$(GOARCH)

BINARY= $(BUILD_DIR)/$(NAME)
BUILD_FLAGS=

SOURCE_FOLDERS := $(shell go list -f {{.Dir}} ./...)

all: build

clean:
	rm -Rf bin/

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v $(BUILD_FLAGS) -o "$(BINARY)" $(SOURCE)

build_windows:
	export GOOS=windows
	go build -v -o "bin/windows.$(GOARCH)/$(NAME).exe" $(SOURCE)

lint:
	golangci-lint run ./...
