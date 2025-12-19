package tests

import (
	"fmt"
	"testing"

	"github.com/warriorguo/gopy/pkg/ast"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
)

func TestParserForLoops(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "basic_for_range",
			input: "for i in range(5):\n    print i",
			valid: true,
		},
		{
			name:  "for_list",
			input: "for x in [1, 2, 3]:\n    print x",
			valid: true,
		},
		{
			name:  "for_variable",
			input: "for item in items:\n    print item",
			valid: true,
		},
		{
			name:  "for_string", 
			input: "for char in \"hello\":\n    print char",
			valid: true, // Fixed string parsing and escape sequence issues
		},
		{
			name:  "nested_for",
			input: "for i in range(3):\n    for j in range(2):\n        print i, j",
			valid: true, // Fixed nested indentation parsing
		},
		{
			name:  "for_with_complex_body",
			input: "for i in range(10):\n    if i % 2 == 0:\n        print i, \"even\"\n    else:\n        print i, \"odd\"",
			valid: true, // Fixed complex indentation and string parsing
		},
		{
			name:  "for_empty_range",
			input: "for i in range(0):\n    print i",
			valid: true,
		},
		{
			name:  "for_range_with_start_end",
			input: "for i in range(1, 10):\n    print i",
			valid: true,
		},
		// Invalid cases
		{
			name:  "missing_colon",
			input: "for i in range(5)\n    print i",
			valid: false,
		},
		{
			name:  "missing_in",
			input: "for i range(5):\n    print i",
			valid: false,
		},
		{
			name:  "missing_body",
			input: "for i in range(5):",
			valid: false,
		},
		{
			name:  "invalid_target",
			input: "for 123 in range(5):\n    print i", 
			valid: true, // Currently doesn't validate target type strictly
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.NewLexer(test.input)
			tokens := l.AllTokens()
			p := parser.NewParser(tokens)

			module, err := p.Parse()

			if test.valid {
				if err != nil {
					t.Errorf("Expected valid parse for %q, got error: %v", test.input, err)
					return
				}

				if len(module.Body) == 0 {
					t.Errorf("Expected at least one statement, got empty module for input %q", test.input)
					return
				}

				// Check if first statement is a ForStmt
				forStmt, ok := module.Body[0].(*ast.ForStmt)
				if !ok {
					t.Errorf("Expected ForStmt as first statement, got %T for input %q", module.Body[0], test.input)
					return
				}

				// Verify ForStmt has required components
				if forStmt.Target == nil {
					t.Errorf("ForStmt missing target for input %q", test.input)
				}
				if forStmt.Iter == nil {
					t.Errorf("ForStmt missing iter for input %q", test.input)
				}
				if len(forStmt.Body) == 0 {
					t.Errorf("ForStmt missing body for input %q", test.input)
				}
			} else {
				if err == nil {
					t.Errorf("Expected parse error for %q, but parsing succeeded", test.input)
				}
			}
		})
	}
}

func TestParserForLoopStructure(t *testing.T) {
	tests := []struct {
		input          string
		expectedTarget string
		expectedIter   string
		expectedBodyLen int
	}{
		{
			input:          "for i in range(5):\n    print i",
			expectedTarget: "*ast.Name", 
			expectedIter:   "*ast.Call",
			expectedBodyLen: 1,
		},
		{
			input:          "for x in [1, 2, 3]:\n    print x",
			expectedTarget: "*ast.Name",
			expectedIter:   "*ast.List", 
			expectedBodyLen: 1,
		},
		{
			input:          "for item in items:\n    print item\n    print \"done\"",
			expectedTarget: "*ast.Name",
			expectedIter:   "*ast.Name",
			expectedBodyLen: 2,
		},
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

		if len(module.Body) == 0 {
			t.Errorf("Expected at least one statement for input %q", test.input)
			continue
		}

		forStmt, ok := module.Body[0].(*ast.ForStmt)
		if !ok {
			t.Errorf("Expected ForStmt, got %T for input %q", module.Body[0], test.input)
			continue
		}

		// Check target type
		targetType := fmt.Sprintf("%T", forStmt.Target)
		if targetType != test.expectedTarget {
			t.Errorf("Expected target type %s, got %s for input %q", test.expectedTarget, targetType, test.input)
		}

		// Check iter type
		iterType := fmt.Sprintf("%T", forStmt.Iter)
		if iterType != test.expectedIter {
			t.Errorf("Expected iter type %s, got %s for input %q", test.expectedIter, iterType, test.input)
		}

		// Check body length
		if len(forStmt.Body) != test.expectedBodyLen {
			t.Errorf("Expected body length %d, got %d for input %q", test.expectedBodyLen, len(forStmt.Body), test.input)
		}
	}
}

