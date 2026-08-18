ARTIFACT_NAME := external-dns-anexia-webhook

TESTPARALLELISM := 4

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GITSHA  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS := -X 'main.Version=$(VERSION)' -X 'main.Gitsha=$(GITSHA)'

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the webhook binary into bin/
	go build -ldflags "$(LDFLAGS)" -o $(CURDIR)/bin/$(ARTIFACT_NAME) ./cmd/webhook

.PHONY: test
test: ## Run unit tests with coverage
	go test -v -parallel $(TESTPARALLELISM) -timeout 5m -covermode atomic -coverprofile=covprofile ./...

.PHONY: fmt
fmt: ## Apply gci/goimports formatting
	golangci-lint fmt -c .golangci.yml

.PHONY: lint
lint: ## Run golangci-lint and go vet
	golangci-lint run -c .golangci.yml
	go vet ./...

.PHONY: tidy
tidy: ## Tidy go.mod/go.sum
	go mod tidy

.PHONY: clean
clean: ## Remove build and test artifacts
	rm -rf $(CURDIR)/bin $(CURDIR)/dist $(CURDIR)/covprofile
