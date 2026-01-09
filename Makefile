.PHONY: test bench coverage lint help

# Default target
all: test

## test: run all tests
test:
	go test -v -race ./...

## bench: run benchmarks
bench:
	go test -bench=. -benchmem ./...

## coverage: run tests and generate coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## lint: run golangci-lint (if installed)
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/"; \
		exit 1; \
	fi

## help: show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' |  sed -e 's/^/ /'