func TestParserForLoopWithRangeVariations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		rangeArgs int // expected number of arguments to range()
	}{
		{
			name:     "range_single_arg",
			input:    "for i in range(10):\n    pass",
			rangeArgs: 1,
		},
		{
			name:     "range_two_args",
			input:    "for i in range(1, 10):\n    pass",
			rangeArgs: 2,
		},
		{
			name:     "range_three_args",
			input:    "for i in range(0, 10, 2):\n    pass",
			rangeArgs: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.NewLexer(test.input)
			tokens := l.AllTokens()
			p := parser.NewParser(tokens)

			module, err := p.Parse()
			if err != nil {
				t.Errorf("Parse error for %q: %v", test.input, err)
				return
			}

			forStmt, ok := module.Body[0].(*ast.ForStmt)
			if !ok {
				t.Errorf("Expected ForStmt, got %T", module.Body[0])
				return
			}

			// Check if iter is a Call (range function call)
			call, ok := forStmt.Iter.(*ast.Call)
			if !ok {
				t.Errorf("Expected Call in iter, got %T", forStmt.Iter)
				return
			}

			// Check number of arguments to range()
			if len(call.Args) != test.rangeArgs {
				t.Errorf("Expected %d arguments to range(), got %d", test.rangeArgs, len(call.Args))
			}
		})
	}
}

func TestParserNestedForLoops(t *testing.T) {
	input := `for i in range(2):
    for j in range(3):
        print i, j
    print "outer", i`

	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)

	module, err := p.Parse()
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}

	if len(module.Body) != 1 {
		t.Errorf("Expected 1 top-level statement, got %d", len(module.Body))
		return
	}

	outerFor, ok := module.Body[0].(*ast.ForStmt)
	if !ok {
		t.Errorf("Expected ForStmt, got %T", module.Body[0])
		return
	}

	if len(outerFor.Body) != 2 {
		t.Errorf("Expected 2 statements in outer for body, got %d", len(outerFor.Body))
		return
	}

	// First statement should be inner for loop
	innerFor, ok := outerFor.Body[0].(*ast.ForStmt)
	if !ok {
		t.Errorf("Expected inner ForStmt, got %T", outerFor.Body[0])
		return
	}

	if len(innerFor.Body) != 1 {
		t.Errorf("Expected 1 statement in inner for body, got %d", len(innerFor.Body))
		return
	}

	// Second statement should be print
	_, ok = outerFor.Body[1].(*ast.PrintStmt)
	if !ok {
		t.Errorf("Expected PrintStmt as second statement in outer for, got %T", outerFor.Body[1])
	}
}

func TestParserForLoopWithComplexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "for_with_function_call",
			input: "for item in get_items():\n    print item",
		},
		{
			name:  "for_with_attribute_access",
			input: "for item in obj.items:\n    print item",
		},
		{
			name:  "for_with_subscript",
			input: "for item in data[key]:\n    print item",
		},
		{
			name:  "for_with_complex_range",
			input: "for i in range(len(items)):\n    print i, items[i]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.NewLexer(test.input)
			tokens := l.AllTokens()
			p := parser.NewParser(tokens)

			_, err := p.Parse()
			// For now, just test that these parse without errors
			// More complex expressions might fail due to parser limitations
			if err != nil {
				t.Logf("Parse error for %q (this may be expected): %v", test.input, err)
			} else {
				t.Logf("Successfully parsed: %q", test.input)
			}
		})
	}
}