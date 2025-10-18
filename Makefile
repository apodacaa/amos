.PHONY: help build run fmt vet test check check-all staticcheck install-air clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o amos

run: ## Run the app
	go run .

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

test: ## Run tests
	go test ./...

check: fmt vet ## Run fmt and vet
	@echo "✓ Code formatted and vetted"

check-all: fmt vet staticcheck ## Run fmt, vet, and staticcheck
	@echo "✓ All checks passed"

staticcheck: ## Run staticcheck linter
	@which staticcheck > /dev/null || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
	staticcheck ./...

install-air: ## Install air for hot reload
	@echo "Installing air (hot reload)..."
	go install github.com/air-verse/air@latest
	@echo "✓ Air installed"

clean: ## Remove built binaries
	rm -f amos

.DEFAULT_GOAL := help
