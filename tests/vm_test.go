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

func compileAndRun(input string) (object.Object, error) {
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

func TestVMArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"1 + 2", 3},
		{"5 - 3", 2},
		{"2 * 4", 8},
		{"8 / 2", 4},
		{"10 % 3", 1},
		{"2 + 3 * 4", 14},        // precedence test
		{"(2 + 3) * 4", 20},     // parentheses test
		{"-5", -5},               // unary minus
		{"+10", 10},              // unary plus
	}

	for _, test := range tests {
		result, err := compileAndRun(test.input)
		if err != nil {
			t.Errorf("Execution error for %q: %v", test.input, err)
			continue
		}

		switch expected := test.expected.(type) {
		case int:
			if intObj, ok := result.(*runtime.PyInt); ok {
				if intObj.Value != expected {
					t.Errorf("Expected %d, got %d for input %q", expected, intObj.Value, test.input)
				}
			} else {
				t.Errorf("Expected PyInt, got %T for input %q", result, test.input)
			}
		}
	}
}

func TestVMFloatArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1.5 + 2.5", 4.0},
		{"5.0 - 2.5", 2.5},
		{"2.0 * 3.0", 6.0},
		{"8.0 / 2.0", 4.0},
	}

	for _, test := range tests {
		result, err := compileAndRun(test.input)
		if err != nil {
			t.Errorf("Execution error for %q: %v", test.input, err)
			continue
		}

		if floatObj, ok := result.(*runtime.PyFloat); ok {
			if floatObj.Value != test.expected {
				t.Errorf("Expected %f, got %f for input %q", test.expected, floatObj.Value, test.input)
			}
		} else {
			t.Errorf("Expected PyFloat, got %T for input %q", result, test.input)
		}
	}
}

func TestVMComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5 == 5", true},
		{"5 == 3", false},
		{"5 != 3", true},
		{"5 != 5", false},
		{"5 > 3", true},
		{"3 > 5", false},
		{"5 >= 5", true},
		{"5 >= 3", true},
		{"3 >= 5", false},
		{"3 < 5", true},
		{"5 < 3", false},
		{"5 <= 5", true},
		{"3 <= 5", true},
		{"5 <= 3", false},
	}

	for _, test := range tests {
		result, err := compileAndRun(test.input)
		if err != nil {
			t.Errorf("Execution error for %q: %v", test.input, err)
			continue
		}

		if boolObj, ok := result.(*runtime.PyBool); ok {
			if boolObj.Value != test.expected {
				t.Errorf("Expected %t, got %t for input %q", test.expected, boolObj.Value, test.input)
			}
		} else {
			t.Errorf("Expected PyBool, got %T for input %q", result, test.input)
		}
	}
}

func TestVMStringOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello" + " world"`, "hello world"},
		{`"test"`, "test"},
	}

	for _, test := range tests {
		result, err := compileAndRun(test.input)
		if err != nil {
			t.Errorf("Execution error for %q: %v", test.input, err)
			continue
		}

		if strObj, ok := result.(*runtime.PyString); ok {
			if strObj.Value != test.expected {
				t.Errorf("Expected %q, got %q for input %q", test.expected, strObj.Value, test.input)
			}
		} else {
			t.Errorf("Expected PyString, got %T for input %q", result, test.input)
		}
	}
}

func TestVMVariables(t *testing.T) {
	input := `x = 42
y = x + 8
y`

	result, err := compileAndRun(input)
	if err != nil {
		t.Errorf("Execution error: %v", err)
		return
	}

	if intObj, ok := result.(*runtime.PyInt); ok {
		if intObj.Value != 50 {
			t.Errorf("Expected 50, got %d", intObj.Value)
		}
	} else {
		t.Errorf("Expected PyInt, got %T", result)
	}
}

func TestVMBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, 5},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`str(42)`, "42"},
		{`type(42)`, "int"},
		{`type("hello")`, "str"},
	}

	for _, test := range tests {
		result, err := compileAndRun(test.input)
		if err != nil {
			t.Errorf("Execution error for %q: %v", test.input, err)
			continue
		}

		switch expected := test.expected.(type) {
		case int:
			if intObj, ok := result.(*runtime.PyInt); ok {
				if intObj.Value != expected {
					t.Errorf("Expected %d, got %d for input %q", expected, intObj.Value, test.input)
				}
			} else {
				t.Errorf("Expected PyInt, got %T for input %q", result, test.input)
			}
		case string:
			if strObj, ok := result.(*runtime.PyString); ok {
				if strObj.Value != expected {
					t.Errorf("Expected %q, got %q for input %q", expected, strObj.Value, test.input)
				}
			} else {
				t.Errorf("Expected PyString, got %T for input %q", result, test.input)
			}
		}
	}
}

func TestVMLists(t *testing.T) {
	tests := []struct {
		input    string
		expected int // length
	}{
		{"[]", 0},
		{"[1, 2, 3]", 3},
		{"[1, 2, 3, 4, 5]", 5},
	}

	for _, test := range tests {
		result, err := compileAndRun(test.input)
		if err != nil {
			t.Errorf("Execution error for %q: %v", test.input, err)
			continue
		}

		if listObj, ok := result.(*runtime.PyList); ok {
			if len(listObj.Elements) != test.expected {
				t.Errorf("Expected list length %d, got %d for input %q", test.expected, len(listObj.Elements), test.input)
			}
		} else {
			t.Errorf("Expected PyList, got %T for input %q", result, test.input)
		}
	}
}

func TestVMListIndexing(t *testing.T) {
	input := `arr = [10, 20, 30]
arr[1]`

	result, err := compileAndRun(input)
	if err != nil {
		t.Errorf("Execution error: %v", err)
		return
	}

	if intObj, ok := result.(*runtime.PyInt); ok {
		if intObj.Value != 20 {
			t.Errorf("Expected 20, got %d", intObj.Value)
		}
	} else {
		t.Errorf("Expected PyInt, got %T", result)
	}
}

func TestVMDictionaries(t *testing.T) {
	input := `d = {1: "one", 2: "two"}
d[1]`

	result, err := compileAndRun(input)
	if err != nil {
		t.Errorf("Execution error: %v", err)
		return
	}

	if strObj, ok := result.(*runtime.PyString); ok {
		if strObj.Value != "one" {
			t.Errorf("Expected 'one', got %q", strObj.Value)
		}
	} else {
		t.Errorf("Expected PyString, got %T", result)
	}
}

func TestVMBooleanLogic(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"True and True", true},
		{"True and False", false},
		{"False or True", true},
		{"False or False", false},
		{"not True", false},
		{"not False", true},
	}

	for _, test := range tests {
		result, err := compileAndRun(test.input)
		if err != nil {
			t.Errorf("Execution error for %q: %v", test.input, err)
			continue
		}

		if boolObj, ok := result.(*runtime.PyBool); ok {
			if boolObj.Value != test.expected {
				t.Errorf("Expected %t, got %t for input %q", test.expected, boolObj.Value, test.input)
			}
		} else {
			t.Errorf("Expected PyBool, got %T for input %q", result, test.input)
		}
	}
}