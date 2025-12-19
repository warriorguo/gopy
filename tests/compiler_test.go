package tests

import (
	"testing"

	"github.com/warriorguo/gopy/pkg/compiler"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
)

func TestCompilerBasicExpression(t *testing.T) {
	tests := []struct {
		input              string
		expectedInstrCount int
	}{
		{"42", 2},           // LOAD_CONST, POP_TOP
		{"x", 2},            // LOAD_NAME, POP_TOP
		{"1 + 2", 4},        // LOAD_CONST, LOAD_CONST, BINARY_ADD, POP_TOP
		{"x = 42", 3},       // LOAD_CONST, STORE_NAME, LOAD_CONST (None)
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		tokens := l.AllTokens()
		p := parser.NewParser(tokens)
		
		module, err := p.Parse()
		if err != nil {
			t.Errorf("Parse error for %q: %v", test.input, err)
			continue
		}

		code, err := compiler.Compile(module)
		if err != nil {
			t.Errorf("Compile error for %q: %v", test.input, err)
			continue
		}

		if len(code.Instructions) < test.expectedInstrCount {
			t.Errorf("Expected at least %d instructions, got %d for input %q", 
				test.expectedInstrCount, len(code.Instructions), test.input)
		}
	}
}

func TestCompilerArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected compiler.OpCode
	}{
		{"1 + 2", compiler.OpBinaryAdd},
		{"5 - 3", compiler.OpBinarySub},
		{"2 * 4", compiler.OpBinaryMul},
		{"8 / 2", compiler.OpBinaryDiv},
		{"10 % 3", compiler.OpBinaryMod},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		tokens := l.AllTokens()
		p := parser.NewParser(tokens)
		
		module, err := p.Parse()
		if err != nil {
			t.Errorf("Parse error for %q: %v", test.input, err)
			continue
		}

		code, err := compiler.Compile(module)
		if err != nil {
			t.Errorf("Compile error for %q: %v", test.input, err)
			continue
		}

		// Find the binary operation instruction
		found := false
		for _, instr := range code.Instructions {
			if instr.Op == test.expected {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected instruction %v not found for input %q", test.expected, test.input)
		}
	}
}

func TestCompilerComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected compiler.OpCode
	}{
		{"x == y", compiler.OpCompareEq},
		{"x != y", compiler.OpCompareNe},
		{"x < y", compiler.OpCompareLt},
		{"x <= y", compiler.OpCompareLe},
		{"x > y", compiler.OpCompareGt},
		{"x >= y", compiler.OpCompareGe},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		tokens := l.AllTokens()
		p := parser.NewParser(tokens)
		
		module, err := p.Parse()
		if err != nil {
			t.Errorf("Parse error for %q: %v", test.input, err)
			continue
		}

		code, err := compiler.Compile(module)
		if err != nil {
			t.Errorf("Compile error for %q: %v", test.input, err)
			continue
		}

		// Find the comparison instruction
		found := false
		for _, instr := range code.Instructions {
			if instr.Op == test.expected {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected instruction %v not found for input %q", test.expected, test.input)
		}
	}
}

func TestCompilerPrintStatements(t *testing.T) {
	input := `print "hello", "world"`

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)
	
	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	code, err := compiler.Compile(module)
	if err != nil {
		t.Errorf("Compile error: %v", err)
		return
	}

	// Should contain PRINT_EXPR and PRINT_NEWLINE instructions
	hasExpr := false
	hasNewline := false

	for _, instr := range code.Instructions {
		if instr.Op == compiler.OpPrintExpr {
			hasExpr = true
		}
		if instr.Op == compiler.OpPrintNewline {
			hasNewline = true
		}
	}

	if !hasExpr {
		t.Error("Expected PRINT_EXPR instruction not found")
	}
	if !hasNewline {
		t.Error("Expected PRINT_NEWLINE instruction not found")
	}
}

func TestCompilerAssignments(t *testing.T) {
	input := "x = 42"

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)
	
	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	code, err := compiler.Compile(module)
	if err != nil {
		t.Errorf("Compile error: %v", err)
		return
	}

	// Should contain STORE_NAME instruction
	found := false
	for _, instr := range code.Instructions {
		if instr.Op == compiler.OpStoreName {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected STORE_NAME instruction not found")
	}

	// Should have the variable name in Names
	if len(code.Names) == 0 || code.Names[0] != "x" {
		t.Errorf("Expected variable 'x' in Names, got %v", code.Names)
	}
}

func TestCompilerFunctionCalls(t *testing.T) {
	input := "len([1, 2, 3])"

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)
	
	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	code, err := compiler.Compile(module)
	if err != nil {
		t.Errorf("Compile error: %v", err)
		return
	}

	// Should contain CALL_FUNCTION instruction
	hasCall := false
	hasBuildList := false

	for _, instr := range code.Instructions {
		if instr.Op == compiler.OpCallFunction {
			hasCall = true
		}
		if instr.Op == compiler.OpBuildList {
			hasBuildList = true
		}
	}

	if !hasCall {
		t.Error("Expected CALL_FUNCTION instruction not found")
	}
	if !hasBuildList {
		t.Error("Expected BUILD_LIST instruction not found")
	}
}

func TestCompilerIfStatements(t *testing.T) {
	input := `if x > 0:
    print "positive"`

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)
	
	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	code, err := compiler.Compile(module)
	if err != nil {
		t.Errorf("Compile error: %v", err)
		return
	}

	// Should contain conditional jump instruction
	found := false
	for _, instr := range code.Instructions {
		if instr.Op == compiler.OpPopJumpIfFalse {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected conditional jump instruction not found")
	}
}

func TestCompilerConstants(t *testing.T) {
	input := `x = 42
y = 3.14
name = "hello"`

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)
	
	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	code, err := compiler.Compile(module)
	if err != nil {
		t.Errorf("Compile error: %v", err)
		return
	}

	// Should have constants for 42, 3.14, "hello", and None
	if len(code.Consts) < 4 {
		t.Errorf("Expected at least 4 constants, got %d", len(code.Consts))
	}
}