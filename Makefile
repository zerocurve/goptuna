VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'

.DEFAULT_GOAL := help

PKGS := $(shell go list ./...)
SOURCES := $(shell find . -name "*.go" -not -name '*_test.go')
ENV := GO111MODULE=on

.PHONY: setup
setup:  ## Setup for required tools.
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u golang.org/x/tools/cmd/stringer
	go get golang.org/x/tools/cmd/godoc


.PHONY: fmt
fmt: $(SOURCES) ## Formatting source codes.
	@$(GO) goimports -w $^

.PHONY: lint
lint: ## Run golint and go vet.
	@$(ENV) golint -set_exit_status=1 $(PKGS)
	@$(ENV) go vet $(PKGS)
	@$(ENV) misspell $(SOURCES)

.PHONY: test
test:  ## Run tests with race condition checking.
	@$(ENV) go test -race ./...

.PHONY: bench
bench:  ## Run benchmarks.
	@$(ENV) go test -bench=. -run=- -benchmem ./...

.PHONY: coverage
cover:  ## Run the tests.
	@$(ENV) go test -coverprofile=coverage.o ./...
	@$(ENV) go tool cover -func=coverage.o

.PHONY: godoc
godoc: ## Run godoc http server
	@$(ENV) echo "Please open http://localhost:6060/pkg/github.com/c-bata/goptuna/"
	$(ENV) godoc -http=localhost:6060

.PHONY: generate
generate: ## Run go generate
	@$(ENV) go generate ./...

.PHONY: build
build: ## Build example command lines.
	mkdir -p ./bin/
	$(ENV) go build -o ./bin/goptuna -ldflags "$(LDFLAGS)" cmd/main.go
	$(ENV) ./_examples/build.sh

.PHONY: help
help: ## Show help text
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
