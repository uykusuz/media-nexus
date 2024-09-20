GOPATH?=$(HOME)/go
SHELL=bash
GOFMT=gofmt

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

GOLINES=$(GOPATH)/bin/golines -t 2 -m 120 --ignored-dirs=openapi -w --chain-split-dots --base-formatter gofmt

BINARY_NAME=media-nexus
TEST_TIMEOUT=0

PATH := $(PATH):$(GOPATH)/bin

# options for Go's build process
export CGO_ENABLED=0
export GO111MODULE=on
export GOFLAGS=$(GO_TAGS)

all: deps compile

deps:
	$(GOGET) -v -t ./... || true
	go install github.com/segmentio/golines@latest

.PHONY: compile
compile:
	$(GOBUILD) -v -o $(BINARY_NAME) .

.PHONY: test
test:
	# count=1: disable test caching
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) -count=1 $(shell go list ./...) | tee gotest-report.out ; exit $${PIPESTATUS[0]}
	cat gotest-report.out | go tool test2json > gotest-report.json

.PHONY: clean
clean:
	rm -f gotest-report.out gotest-report.json golint-report.out
	$(GOCLEAN) -v .
	rm -f $(BINARY_NAME)

.PHONY: lint
lint:
	set -o pipefail; golangci-lint run | tee golint-report.out

.PHONY: apply_format
apply_format:
	$(GOFMT) -s -w .
	$(GOLINES) .

.PHONY: check_format
check_format:
	@echo -n "check format: "
	test -z "$(shell $(GOFMT) -s -l .)"
	@echo -n "check line length: "
	$(GOLINES) --dry-run .
	test "$(shell $(GOLINES) --dry-run . | wc -l)" -eq "0"
