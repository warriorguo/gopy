package tests

import (
	"testing"

	"github.com/warriorguo/gopy/pkg/compiler"
	"github.com/warriorguo/gopy/pkg/lexer"
	"github.com/warriorguo/gopy/pkg/object"
	"github.com/warriorguo/gopy/pkg/parser"
	"github.com/warriorguo/gopy/pkg/runtime"
	"github.com/warriorguo/gopy/pkg/vm"
)

func compileAndRunForTest(input string) (object.Object, error) {
	l := lexer.NewLexer(input)
	tokens := l.AllTokens()
	p := parser.NewParser(tokens)

	module, err := p.Parse()
	if err != nil {
		return nil, err
	}

	code, err := compiler.Compile(module)
	if err != nil {
		return nil, err
	}

	machine := vm.NewVM()
	return machine.Run(code)
}

func TestVMForLoopBasic(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		shouldSucceed bool
		description   string
	}{
		{
			name:          "simple_for_range",
			input:         "for i in range(3):\n    print i",
			shouldSucceed: true,
			description:   "Basic for loop with range should execute without errors",
		},
		{
			name:          "for_with_list",
			input:         "for x in [1, 2, 3]:\n    print x",
			shouldSucceed: true,
			description:   "For loop with list should iterate through elements",
		},
		{
			name:          "for_with_string",
			input:         "for c in \"abc\":\n    print c",
			shouldSucceed: true,
			description:   "For loop with string should iterate through characters",
		},
		{
			name:          "empty_range",
			input:         "for i in range(0):\n    print i",
			shouldSucceed: true,
			description:   "For loop with empty range should not execute body",
		},
		{
			name:          "for_with_variable_assignment",
			input:         "total = 0\nfor i in range(3):\n    total = total + i\nprint total",
			shouldSucceed: true,
			description:   "For loop should be able to modify variables",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compileAndRunForTest(test.input)

			if test.shouldSucceed {
				if err != nil {
					t.Errorf("Expected success for %q, but got error: %v", test.input, err)
				} else {
					t.Logf("✓ %s: %q executed successfully", test.description, test.input)
					if result != nil {
						t.Logf("  Result: %s", result.String())
					}
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for %q, but execution succeeded", test.input)
				} else {
					t.Logf("✓ %s: %q failed as expected with error: %v", test.description, test.input, err)
				}
			}
		})
	}
}

func TestVMForLoopWithRange(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue interface{}
	}{
		{
			name:  "sum_range",
			input: "total = 0\nfor i in range(5):\n    total = total + i\ntotal",
			expectedValue: 10, // 0+1+2+3+4 = 10
		},
		{
			name:  "count_iterations",
			input: "count = 0\nfor i in range(3):\n    count = count + 1\ncount",
			expectedValue: 3,
		},
		{
			name:  "last_value",
			input: "last = 0\nfor i in range(1, 4):\n    last = i\nlast",
			expectedValue: 3, // last value should be 3
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compileAndRunForTest(test.input)
			if err != nil {
				t.Errorf("Execution error for %q: %v", test.input, err)
				return
			}

			switch expected := test.expectedValue.(type) {
			case int:
				if intResult, ok := result.(*runtime.PyInt); ok {
					if intResult.Value != expected {
						t.Errorf("Expected %d, got %d for input %q", expected, intResult.Value, test.input)
					}
				} else {
					t.Errorf("Expected PyInt result, got %T for input %q", result, test.input)
				}
			default:
				t.Errorf("Unsupported expected type %T", expected)
			}
		})
	}
}

func TestVMForLoopWithList(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue interface{}
	}{
		{
			name:  "sum_list_elements",
			input: "total = 0\nfor x in [1, 2, 3, 4]:\n    total = total + x\ntotal",
			expectedValue: 10,
		},
		{
			name:  "count_list_elements",
			input: "count = 0\nfor item in [10, 20, 30]:\n    count = count + 1\ncount",
			expectedValue: 3,
		},
		{
			name:  "find_max",
			input: "max_val = 0\nfor x in [3, 7, 2, 9, 1]:\n    if x > max_val:\n        max_val = x\nmax_val",
			expectedValue: 9,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compileAndRunForTest(test.input)
			if err != nil {
				t.Errorf("Execution error for %q: %v", test.input, err)
				return
			}

			switch expected := test.expectedValue.(type) {
			case int:
				if intResult, ok := result.(*runtime.PyInt); ok {
					if intResult.Value != expected {
						t.Errorf("Expected %d, got %d for input %q", expected, intResult.Value, test.input)
					}
				} else {
					t.Errorf("Expected PyInt result, got %T for input %q", result, test.input)
				}
			}
		})
	}
}

