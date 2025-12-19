# GoPy - Python 2 Interpreter in Go

A simplified Python 2 interpreter implementation in Go, featuring a clean separation between compilation frontend and virtual machine backend.

## Module

```
module github.com/warriorguo/gopy
```

## Architecture

The interpreter follows a traditional compilation pipeline:

- **Frontend**: Lexer → Parser → AST → Compiler → Bytecode
- **Backend**: Virtual Machine executing bytecode instructions

## Supported Python 2 Features

### Data Types
- Basic types: `int`, `float`, `bool`, `str`, `list`, `dict`, `None`
- Type conversion and operations

### Control Flow
- Conditional statements: `if`/`elif`/`else`
- Loop constructs: `while`, `for...in` (with `range()`, lists, strings)
- Function definitions: `def`, `return`, recursive calls

### Operators
- Arithmetic: `+`, `-`, `*`, `/`, `%`, unary `+/-`
- **Augmented assignment**: `+=`, `-=` ✨
- Comparison: `==`, `!=`, `<`, `<=`, `>`, `>=`
- Membership: `in`
- Boolean: `and`, `or`, `not`

### Built-in Functions
- `print` - Output to console
- `len()` - Get length of sequences
- `range()` - Generate integer sequences (1-3 arguments)
- `type()` - Get object type
- `str()` - Convert to string representation

### Advanced Features
- **Recursive function calls** (fixed scope handling) ✨
- Nested function scopes
- Variable scope resolution (local, global, builtin)

## Project Structure

```
github.com/warriorguo/gopy/
├── cmd/
│   ├── astprint/      # AST visualization tool
│   ├── py2c/          # Python to bytecode compiler
│   └── py2vm/         # Bytecode virtual machine
├── pkg/
│   ├── ast/           # AST node definitions and printing
│   ├── compiler/      # Bytecode generation and objects  
│   ├── lexer/         # Tokenization and lexical analysis
│   ├── object/        # Object interface definitions
│   ├── parser/        # AST generation from tokens
│   ├── runtime/       # Built-in Python objects
│   └── vm/            # Virtual machine and execution
├── examples/          # Example Python programs
├── tests/            # Comprehensive test suite
├── Makefile          # Build and test automation
└── README.md
```

## Quick Start

### Build

```bash
# Build all tools
make build

# Or build individually
go build -o py2c ./cmd/py2c
go build -o py2vm ./cmd/py2vm
go build -o astprint ./cmd/astprint
```

### Usage

```bash
# Compile Python to bytecode
./py2c examples/hello.py -o hello.pyc

# Execute bytecode
./py2vm hello.pyc

# Or use Makefile shortcuts
make run-hello
make run-fibonacci  
make run-augmented-assignment
```

### Examples

```python
# Recursive Functions
def fibonacci(n):
    if n <= 1:
        return n
    else:
        return fibonacci(n-1) + fibonacci(n-2)

print fibonacci(6)  # Output: 8

# Augmented Assignment  
x = 10
x += 5    # x becomes 15
x -= 3    # x becomes 12

# For Loops
for i in range(3):
    print i        # Output: 0, 1, 2

for char in "abc":
    print char     # Output: a, b, c

# List operations
numbers = [1, 2, 3]
for num in numbers:
    print num * num
```

## Testing

```bash
# Run unit tests
make test

# Run specific test categories
make test-for          # FOR loop tests
make test-integration  # Integration tests
make test-examples     # Test example files

# Run tests with coverage
make coverage

# Test all examples
make test-examples
```

## Development

```bash
# Full development setup
make dev-setup

# Run quality assurance
make qa

# Clean build artifacts
make clean

# Show all available targets
make help
```

## Requirements

- Go 1.21 or later
- No external dependencies (pure Go implementation)