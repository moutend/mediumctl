EXECUTABLE_NAME := mediumctl
GOVERSION := $(shell go version)
RELEASE_DIR=bin
DEVTOOL_DIR=_devtools
REVISION=$(shell git rev-parse --verify HEAD | cut -c-6)

.PHONY: clean build build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-386 build-windows-amd64 build-windows-386 $(RELEASE_DIR)/$(EXECUTABLE_NAME)_$(GOOS)_$(GOARCH) all

all: build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-386 build-windows-amd64 build-windows-386

build: $(RELEASE_DIR)/$(EXECUTABLE_NAME)_$(GOOS)_$(GOARCH)

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-386:
	@$(MAKE) build GOOS=linux GOARCH=386

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-386:
	@$(MAKE) build GOOS=darwin GOARCH=386

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64

build-windows-386:
	@$(MAKE) build GOOS=windows GOARCH=386

$(RELEASE_DIR)/$(EXECUTABLE_NAME)_$(GOOS)_$(GOARCH):
ifndef VERSION
	@echo '[ERROR] $$VERSION must be specified'
	exit 255
endif
	go build -ldflags "-X main.revision=$(REVISION) -X main.version=$(VERSION)" \
		-o $(RELEASE_DIR)/$(EXECUTABLE_NAME)_$(GOOS)_$(GOARCH)_$(VERSION) main.go

clean:
	rm -rf $(RELEASE_DIR)/*
