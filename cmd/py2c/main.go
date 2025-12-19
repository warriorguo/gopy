package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/warriorguo/gopy/pkg/compiler"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
)

func main() {
	var outputFile = flag.String("o", "", "output bytecode file")
	var verbose = flag.Bool("v", false, "verbose output")
	var disasm = flag.Bool("d", false, "disassemble bytecode")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] source.py\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	sourceFile := flag.Arg(0)

	if *outputFile == "" {
		ext := filepath.Ext(sourceFile)
		base := strings.TrimSuffix(sourceFile, ext)
		*outputFile = base + ".pyc"
	}

	if *verbose {
		fmt.Printf("Compiling %s -> %s\n", sourceFile, *outputFile)
	}

	source, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(string(source))
	tokens := l.AllTokens()

	if *verbose {
		fmt.Printf("Lexed %d tokens\n", len(tokens))
	}

	module, err := parser.Parse(tokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Parsed AST with %d statements\n", len(module.Body))
	}

	code, err := compiler.Compile(module)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Compile error: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Generated %d instructions\n", len(code.Instructions))
	}

	if *disasm {
		fmt.Println(code.Disassemble())
	}

	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	err = code.Serialize(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing bytecode: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Successfully compiled to %s\n", *outputFile)
	}
}
