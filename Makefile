# Based on https://git.io/fjkGc

# Define GOPATH in case it is not installed,
# it is necessary to install auxiliary tools.
GOPATH ?= $(HOME)/go
GOPATH_BIN = $(GOPATH)/bin

# The full path to the main package is used in the
# imports tool to format imports correctly.
NAMESPACE = github.com/xorcare/golden

# The name of the file recommended in the standard
# documentation go test -cover and used codecov.io
# to check code coverage.
COVER_FILE ?= coverage.out

# Maximum life time of the compiled tool in minutes.
# This limitation is needed to periodically update
# auxiliary tools and at the same time have a cache
# to speed up statically checks and build.
TOOL_LIFETIME ?= 60

# The section contains installation instructions for auxiliary tools.
GOLINT = $(GOPATH_BIN)/golint
.PHONY: $(GOLINT)
$(GOLINT):
	@if [ $$(find $(GOLINT) -type f -maxdepth 0 -cmin -$(TOOL_LIFETIME) 2> /dev/null | wc -l) -eq 0 ]; \
		then GO111MODULE=off go get -u golang.org/x/lint/golint; fi

GOIMPORTS = $(GOPATH_BIN)/goimports
.PHONY: $(GOIMPORTS)
$(GOIMPORTS):
	@if [ $$(find $(GOIMPORTS) -type f -maxdepth 0 -cmin -$(TOOL_LIFETIME) 2> /dev/null | wc -l) -eq 0 ]; \
		then GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports; fi

# Main targets.
.DEFAULT_GOAL := help

.PHONY: bench
bench: ## Run benchmarks
	@go test ./... -bench=. -run="Benchmark*"

.PHONY: build
build: ## Build the project binary
	@go build ./...

.PHONY: check
check: static test ## Check project with static checks and unit tests

$(COVER_FILE):
	@$(MAKE) test

.PHONY: cover
cover: $(COVER_FILE) ## Output coverage in human readable form in html
	@go tool cover -html=$(COVER_FILE)
	@rm -f $(COVER_FILE)

.PHONY: dep
dep: ## Install and sync go modules dependencies, beautify go.mod and go.sum files
	@go mod download
	@go mod tidy
	@go mod download
	@go mod verify

.PHONY: fmt
fmt: ## Run go fmt for the whole project
	@test -z $$(for d in $$(go list -f {{.Dir}} ./...); do go fmt $$d/*.go; done)

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

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
	@go test ./... $(ARGS) -coverprofile=$(COVER_FILE) -covermode=atomic $d
	@go tool cover -func=$(COVER_FILE) | grep ^total

.PHONY: testin
testin: ## Run integration tests
	@go test ./... -tags=integration -coverprofile=$(COVER_FILE) -covermode=atomic $d
	@go tool cover -func=$(COVER_FILE) | grep ^total

.PHONY: testup
testup: ## Run unit tests with golden files update
	@find . -type f -name '*.golden' -exec rm -f {} \;
	@go test ./... -update

.PHONY: tools
tools: ## Install all needed tools, e.g. for static checks
	@for tool in $(GOLINT) $(GOIMPORTS); \
	do rm -f $$tool && $(MAKE) $$tool; done

.PHONY: vet
vet: ## Check the project with vet
	@go vet ./...
