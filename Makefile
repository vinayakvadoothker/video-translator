# Variables
GO := go
SERVER_ENTRY := cmd/server/main.go
TEST_DIR := tests
BIN_NAME := video-translator

# Run the server directly
run-server:
	@echo "Starting the server..."
	$(GO) run $(SERVER_ENTRY)

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BIN_NAME)

# Build the server binary
build:
	@echo "Building the server binary..."
	$(GO) build -o $(BIN_NAME) $(SERVER_ENTRY)

# Run the server binary
run-binary:
	@echo "Running the server binary..."
	./$(BIN_NAME)

# Help
help:
	@echo "Available commands:"
	@echo "  make run-server   - Run the server directly"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean up build artifacts"
	@echo "  make build        - Build the server binary"
	@echo "  make run-binary   - Run the built server binary"
