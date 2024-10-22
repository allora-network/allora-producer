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
PWD=$(shell pwd)

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
	docker run -t --rm -v $(PWD):/app -v ~/.cache/golangci-lint/v1.61.0:/root/.cache -w /app golangci/golangci-lint:v1.61.0 golangci-lint run -v

# golangci-lint must be installed locally. It's the fastest way to run the linter.
lint-local:
	@echo "--> Running linter"
	golangci-lint run --timeout=10m


cover:
	mkdir -p coverage
	rm -rf coverage/*
	$(GOTEST) -coverprofile=./coverage/coverage.out.tmp ./...
	grep -v mock ./coverage/coverage.out.tmp | grep -v allora-chain > ./coverage/coverage.out
	$(GOCMD) tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
	$(GOCMD) tool cover -func=./coverage/coverage.out

# Default target
.DEFAULT_GOAL := build
