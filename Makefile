FILENAME = rmm
STATICCHECK_VERSION = 2026.2.1
GOVULNCHECK_VERSION = v1.8.0

ifeq ($(OS),Windows_NT)
	FILENAME := $(FILENAME).exe
endif

GO_BIN_DIR := $(or $(shell go env GOBIN),$(shell go env GOPATH)/bin)

.PHONY: build install test lint vulncheck

build:
	go build -o $(FILENAME) .

install:
	@echo "Installing Remindme as $(FILENAME) in $(GO_BIN_DIR)"
	go build -o $(GO_BIN_DIR)/$(FILENAME) .

test:
	go test -race ./...

lint:
	@test -z "$$(gofmt -l .)" || (echo "Files not gofmt-formatted:"; gofmt -l .; exit 1)
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION) ./...

vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...
