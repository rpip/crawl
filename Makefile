# Project name
PROJECT = crawl
BUILDTAGS=

.PHONY: clean all fmt vet lint build test install static deps docker
.DEFAULT: default

all: clean build fmt lint test vet install

build:
	@echo "+ $@"
	@go build -tags "$(BUILDTAGS) cgo" .

static:
	@echo "+ $@"
	CGO_ENABLED=1 go build -tags "$(BUILDTAGS) cgo static_build" -ldflags "-w -extldflags -static" -o reg .

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	@echo "+ $@"
	@golint ./... | grep -v vendor | tee /dev/stderr

test: fmt lint vet
	@echo "+ $@"
	@go test -v -tags "$(BUILDTAGS) cgo" $(shell go list ./... | grep -v vendor)

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor)

clean:
	@echo "+ $@"
	@rm -rf reg

install:
	@echo "+ $@"
	@go install .

deps:
	@echo "Installing dependencies..."
	@deps ensure

docker:
	@docker build . -t $(PROJECT)