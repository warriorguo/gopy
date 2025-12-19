package tests

import (
	"testing"

	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
)

func TestParserSimple(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		// Basic expressions
		{"42", true},
		{"3.14", true},
		{`"hello"`, true},
		{"True", true},
		{"x", true},

		// Binary operations  
		{"1 + 2", true},
		{"x - y", true},
		{"a * b", true},
		{"10 / 5", true},
		{"x % 3", true},

		// Comparisons
		{"a == b", true},
		{"x < y", true},
		{"x > y", true},
		{"x <= y", true},
		{"x >= y", true},

		// Assignments
		{"x = 42", true},
		{`name = "hello"`, true},
		{"result = a + b", true},

		// Print statements
		{`print "hello"`, true},
		{`print x, y`, true},

		// If statements
		{`if x > 0:
    print "positive"`, true},

		// For loops
		{`for i in range(5):
    print i`, true},

		// Function definitions
		{`def add(a, b):
    return a + b`, true},

		// Lists
		{"[]", true},
		{"[1, 2, 3]", true},

		// Dictionaries
		{"{}", true},
		{`{"a": 1, "b": 2}`, true},

		// Invalid syntax
		{"if x > 0 print", false},
		{"def ()", false},
		{"[1, 2,", false},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		tokens := l.AllTokens()
		p := parser.NewParser(tokens)
		
		_, err := p.Parse()
		
		if test.valid && err != nil {
			t.Errorf("Expected valid parse for %q, got error: %v", test.input, err)
		}
		
		if !test.valid && err == nil {
			t.Errorf("Expected parse error for %q, but parsing succeeded", test.input)
		}
	}
}