.PHONY: help install dev build test lint clean data spike

help: ## Show this help
	@grep -E '^[a-z]+:.*##' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-10s %s\n", $$1, $$2}'

install: ## Install web dependencies and Go modules
	cd web && npm ci
	cd pipeline && go mod download

dev: ## Run the web dev server
	cd web && npm run dev

build: ## Build the static site into web/build
	cd web && npm run build

test: ## Run the Go tests and svelte-check
	cd pipeline && go test ./...
	cd web && npm run check

lint: ## go vet + gofmt check + prettier check
	cd pipeline && go vet ./... && test -z "$$(gofmt -l .)"
	cd web && npm run format:check

data: ## Download pro drafts (cached) into pipeline/data/drafts.json
	cd pipeline && go run ./cmd/fetch

spike: ## Measure how predictable R5 is from baseline signals
	cd pipeline && go run ./cmd/spike

clean: ## Remove build output (keeps the API cache)
	rm -rf web/build web/.svelte-kit pipeline/data/drafts.json
