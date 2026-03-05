# Makefile for go-run-gopher game (WebAssembly focused)

# Include variable definitions and internal tasks
include mk/variables.mk
include mk/dev.mk
include mk/build.mk
include mk/licenses.mk

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sed 's/^[^:]*://' | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: wasm
wasm: $(DIST_WASM) ## Build WebAssembly binary

.PHONY: wasm-dist
wasm-dist: $(DIST_WASM) $(DIST_WASM_EXEC) $(DIST_INDEX) ## Create distribution package for web hosting
	@echo "✓ Distribution package ready in $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/

.PHONY: serve
serve: wasm-dist ## Start local development server at http://localhost:8080
	@echo "Starting development server at http://localhost:8080"
	@echo "Press Ctrl+C to stop"
	@cd $(DIST_DIR) && python3 -m http.server 8080

.PHONY: run
run: ## Run the game (native version)
	$(GO) run ./$(CMD_DIR)

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "✓ Code formatted"

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "✓ No issues found"

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./...

.PHONY: tidy
tidy: ## Clean up dependencies
	@echo "Tidying dependencies..."
	$(GO) mod tidy
	@echo "✓ Dependencies tidied"

.PHONY: licenses
licenses: $(LICENSES_MD) ## Regenerate third-party licenses markdown

.PHONY: clean
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(DIST_DIR)
	@echo "✓ Cleaned"
