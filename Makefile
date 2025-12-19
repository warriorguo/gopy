# GoPy Makefile
# Build and test automation for the GoPy Python interpreter

.PHONY: all build test test-unit test-integration test-e2e test-all clean coverage benchmark help

# Default target
all: build test

# Build targets
build: build-py2c build-py2vm build-astprint

build-py2c:
	@echo "Building py2c compiler..."
	@go build -o py2c ./cmd/py2c

build-py2vm:
	@echo "Building py2vm virtual machine..."
	@go build -o py2vm ./cmd/py2vm

build-astprint:
	@echo "Building astprint AST printer..."
	@go build -o astprint ./cmd/astprint

# Test targets
test: test-unit

test-unit:
	@echo "Running unit tests..."
	@go test ./tests -run "TestLexer|TestParserSimple"

test-for:
	@echo "Running FOR keyword specific tests..."
	@go test ./tests -run "TestLexerForKeyword|TestParserForLoops|TestCompilerForLoops|TestVMForLoop|TestForLoopScenarios"

test-integration:
	@echo "Running integration tests..."
	@go test ./tests -run TestEndToEnd

test-e2e:
	@echo "Running end-to-end tests..."
	@go test ./tests -run TestE2E

test-all: test-unit test-integration test-e2e

test-verbose:
	@echo "Running all tests with verbose output..."
	@go test ./tests -v

test-short:
	@echo "Running quick tests..."
	@go test ./tests -short

# Coverage and benchmarking
coverage:
	@echo "Generating test coverage report..."
	@go test ./tests -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

benchmark:
	@echo "Running performance benchmarks..."
	@go test ./tests -run TestE2EBenchmarks -timeout 60s

# Development helpers
lint:
	@echo "Running go fmt and vet..."
	@go fmt ./...
	@go vet ./...

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Example testing
test-examples: build
	@echo "Testing example files..."
	@for example in examples/*.py; do \
		echo "Testing $$example..."; \
		./py2c -o temp.pyc "$$example" && ./py2vm temp.pyc || echo "Failed: $$example"; \
		rm -f temp.pyc; \
	done

test-for-examples: build
	@echo "Testing FOR loop example..."
	@./py2c -o for_test.pyc examples/for_loops.py
	@./py2vm for_test.pyc
	@rm -f for_test.pyc

# Run specific examples
run-hello: build
	@./py2c -o hello.pyc examples/hello.py
	@./py2vm hello.pyc
	@rm -f hello.pyc

run-fibonacci: build  
	@./py2c -o fib.pyc examples/fibonacci.py
	@./py2vm fib.pyc
	@rm -f fib.pyc

run-calculator: build
	@./py2c -o calc.pyc examples/calculator.py  
	@./py2vm calc.pyc
	@rm -f calc.pyc

run-augmented-assignment: build
	@./py2c -o augassign.pyc examples/augmented_assignment.py
	@./py2vm augassign.pyc
	@rm -f augassign.pyc

# Cleanup targets
clean:
	@echo "Cleaning up build artifacts..."
	@rm -f py2c py2vm astprint
	@rm -f *.pyc
	@rm -f coverage.out coverage.html
	@rm -f *_coverage.out

clean-all: clean
	@echo "Cleaning up all generated files..."
	@find . -name "*.pyc" -delete
	@find . -name "*_test.go.tmp" -delete

# Development workflow
dev-setup: deps lint build test

# CI/CD targets
ci: lint build test-all coverage

# Quality assurance
qa: lint test-all coverage benchmark

# Install targets (optional)
install: build
	@echo "Installing GoPy tools to /usr/local/bin..."
	@cp py2c /usr/local/bin/
	@cp py2vm /usr/local/bin/

uninstall:
	@echo "Removing GoPy tools from /usr/local/bin..."
	@rm -f /usr/local/bin/py2c
	@rm -f /usr/local/bin/py2vm

# Help target
help:
	@echo "GoPy Makefile Help"
	@echo "=================="
	@echo ""
	@echo "Build targets:"
	@echo "  build          - Build all tools (py2c, py2vm, astprint)"
	@echo "  build-py2c     - Build only the compiler"
	@echo "  build-py2vm    - Build only the virtual machine"
	@echo "  build-astprint - Build only the AST printer"
	@echo ""
	@echo "Test targets:"
	@echo "  test           - Run unit tests (default)"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests"
	@echo "  test-e2e       - Run end-to-end tests"
	@echo "  test-all       - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-short     - Run quick tests only"
	@echo ""
	@echo "Quality targets:"
	@echo "  coverage       - Generate test coverage report"
	@echo "  benchmark      - Run performance benchmarks"
	@echo "  lint          - Run go fmt and vet"
	@echo ""
	@echo "Example targets:"
	@echo "  test-examples  - Test all example files"
	@echo "  run-hello      - Run hello.py example"
	@echo "  run-fibonacci  - Run fibonacci.py example"
	@echo "  run-calculator - Run calculator.py example"
	@echo "  run-augmented-assignment - Run augmented_assignment.py example"
	@echo ""
	@echo "Utility targets:"
	@echo "  clean          - Remove build artifacts"
	@echo "  clean-all      - Remove all generated files"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  install        - Install tools to /usr/local/bin"
	@echo "  uninstall      - Remove tools from /usr/local/bin"
	@echo ""
	@echo "Workflow targets:"
	@echo "  dev-setup      - Full development setup"
	@echo "  ci             - CI/CD pipeline"
	@echo "  qa             - Quality assurance checks"