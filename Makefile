.PHONY: help build run fmt vet test test-v test-cover check check-all ci ci-cover staticcheck install-air clean

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

test-v: ## Run tests with verbose output
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -cover ./...

check: fmt vet ## Run fmt and vet
	@echo "✓ Code formatted and vetted"

check-all: fmt vet staticcheck ## Run fmt, vet, and staticcheck
	@echo "✓ All checks passed"

ci: fmt vet staticcheck test ## Run all checks + tests (CI pipeline)
	@echo ""
	@echo "========================================"
	@echo "✓ All checks passed"
	@echo "✓ All tests passed"
	@echo "========================================"
	@echo "Ready to commit!"

ci-cover: fmt vet staticcheck test-cover ## Run all checks + tests with coverage
	@echo ""
	@echo "========================================"
	@echo "✓ All checks passed"
	@echo "✓ All tests passed (with coverage)"
	@echo "========================================"
	@echo "Ready to commit!"

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
