.PHONY: default
default: build lint test

.PHONY: build
build:
	go build -v ./...

GOLANG_TOOL_PATH_TO_BIN=$(shell go env GOPATH)
GOLANGCI_LINT_CLI_VERSION?=latest
GOLANGCI_LINT_CLI_MODULE=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI_LINT_CLI=$(GOLANG_TOOL_PATH_TO_BIN)/bin/golangci-lint
$(GOLANGCI_LINT_CLI):
	go install $(GOLANGCI_LINT_CLI_MODULE)@$(GOLANGCI_LINT_CLI_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT_CLI)
	golangci-lint run

COVER_OUT=coverage.out
COVER_HTML=coverage.html

.PHONY: test
test: dbmigrate-up
	go test -v -cover ./... -coverprofile=$(COVER_OUT)
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

.PHONY: containers-up
containers-up:
	$(MAKE) -C examples/containers up
	sleep 5

.PHONY: dbmigrate-up
dbmigrate-up: containers-up
	$(MAKE) -C examples/dbmigrations up
