.PHONY: help install dev build test lint clean data spike export

help: ## Show this help
	@grep -E '^[a-z]+:.*##' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-10s %s\n", $$1, $$2}'

install: ## Install ui dependencies and Go modules
	cd ui && npm ci
	cd pipeline && go mod download

dev: ## Run the ui dev server
	cd ui && npm run dev

build: ## Build the static site into ui/build
	cd ui && npm run build

test: ## Run the Go tests and svelte-check
	cd pipeline && go test ./...
	cd ui && npm run check

lint: ## go vet + gofmt check + prettier check
	cd pipeline && go vet ./... && test -z "$$(gofmt -l .)"
	cd ui && npm run format:check

data: ## Download pro drafts (cached) into pipeline/data/drafts.json
	cd pipeline && go run ./cmd/fetch

export: ## Write the puzzle set to ui/static/data/puzzles.json
	cd pipeline && go run ./cmd/export

spike: ## Measure how predictable R5 is from baseline signals
	cd pipeline && go run ./cmd/spike

clean: ## Remove build output (keeps the API cache)
	rm -rf ui/build ui/.svelte-kit pipeline/data/drafts.json
