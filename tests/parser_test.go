package tests

import (
	"fmt"
	"testing"

	"github.com/warriorguo/gopy/pkg/ast"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
)

func TestParserBasicExpressions(t *testing.T) {
	tests := []struct {
		input        string
		expectedType string
	}{
		{"42", "*ast.Num"},
		{"3.14", "*ast.Num"},
		{`"hello"`, "*ast.Str"},
		{"True", "*ast.NameConstant"},
		{"x", "*ast.Name"},
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

		if len(module.Body) != 1 {
			t.Errorf("Expected 1 statement, got %d for input %q", len(module.Body), test.input)
			continue
		}

		exprStmt, ok := module.Body[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("Expected ExprStmt, got %T for input %q", module.Body[0], test.input)
			continue
		}

		actualType := fmt.Sprintf("%T", exprStmt.Expr)
		if actualType != test.expectedType {
			t.Errorf("Expected expression type %q, got %q for input %q", test.expectedType, actualType, test.input)
		}
	}
}

func TestParserBinaryOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + 2", "(1 + 2)"},
		{"x - y", "(x - y)"},
		{"a * b", "(a * b)"},
		{"10 / 5", "(10 / 5)"},
		{"x % 3", "(x % 3)"},
		{"a == b", "(a == b)"},
		{"x < y", "(x < y)"},
		{"x > y", "(x > y)"},
		{"x <= y", "(x <= y)"},
		{"x >= y", "(x >= y)"},
		{"x != y", "(x != y)"},
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

		if len(module.Body) != 1 {
			t.Errorf("Expected 1 statement, got %d for input %q", len(module.Body), test.input)
			continue
		}

		exprStmt, ok := module.Body[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("Expected ExprStmt, got %T for input %q", module.Body[0], test.input)
			continue
		}

		if exprStmt.Expr.String() != test.expected {
			t.Errorf("Expected expression %q, got %q for input %q", test.expected, exprStmt.Expr.String(), test.input)
		}
	}
}

func TestParserAssignments(t *testing.T) {
	tests := []struct {
		input string
		name  string
		value string
	}{
		{"x = 42", "x", "42"},
		{"name = \"hello\"", "name", `"hello"`},
		{"result = a + b", "result", "(a + b)"},
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

		if len(module.Body) != 1 {
			t.Errorf("Expected 1 statement, got %d for input %q", len(module.Body), test.input)
			continue
		}

		assignStmt, ok := module.Body[0].(*ast.AssignStmt)
		if !ok {
			t.Errorf("Expected AssignStmt, got %T for input %q", module.Body[0], test.input)
			continue
		}

		if assignStmt.Target.String() != test.name {
			t.Errorf("Expected target %q, got %q for input %q", test.name, assignStmt.Target.String(), test.input)
		}

		if assignStmt.Value.String() != test.value {
			t.Errorf("Expected value %q, got %q for input %q", test.value, assignStmt.Value.String(), test.input)
		}
	}
}

func TestParserPrintStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int // number of arguments
	}{
		{`print "hello"`, 1},
		{`print x, y`, 2},
		{`print "x =", x`, 2},
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

		if len(module.Body) != 1 {
			t.Errorf("Expected 1 statement, got %d for input %q", len(module.Body), test.input)
			continue
		}

		printStmt, ok := module.Body[0].(*ast.PrintStmt)
		if !ok {
			t.Errorf("Expected PrintStmt, got %T for input %q", module.Body[0], test.input)
			continue
		}

		if len(printStmt.Values) != test.expected {
			t.Errorf("Expected %d arguments, got %d for input %q", test.expected, len(printStmt.Values), test.input)
		}
	}
}

func TestParserIfStatements(t *testing.T) {
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

	if len(module.Body) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(module.Body))
		return
	}

	ifStmt, ok := module.Body[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("Expected IfStmt, got %T", module.Body[0])
		return
	}

	if ifStmt.Test.String() != "(x > 0)" {
		t.Errorf("Expected test '(x > 0)', got %q", ifStmt.Test.String())
	}

	if len(ifStmt.Body) != 1 {
		t.Errorf("Expected 1 body statement, got %d", len(ifStmt.Body))
	}
}

func TestParserForStatements(t *testing.T) {
	input := `for i in range(5):
    print i`

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)
	
	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	if len(module.Body) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(module.Body))
		return
	}

	forStmt, ok := module.Body[0].(*ast.ForStmt)
	if !ok {
		t.Errorf("Expected ForStmt, got %T", module.Body[0])
		return
	}

	if forStmt.Target.String() != "i" {
		t.Errorf("Expected target 'i', got %q", forStmt.Target.String())
	}

	if len(forStmt.Body) != 1 {
		t.Errorf("Expected 1 body statement, got %d", len(forStmt.Body))
	}
}

func TestParserFunctionDefinitions(t *testing.T) {
	input := `def add(a, b):
    return a + b`

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)
	
	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	if len(module.Body) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(module.Body))
		return
	}

	funcDef, ok := module.Body[0].(*ast.FuncDef)
	if !ok {
		t.Errorf("Expected FuncDef, got %T", module.Body[0])
		return
	}

	if funcDef.Name != "add" {
		t.Errorf("Expected function name 'add', got %q", funcDef.Name)
	}

	if len(funcDef.Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(funcDef.Args))
	}

	if len(funcDef.Body) != 1 {
		t.Errorf("Expected 1 body statement, got %d", len(funcDef.Body))
	}
}

func TestParserLists(t *testing.T) {
	tests := []struct {
		input    string
		expected int // number of elements
	}{
		{"[]", 0},
		{"[1, 2, 3]", 3},
		{"[x, y]", 2},
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

		if len(module.Body) != 1 {
			t.Errorf("Expected 1 statement, got %d for input %q", len(module.Body), test.input)
			continue
		}

		exprStmt, ok := module.Body[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("Expected ExprStmt, got %T for input %q", module.Body[0], test.input)
			continue
		}

		listExpr, ok := exprStmt.Expr.(*ast.List)
		if !ok {
			t.Errorf("Expected List, got %T for input %q", exprStmt.Expr, test.input)
			continue
		}

		if len(listExpr.Elts) != test.expected {
			t.Errorf("Expected %d elements, got %d for input %q", test.expected, len(listExpr.Elts), test.input)
		}
	}
}