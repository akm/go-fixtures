GOOSE_CLI=$(shell go env GOPATH)/bin/goose
$(GOOSE_CLI):
	go install github.com/pressly/goose/v3/cmd/goose@latest

GOOSE_OPTS="root@/go-fixtures-db?parseTime=true"

.PHONY: up
up: $(GOOSE_CLI)
	$(GOOSE_CLI) mysql $(GOOSE_OPTS) up

.PHONY: down
down: $(GOOSE_CLI)
	$(GOOSE_CLI) mysql $(GOOSE_OPTS) down

