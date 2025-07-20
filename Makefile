# Default target
.DEFAULT_GOAL := help

# Directories
FRONTEND_DIR := ./frontend
BACKEND_DIR := ./backend

# Go commands
GO := go

# Node commands
NPM := npm

# Help display
.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Test targets
.PHONY: test ## Run all tests
test: test-frontend test-backend

.PHONY: test-frontend
test-frontend: ## Run frontend tests
	@echo "Testing frontend..."
	${NPM} run test --prefix $(FRONTEND_DIR)

.PHONY: test-backend
test-backend: ## Run backend tests
	@echo "Testing backend..."
	cd $(BACKEND_DIR) && ${GO} test ./...

.PHONY: test-backend-prod
test-backend-prod: ## Run backend tests with production tag
	@echo "Testing backend with prod tag..."
	cd $(BACKEND_DIR) && ${GO} test -tags prod ./...

# Lint targets
.PHONY: lint ## Run all linters
lint: lint-frontend lint-backend

.PHONY: lint-frontend
lint-frontend: ## Lint the frontend
	@echo "Linting frontend..."
	${NPM} run format --prefix $(FRONTEND_DIR)
	${NPM} run lint --prefix $(FRONTEND_DIR)
	${NPM} run type-check --prefix $(FRONTEND_DIR)

.PHONY: lint-backend
lint-backend: ## Lint the backend
	@echo "Linting backend..."
	cd $(BACKEND_DIR) && ${GO} fmt ./...
	cd $(BACKEND_DIR) && ${GO} vet ./...

# Development environment setup
.PHONY: setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Installing frontend dependencies..."
	${NPM} install --prefix $(FRONTEND_DIR)
	${NPM} run build --prefix $(FRONTEND_DIR)
	@echo "Installing backend dependencies..."
	cd $(BACKEND_DIR) && ${GO} mod tidy
	go install github.com/air-verse/air@latest
	@echo "Development environment setup complete."

# Clean target
.PHONY: clean
clean: ## Clean up build artifacts
	@echo "Cleaning up..."
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/.vite
	rm -rf $(FRONTEND_DIR)/node_modules
	rm -rf ../bin/dmx_viewer_backend
