# Based on https://git.io/fjkGc

# The full path to the main package is use in the
# imports tool to format imports correctly.
NAMESPACE = github.com/xorcare/golden

# The name of the file recommended in the standard
# documentation go test -cover and used codecov.io
# to check code coverage.
COVER_FILE ?= coverage.out

# Main targets.
.DEFAULT_GOAL := help

.PHONY: bench
bench: ## Run benchmarks
	@go test ./... -bench=. -run="Benchmark*"

.PHONY: build
build: ## Build the project binary
	@go build ./...

.PHONY: ci
ci: check checkstate ## Target for integration with ci pipeline

.PHONY: check
check: static test build ## Check project with static checks and unit tests

$(COVER_FILE):
	@$(MAKE) test

.PHONY: cover
cover: $(COVER_FILE) ## Output coverage in human readable an HTML
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
imports: tools ## Check and fix import section by import rules
	@test -z $$(for d in $$(go list -f {{.Dir}} ./...); do goimports -e -l -local $(NAMESPACE) -w $$d/*.go; done)

.PHONY: lint
lint: tools ## Check the project with lint
	@go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

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

CDTOOLS ?= cd internal/tools &&
.PHONY: tools
tools: ## Install all needed tools, e.g. for static checks
	@$(CDTOOLS) go install golang.org/x/lint/golint
	@$(CDTOOLS) go install golang.org/x/tools/cmd/goimports

.PHONY: toolsup
toolsup: ## Update all needed tools, e.g. for static checks
	@$(CDTOOLS) go mod tidy
	@$(CDTOOLS) go get golang.org/x/lint/golint@latest
	@$(CDTOOLS) go get golang.org/x/tools/cmd/goimports@latest
	@$(CDTOOLS) go mod download
	@$(CDTOOLS) go mod verify
	@$(MAKE) tools

.PHONY: vet
vet: ## Check the project with vet
	@go vet ./...

.PHONY: checkstate
checkstate: ## Checking the relevance of dependencies, and tools. Also, the absence of arbitrary changes when performing checks.
	@echo 'checking the relevance of the dependency list'
	@go mod tidy
	@git diff --exit-code go.mod go.sum
	@echo 'checking the relevance of the committed dependencies'
	@go mod vendor
	@git diff --exit-code vendor
	@echo 'checking the relevance of the committed generated files'
	@go generate
	@exit $$(git status -s | wc -l)
	@go mod verify
