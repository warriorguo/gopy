package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/warriorguo/gopy/pkg/compiler"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
	"github.com/warriorguo/gopy/pkg/vm"
)

func main() {
	var verbose = flag.Bool("v", false, "verbose output")
	var disasm = flag.Bool("d", false, "disassemble bytecode before execution")
	
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [bytecode.pyc]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "If no file is provided, starts REPL mode\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}
	
	flag.Parse()
	
	if flag.NArg() == 0 {
		runREPL(*verbose)
		return
	}
	
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	
	bytecodeFile := flag.Arg(0)
	
	file, err := os.Open(bytecodeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	
	code, err := compiler.DeserializeCodeObject(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading bytecode: %v\n", err)
		os.Exit(1)
	}
	
	if *verbose {
		fmt.Printf("Loaded code object: %s\n", code.Name)
	}
	
	if *disasm {
		fmt.Println(code.Disassemble())
	}
	
	vm := vm.NewVM()
	result, err := vm.Run(code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
	
	if *verbose {
		fmt.Printf("Execution completed, result: %s\n", result)
	}
}

func runREPL(verbose bool) {
	fmt.Println("GoPy REPL - Python 2 Interpreter")
	fmt.Println("Type 'exit' or 'quit' to exit")
	
	vm := vm.NewVM()
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		if line == "exit" || line == "quit" {
			break
		}
		
		if err := executeREPLLine(line, vm, verbose); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func executeREPLLine(source string, vm *vm.VM, verbose bool) error {
	l := lexer.NewLexer(source)
	tokens := l.AllTokens()
	
	if verbose {
		fmt.Printf("Tokens: %d\n", len(tokens))
	}
	
	module, err := parser.Parse(tokens)
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}
	
	code, err := compiler.Compile(module)
	if err != nil {
		return fmt.Errorf("compile error: %v", err)
	}
	
	if verbose {
		fmt.Printf("Instructions: %d\n", len(code.Instructions))
	}
	
	result, err := vm.Run(code)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	
	if verbose && result != nil && result.String() != "None" {
		fmt.Printf("Result: %s\n", result)
	}
	
	return nil
}