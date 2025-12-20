BINARY_NAME := coralie-clip
CMD_DIR := ./cmd/coralie-clip
BIN_DIR := ./bin
PREFIX ?= /usr/local
INSTALL_DIR := $(PREFIX)/bin
GO ?= go

.PHONY: help build run test clean install uninstall

help: ## Show this help message
	@echo "Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make build              Build the binary"
	@echo "  make run ARGS='fetch \"Hello\"'  Run from source"
	@echo "  make install           Install to $(INSTALL_DIR)"
	@echo "  make test              Run tests"

build: ## Build the binary
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

run: ## Run from source (use ARGS='...' to pass arguments)
	$(GO) run $(CMD_DIR) -- $(ARGS)

test: ## Run unit tests
	$(GO) test ./...

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)

install: build ## Install binary to system (default: $(INSTALL_DIR))
	@mkdir -p $(INSTALL_DIR)
	install $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo ""
	@echo "Installation complete!"
	@echo ""
	@if ! echo $$PATH | grep -q "$(PREFIX)/bin"; then \
		echo "Note: $(PREFIX)/bin is not in your PATH."; \
		echo ""; \
		echo "To add it, run one of the following:"; \
		echo "  For zsh: echo 'export PATH=\"$(PREFIX)/bin:\$$PATH\"' >> ~/.zshrc && source ~/.zshrc"; \
		echo "  For bash: echo 'export PATH=\"$(PREFIX)/bin:\$$PATH\"' >> ~/.bashrc && source ~/.bashrc"; \
		echo ""; \
		echo "Or open a new terminal window."; \
	fi

uninstall: ## Remove installed binary
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

