PACKAGE=github.com/moutend/mediumctl
EXECUTABLE_NAME := $(shell pwd | sed "s/.*\/\(.*\)$$/\1/")
GOVERSION := $(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
RELEASE_DIR=bin
DEVTOOL_DIR=_devtools
REVISION=$(shell git rev-parse --verify HEAD | cut -c-6)

.PHONY: test clean build build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-386 build-windows-amd64 build-windows-386 $(RELEASE_DIR)/$(EXECUTABLE_NAME)_$(GOOS)_$(GOARCH) all

all: install-deps build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-386 build-windows-amd64 build-windows-386
fvl: fmt vet lint test

fmt:
	@go fmt

vet:
	@go vet

lint:
	@golint -min_confidence=0.1

test:
	@go test | tee -a log

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

$(DEVTOOL_DIR)/$(GOOS)/$(GOARCH)/glide:
ifndef HAVE_GLIDE
	@echo "Installing glide for $(GOOS)/$(GOARCH)..."
	mkdir -p $(DEVTOOL_DIR)/$(GOOS)/$(GOARCH)
	wget -q -O - https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-$(GOOS)-$(GOARCH).tar.gz | tar xvz
	mv $(GOOS)-$(GOARCH)/glide $(DEVTOOL_DIR)/$(GOOS)/$(GOARCH)/glide
	rm -rf $(GOOS)-$(GOARCH)
endif

glide: $(DEVTOOL_DIR)/$(GOOS)/$(GOARCH)/glide

install-deps: glide
	@PATH=$(DEVTOOL_DIR)/$(GOOS)/$(GOARCH):$(PATH) glide install

clean:
	rm -rf $(RELEASE_DIR)/*
	rm -rf testdata/json/count
	rm -rf testdata/json/issue
	rm -rf testdata/json/issues