func TestVMForLoopVariableScoping(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue int
	}{
		{
			name: "loop_variable_persists",
			input: `i = 100
for i in range(3):
    pass
i`, // i should be 2 (last value from range)
			expectedValue: 2,
		},
		{
			name: "inner_variable_modification",
			input: `x = 0
for i in range(3):
    x = x + i
x`,
			expectedValue: 3, // 0 + 0 + 1 + 2
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compileAndRunForTest(test.input)
			if err != nil {
				t.Errorf("Execution error for %q: %v", test.input, err)
				return
			}

			if intResult, ok := result.(*runtime.PyInt); ok {
				if intResult.Value != test.expectedValue {
					t.Errorf("Expected %d, got %d for input %q", test.expectedValue, intResult.Value, test.input)
				}
			} else {
				t.Errorf("Expected PyInt result, got %T for input %q", result, test.input)
			}
		})
	}
}

func TestVMNestedForLoops(t *testing.T) {
	input := `total = 0
for i in range(2):
    for j in range(3):
        total = total + 1
total`

	result, err := compileAndRunForTest(input)
	if err != nil {
		t.Errorf("Execution error: %v", err)
		return
	}

	if intResult, ok := result.(*runtime.PyInt); ok {
		expected := 6 // 2 * 3 = 6 iterations
		if intResult.Value != expected {
			t.Errorf("Expected %d, got %d for nested for loops", expected, intResult.Value)
		}
	} else {
		t.Errorf("Expected PyInt result, got %T", result)
	}
}

func TestVMForLoopWithConditionals(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue int
	}{
		{
			name: "count_even_numbers",
			input: `count = 0
for i in range(6):
    if i % 2 == 0:
        count = count + 1
count`,
			expectedValue: 3, // 0, 2, 4 are even
		},
		{
			name: "sum_positive_numbers",
			input: `total = 0
for x in [-2, -1, 0, 1, 2]:
    if x > 0:
        total = total + x
total`,
			expectedValue: 3, // 1 + 2 = 3
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compileAndRunForTest(test.input)
			if err != nil {
				t.Errorf("Execution error for %q: %v", test.input, err)
				return
			}

			if intResult, ok := result.(*runtime.PyInt); ok {
				if intResult.Value != test.expectedValue {
					t.Errorf("Expected %d, got %d for input %q", test.expectedValue, intResult.Value, test.input)
				}
			} else {
				t.Errorf("Expected PyInt result, got %T for input %q", result, test.input)
			}
		})
	}
}

func TestVMForLoopEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		description string
	}{
		{
			name:        "empty_list",
			input:       "for x in []:\n    print x\n\"done\"",
			shouldError: false,
			description: "For loop with empty list should not execute body",
		},
		{
			name:        "single_element_list",
			input:       "for x in [42]:\n    print x\n\"done\"",
			shouldError: false,
			description: "For loop with single element should execute once",
		},
		{
			name:        "range_zero",
			input:       "for i in range(0):\n    print \"should not print\"\n\"done\"",
			shouldError: false,
			description: "For loop with range(0) should not execute body",
		},
		{
			name:        "range_one",
			input:       "count = 0\nfor i in range(1):\n    count = count + 1\ncount",
			shouldError: false,
			description: "For loop with range(1) should execute once",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compileAndRunForTest(test.input)

			if test.shouldError {
				if err == nil {
					t.Errorf("Expected error for %q, but execution succeeded", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("Expected success for %q, but got error: %v", test.input, err)
				} else {
					t.Logf("✓ %s: executed successfully", test.description)
					if result != nil {
						t.Logf("  Result: %s", result.String())
					}
				}
			}
		})
	}
}

func TestVMForLoopPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Test with a moderately large range to ensure the VM can handle it
	input := `total = 0
for i in range(100):
    total = total + i
total`

	result, err := compileAndRunForTest(input)
	if err != nil {
		t.Errorf("Performance test failed: %v", err)
		return
	}

	if intResult, ok := result.(*runtime.PyInt); ok {
		// Sum of 0 to 99 = 99 * 100 / 2 = 4950
		expected := 4950
		if intResult.Value != expected {
			t.Errorf("Expected %d, got %d for performance test", expected, intResult.Value)
		} else {
			t.Logf("✓ Performance test passed: computed sum of 0-99 = %d", intResult.Value)
		}
	} else {
		t.Errorf("Expected PyInt result, got %T", result)
	}
}