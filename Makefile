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
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) -count=1 $(shell go list ./... | grep -v /integrationtests) | tee gotest-report.out ; exit $${PIPESTATUS[0]}
	cat gotest-report.out | go tool test2json > gotest-report.json

.PHONY: test.integration
test.integration:
# count=1: disable test caching
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) -count=1 ./integrationtests/... | tee gotest-report-integration.out ; exit $${PIPESTATUS[0]}
	cat gotest-report-integration.out | go tool test2json > gotest-report-integration.json

.PHONY: clean
clean:
	rm -f gotest-report.out gotest-report.json golint-report.out gotest-report-integration.out gotest-report-integration.json
	$(GOCLEAN) -v .
	rm -f $(BINARY_NAME)

.PHONY: lint
lint:
	set -o pipefail; golangci-lint run | tee golint-report.out

.PHONY: format.apply
format.apply:
	$(GOFMT) -s -w .
	$(GOLINES) .

.PHONY: format.check
format.check:
	@echo -n "check format: "
	test -z "$(shell $(GOFMT) -s -l .)"
	@echo -n "check line length: "
	$(GOLINES) --dry-run .
	test "$(shell $(GOLINES) --dry-run . | wc -l)" -eq "0"

.PHONY: docs
docs:
	swag fmt adapters/primary/ahttp
	swag init -g adapters/primary/ahttp/api.go
