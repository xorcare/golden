GOPATH ?= $(HOME)/go
GOBIN = $(GOPATH)/bin

PACKAGE = golden
NAMESPACE = github.com/xorcare/$(PACKAGE)
COVER_FILE ?= coverage.out

# Tools
GOLINT = $(GOBIN)/golint
$(GOLINT):
	GO111MODULE=off go get -u golang.org/x/lint/golint

GOIMPORTS = $(GOBIN)/goimports
$(GOIMPORTS):
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

# Main targets
all: depin check build
.DEFAULT_GOAL := all

.PHONY: bench
bench: ## Run benchmarks
	@go test ./... -bench=. -run="Benchmark*"

.PHONY: build
build: ## Build the project binary
	@go build

.PHONY: check
check: static test ## Check project with static checks and unit tests

$(COVER_FILE):
	@$(MAKE) test

.PHONY: cover
cover: $(COVER_FILE) ## Output coverage in human readable form in html
	@go tool cover -html=$(COVER_FILE)
	@rm -f $(COVER_FILE)

.PHONY: fmt
fmt: ## Run go fmt for the whole project
	@test -z $$(for d in $$(go list -f {{.Dir}} ./...); do go fmt $$d/*.go; done)

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: imports
imports: $(GOIMPORTS) ## Check and fix import section by import rules
	@test -z $$(for d in $$(go list -f {{.Dir}} ./...); do goimports -e -l -local $(NAMESPACE) -w $$d/*.go; done)

.PHONY: lint
lint: $(GOLINT) ## Check the project with lint
	@golint -set_exit_status ./...

.PHONY: static
static: fmt imports vet lint ## Run static checks (fmt, lint, imports, vet, ...) all over the project

.PHONY: test
test: ## Run unit tests
	@go test ./... -coverprofile=$(COVER_FILE) -covermode=atomic $d
	@go tool cover -func=$(COVER_FILE) | grep ^total

.PHONY: testup
testup: ## Run unit tests with golden files update
	@go test -v ./... -update

.PHONY: tools
tools: $(GOLINT) $(GOIMPORTS) ## Install all needed tools, e.g. for static checks

.PHONY: depin
depin: ## Install go mod dependencies, beautify go.mod and go.sum files
	@go mod tidy
	@go mod verify

.PHONY: depup
depup: ## Update go mod dependencies, beautify go.mod and go.sum files
	@go mod download
	@go mod verify

.PHONY: vet
vet: ## Check the project with vet
	@go vet ./...
