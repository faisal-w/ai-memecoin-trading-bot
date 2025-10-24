.PHONY: all build-backend build-frontend build run-server run-discord run-telegram clean help

# Build all components
all: build

# Build Go backend
build-backend:
	@echo "Building Go backend..."
	@go build -o bin/server ./cmd/server
	@go build -o bin/discord ./cmd/discord
	@go build -o bin/telegram ./cmd/telegram
	@echo "Backend build complete!"

# Build Rust frontend
build-frontend:
	@echo "Building Rust frontend..."
	@cd frontend && trunk build --release
	@echo "Frontend build complete!"

# Build everything
build: build-backend build-frontend

# Run HTTP server (with frontend)
run-server:
	@./bin/server

# Run Discord bot
run-discord:
	@./bin/discord

# Run Telegram bot
run-telegram:
	@./bin/telegram

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf frontend/dist/
	@rm -rf frontend/target/
	@rm -rf output/
	@echo "Clean complete!"

# Show help
help:
	@echo "Meme Bot - Makefile Commands"
	@echo ""
	@echo "  make build              - Build both backend (Go) and frontend (Rust)"
	@echo "  make build-backend      - Build only the Go backend"
	@echo "  make build-frontend     - Build only the Rust frontend"
	@echo "  make run-server         - Run the HTTP server with frontend"
	@echo "  make run-discord        - Run the Discord bot"
	@echo "  make run-telegram       - Run the Telegram bot"
	@echo "  make clean              - Clean all build artifacts"
	@echo "  make help               - Show this help message"
