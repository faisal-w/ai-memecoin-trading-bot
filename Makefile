.PHONY: all build run-trading test clean help

# Build all components
all: build

# Build Go backend
build:
	@echo "Building trading bot..."
	@go build -o bin/trading ./cmd/trading
	@echo "Build complete!"

# Run trading bot
run-trading:
	@./bin/trading

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v
	@echo "Tests complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf output/
	@echo "Clean complete!"

# Show help
help:
	@echo "Meme Coin Trading Bot - Makefile Commands"
	@echo ""
	@echo "  make build              - Build the trading bot"
	@echo "  make run-trading        - Run the meme coin trading bot"
	@echo "  make test               - Run all tests"
	@echo "  make clean              - Clean all build artifacts"
	@echo "  make help               - Show this help message"
