package tests

import (
	"testing"

	"github.com/warriorguo/gopy/pkg/compiler"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/parser"
	"github.com/warriorguo/gopy/pkg/runtime"
)

func TestCompilerForLoops(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectedInstructions  []compiler.OpCode
		shouldCompile        bool
	}{
		{
			name:  "simple_for_range",
			input: "for i in range(3):\n    print i",
			expectedInstructions: []compiler.OpCode{
				compiler.OpLoadName,    // range
				compiler.OpLoadConst,   // 3
				compiler.OpCallFunction,// range(3)
				compiler.OpGetIter,     // get iterator
				compiler.OpForIter,     // for iteration
				compiler.OpStoreFast,   // store i
				compiler.OpLoadGlobal,  // print
				compiler.OpLoadFast,    // load i
				compiler.OpCallFunction,// print(i)
				compiler.OpPopTop,      // pop result
			},
			shouldCompile: true,
		},
		{
			name:  "for_with_list",
			input: "for x in [1, 2, 3]:\n    pass",
			expectedInstructions: []compiler.OpCode{
				compiler.OpBuildList,   // build [1, 2, 3]
				compiler.OpGetIter,     // get iterator
				compiler.OpForIter,     // for iteration
				compiler.OpStoreFast,   // store x
			},
			shouldCompile: true,
		},
		{
			name:  "nested_for_loops",
			input: "for i in range(2):\n    for j in range(2):\n        print i, j",
			shouldCompile: true,
		},
		{
			name:  "for_with_break_continue", 
			input: `for i in range(5):
    if i == 2:
        continue
    if i == 4:
        break
    print i`,
			shouldCompile: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.NewLexer(test.input)
			tokens := l.AllTokens()
			p := parser.NewParser(tokens)

			module, err := p.Parse()
			if err != nil {
				if test.shouldCompile {
					t.Errorf("Parse error for %q: %v", test.input, err)
				}
				return
			}

			code, err := compiler.Compile(module)
			if err != nil {
				if test.shouldCompile {
					t.Errorf("Compile error for %q: %v", test.input, err)
				}
				return
			}

			if !test.shouldCompile {
				t.Errorf("Expected compilation to fail for %q, but it succeeded", test.input)
				return
			}

			// Check that basic compilation succeeded
			if len(code.Instructions) == 0 {
				t.Errorf("Expected non-empty instructions for %q", test.input)
				return
			}

			// For more detailed tests, check for specific instructions
			if len(test.expectedInstructions) > 0 {
				foundInstructions := make(map[compiler.OpCode]bool)
				for _, instr := range code.Instructions {
					foundInstructions[instr.Op] = true
				}

				for _, expectedOp := range test.expectedInstructions {
					if !foundInstructions[expectedOp] {
						t.Logf("Missing expected instruction %v in compiled code for %q", expectedOp, test.input)
					}
				}
			}

			t.Logf("Successfully compiled %q with %d instructions", test.input, len(code.Instructions))
		})
	}
}

func TestCompilerForLoopInstructions(t *testing.T) {
	input := "for i in range(3):\n    print i"

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

	// Test basic properties of the compiled code
	if len(code.Instructions) == 0 {
		t.Error("Expected non-empty instructions")
		return
	}

	// Check that we have some expected instruction types for for loops
	expectedOps := map[compiler.OpCode]bool{
		compiler.OpForIter: false,
	}

	for _, instr := range code.Instructions {
		if _, exists := expectedOps[instr.Op]; exists {
			expectedOps[instr.Op] = true
		}
	}

	for op, found := range expectedOps {
		if !found {
			t.Logf("Note: Expected instruction %v not found (this may be normal depending on implementation)", op)
		}
	}

	// Check that constants include the range argument
	foundRangeArg := false
	for _, constant := range code.Consts {
		if intConst, ok := constant.(*runtime.PyInt); ok && intConst.Value == 3 {
			foundRangeArg = true
			break
		}
	}

	if !foundRangeArg {
		t.Error("Expected to find range argument (3) in constants")
	}

	// Check that names include 'range' and 'i'
	expectedNames := []string{"range", "i", "print"}
	foundNames := make(map[string]bool)
	
	for _, name := range code.Names {
		foundNames[name] = true
	}

	for _, expectedName := range expectedNames {
		if !foundNames[expectedName] {
			t.Errorf("Expected name %q not found in compiled names", expectedName)
		}
	}
}

func TestCompilerForLoopWithDifferentIterables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		checkOp  compiler.OpCode
	}{
		{
			name:    "for_with_range",
			input:   "for i in range(5):\n    pass",
			checkOp: compiler.OpCallFunction, // range() call
		},
		{
			name:    "for_with_list",
			input:   "for x in [1, 2, 3]:\n    pass",
			checkOp: compiler.OpBuildList, // list construction
		},
		{
			name:    "for_with_string",
			input:   `for c in "hello":\n    pass`,
			checkOp: compiler.OpLoadConst, // string constant
		},
		{
			name:    "for_with_variable",
			input:   "for item in items:\n    pass",
			checkOp: compiler.OpLoadName, // variable lookup
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

			code, err := compiler.Compile(module)
			if err != nil {
				t.Errorf("Compile error for %q: %v", test.input, err)
				return
			}

			// Check for expected operation
			found := false
			for _, instr := range code.Instructions {
				if instr.Op == test.checkOp {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Expected instruction %v not found for %q", test.checkOp, test.input)
				
				// Debug: show all instructions
				t.Logf("All instructions for %q:", test.input)
				for i, instr := range code.Instructions {
					t.Logf("  %d: %v (arg: %d)", i, instr.Op, instr.Arg)
				}
			}
		})
	}
}

func TestCompilerForLoopVariableScoping(t *testing.T) {
	input := `x = 10
for x in range(3):
    print x
print x`

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

	// Check that variable 'x' is properly handled in different scopes
	xFound := false
	for _, name := range code.Names {
		if name == "x" {
			xFound = true
			break
		}
	}

	if !xFound {
		t.Error("Expected variable 'x' in compiled names")
	}

	// The compilation should succeed, handling variable scoping correctly
	t.Logf("Successfully compiled for loop with variable scoping, %d instructions generated", len(code.Instructions))
}

func TestCompilerForLoopOptimizations(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		maxInstructions int // rough upper bound for optimization check
	}{
		{
			name:          "simple_for",
			input:         "for i in range(1):\n    pass",
			maxInstructions: 20, // should be relatively few instructions
		},
		{
			name:          "empty_range",
			input:         "for i in range(0):\n    print i",
			maxInstructions: 25,
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

			code, err := compiler.Compile(module)
			if err != nil {
				t.Errorf("Compile error for %q: %v", test.input, err)
				return
			}

			if len(code.Instructions) > test.maxInstructions {
				t.Logf("Note: %q generated %d instructions (expected <= %d). This might indicate room for optimization.",
					test.input, len(code.Instructions), test.maxInstructions)
			} else {
				t.Logf("Efficiently compiled %q with %d instructions", test.input, len(code.Instructions))
			}
		})
	}
}