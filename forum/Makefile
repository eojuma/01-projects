# Variables
BINARY_NAME=forum
MAIN_PATH=./cmd/main.go
DB_PATH=./internals/db/data/forum.db

.PHONY: all build run test clean tidy help

## all: Default target to build and run the project
all: tidy build run

## build: Compile the Go binary
build:
	@echo "Building the binary..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

## run: Run the application directly
run:
	@echo "Starting the forum server..."
	go run $(MAIN_PATH)

## test: Run all tests in the project
test:
	@echo "Running tests..."
	go test ./... -v

## clean: Remove the binary and the SQLite database for a fresh start
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f $(DB_PATH)

## tidy: Clean up go.mod and download missing dependencies
tidy:
	@echo "Tidying up modules..."
	go mod tidy

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets: build, run, test, clean, tidy, help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'