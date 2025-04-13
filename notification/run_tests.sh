#!/bin/bash

# run_tests.sh - Script to run tests and generate coverage report

# Set error handling
set -e

# Print header
echo "=== Running Tests for Notification Service ==="
echo "Started at: $(date)"
echo

# Ensure dependencies are up to date
echo "=== Updating dependencies ==="
go mod tidy
echo

# Run tests with coverage
echo "=== Running tests with coverage ==="
go test -v -race -coverprofile=coverage.out ./...
echo

# Generate coverage report
echo "=== Generating coverage report ==="
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report generated: coverage.html"
echo

# Run specific tests for each package
echo "=== Running specific package tests ==="

echo "API tests:"
go test -v ./api

echo "Store tests:"
go test -v ./store

echo "Stream tests:"
go test -v ./stream

echo "SDK tests:"
go test -v ./sdk

echo "Main tests:"
go test -v .

echo
echo "=== All tests completed ==="
echo "Finished at: $(date)"
