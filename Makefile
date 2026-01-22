.PHONY: help build build-windows build-all run fmt vet test test-v test-cover check check-all ci ci-cover staticcheck install-air clean gen-test-data release

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -ldflags "-X main.Version=1.7.2" -o amos

build-windows: ## Build Windows binary (amd64)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=1.7.2" -o amos.exe

build-all: ## Build binaries for all platforms
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=1.7.2" -o amos-linux-amd64
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=1.7.2" -o amos-darwin-amd64
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=1.7.2" -o amos-darwin-arm64
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=1.7.2" -o amos-windows-amd64.exe
	@echo "✓ All binaries built"

run: ## Run the app
	go build -ldflags "-X main.Version=1.7.2" -o amos && ./amos

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
	rm -f amos amos.exe amos-*

gen-test-data: ## Generate test data (usage: make gen-test-data ENTRIES=1000 TODOS=500)
	@ENTRIES=$${ENTRIES:-100}; \
	TODOS=$${TODOS:-50}; \
	echo "Generating $$ENTRIES entries and $$TODOS todos..."; \
	go run scripts/generate_test_data.go -entries $$ENTRIES -todos $$TODOS

release: ## Create a new release (usage: make release VERSION=1.2.1 [NOTES=release-notes.md] [DRY_RUN=true])
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION required"; \
		echo "Usage: make release VERSION=1.7.2 [NOTES=release-notes.md] [DRY_RUN=true]"; \
		exit 1; \
	fi
	@./scripts/release.sh $(VERSION) "$(NOTES)" "$(DRY_RUN)"

.DEFAULT_GOAL := help
