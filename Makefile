.DEFAULT_GOAL := help

PKGS ?= $(shell go list ./... | grep -v /vendor/)
TEST_FLAGS ?= -race

.PHONY: help
help:
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-12s[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## install go deps
	go mod vendor

.PHONY: build
build: ## build tada
	go build -ldflags "-s -w" -o tada .

.PHONY: install
install: build  ## install tada
	cp tada $(GOPATH)/bin/

.PHONY: test
test: ## run tests
	go test $(TEST_FLAGS) $(PKGS)

.PHONY: vet
vet: ## run go vet
	go vet $(PKGS)

.PHONY: coverage
coverage: ## generate code coverage
	go test $(TEST_FLAGS) -covermode=atomic -coverprofile=coverage.txt $(PKGS)
	go tool cover -func=coverage.txt

.PHONY: lint
lint: ## run golangci-lint
	golangci-lint run
