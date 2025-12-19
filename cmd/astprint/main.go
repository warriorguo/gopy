package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/warriorguo/gopy/pkg/ast"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
)

func main() {
	var (
		indent   = flag.String("indent", "  ", "Indentation string (default: two spaces)")
		filename = flag.String("f", "", "Python file to parse and print AST")
		help     = flag.Bool("h", false, "Show help")
	)
	flag.Parse()

	if *help {
		fmt.Println("AST Printer - Pretty print Python AST")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  astprint [-indent=string] [-f=file.py] [python_code]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -indent=string    Set indentation string (default: \"  \")")
		fmt.Println("  -f=file.py        Read Python code from file")
		fmt.Println("  -h                Show this help")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  astprint 'x = 1 + 2'")
		fmt.Println("  astprint -f=example.py")
		fmt.Println("  echo 'for i in range(3): print i' | astprint")
		fmt.Println("  astprint -indent='    ' 'def foo(): return 42'")
		return
	}

	var input string
	var err error

	// Determine input source
	if *filename != "" {
		// Read from file
		content, err := os.ReadFile(*filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", *filename, err)
			os.Exit(1)
		}
		input = string(content)
	} else if len(flag.Args()) > 0 {
		// Use command line argument
		input = flag.Args()[0]
	} else {
		// Read from stdin
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(stdinBytes)
	}

	if input == "" {
		fmt.Fprintf(os.Stderr, "No input provided. Use -h for help.\n")
		os.Exit(1)
	}

	// Parse the input
	module, err := parseInput(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	// Format and print the AST
	formatter := ast.NewASTFormatter()
	formatter.SetIndent(*indent)
	
	result := formatter.FormatModule(module)
	fmt.Println(result)
}

func parseInput(input string) (*ast.Module, error) {
	// Tokenize
	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	
	// Parse
	p := parser.NewParser(tokens)
	module, err := p.Parse()
	if err != nil {
		return nil, err
	}
	
	return module, nil
}