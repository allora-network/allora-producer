# Makefile for building the Allora Producer project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate
# Main package path
MAIN_PACKAGE=./cmd/producer

# Binary name
BINARY_NAME=bin/allora-producer

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Get dependencies
deps:
	$(GOGET) -v -t -d ./...

generate:
	$(GOGENERATE) ./...

# Build and run
run: build
	./$(BINARY_NAME)

lint:
	@echo "--> Running linter"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61 run --timeout=10m --fix

# Default target
.DEFAULT_GOAL := build