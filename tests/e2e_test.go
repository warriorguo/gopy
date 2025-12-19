package tests

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestE2EPythonFeatures tests comprehensive Python language features
func TestE2EPythonFeatures(t *testing.T) {
	if err := setupE2E(); err != nil {
		t.Fatalf("E2E setup failed: %v", err)
	}
	defer cleanupE2E()

	t.Run("DataTypes", testDataTypes)
	t.Run("ControlFlow", testControlFlow)
	t.Run("Functions", testFunctions)
	t.Run("DataStructures", testDataStructures)
	t.Run("Operators", testOperators)
	t.Run("BuiltinFunctions", testBuiltinFunctions)
	t.Run("ErrorCases", testErrorCases)
}

func testDataTypes(t *testing.T) {
	tests := map[string]struct {
		source   string
		expected string
	}{
		"integers": {
			source: `print 42
print -10
print 0`,
			expected: `42
-10
0`,
		},
		"floats": {
			source: `print 3.14
print -2.5
print 0.0`,
			expected: `3.14
-2.5
0`,
		},
		"strings": {
			source: `print "hello world"
print 'single quotes'
print "escaped \" quote"`,
			expected: `"hello world"
"single quotes"
"escaped \" quote"`,
		},
		"booleans": {
			source: `print True
print False
print not True
print not False`,
			expected: `True
False
False
True`,
		},
		"none": {
			source: `x = None
print x`,
			expected: `None`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := runPythonCode(test.source)
			if err != nil {
				t.Errorf("Failed to run %s: %v", name, err)
				return
			}
			assertOutput(t, test.expected, output)
		})
	}
}

func testControlFlow(t *testing.T) {
	tests := map[string]struct {
		source   string
		expected string
	}{
		"if_else": {
			source: `x = 10
if x > 5:
    print "greater than 5"
else:
    print "less than or equal to 5"

y = 3
if y > 5:
    print "y greater than 5"
else:
    print "y less than or equal to 5"`,
			expected: `"greater than 5"
"y less than or equal to 5"`,
		},
		"elif": {
			source: `score = 85
if score >= 90:
    print "A"
elif score >= 80:
    print "B"
elif score >= 70:
    print "C"
else:
    print "F"`,
			expected: `"B"`,
		},
		"for_loop": {
			source: `for i in range(5):
    print i`,
			expected: `0
1
2
3
4`,
		},
		"while_loop": {
			source: `i = 0
while i < 3:
    print i
    i = i + 1`,
			expected: `0
1
2`,
		},
		"nested_loops": {
			source: `for i in range(2):
    for j in range(2):
        print i, j`,
			expected: `0 0
0 1
1 0
1 1`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := runPythonCode(test.source)
			if err != nil {
				t.Errorf("Failed to run %s: %v", name, err)
				return
			}
			assertOutput(t, test.expected, output)
		})
	}
}

func testFunctions(t *testing.T) {
	tests := map[string]struct {
		source   string
		expected string
	}{
		"simple_function": {
			source: `def greet(name):
    return "Hello, " + name

result = greet("World")
print result`,
			expected: `"Hello, World"`,
		},
		"function_with_multiple_params": {
			source: `def add(a, b):
    return a + b

result = add(5, 3)
print result`,
			expected: `8`,
		},
		"function_calling_function": {
			source: `def square(x):
    return x * x

def sum_of_squares(a, b):
    return square(a) + square(b)

result = sum_of_squares(3, 4)
print result`,
			expected: `25`,
		},
		"recursive_function": {
			source: `def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n - 1)

result = factorial(5)
print result`,
			expected: `120`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := runPythonCode(test.source)
			if err != nil {
				t.Errorf("Failed to run %s: %v", name, err)
				return
			}
			assertOutput(t, test.expected, output)
		})
	}
}

func testDataStructures(t *testing.T) {
	tests := map[string]struct {
		source   string
		expected string
	}{
		"lists": {
			source: `numbers = [1, 2, 3, 4, 5]
print len(numbers)
print numbers[0]
print numbers[2]`,
			expected: `5
1
3`,
		},
		"list_operations": {
			source: `items = []
items = [1, 2]
items[1] = 10
print items[1]`,
			expected: `10`,
		},
		"dictionaries": {
			source: `person = {"name": "Alice", "age": 30}
print person["name"]
print person["age"]`,
			expected: `"Alice"
30`,
		},
		"nested_structures": {
			source: `matrix = [[1, 2], [3, 4]]
print matrix[0][1]
print matrix[1][0]`,
			expected: `2
3`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := runPythonCode(test.source)
			if err != nil {
				t.Errorf("Failed to run %s: %v", name, err)
				return
			}
			assertOutput(t, test.expected, output)
		})
	}
}

func testOperators(t *testing.T) {
	tests := map[string]struct {
		source   string
		expected string
	}{
		"arithmetic": {
			source: `print 10 + 5
print 10 - 5
print 10 * 5
print 10 / 5
print 10 % 3`,
			expected: `15
5
50
2
1`,
		},
		"comparison": {
			source: `print 5 == 5
print 5 != 3
print 5 > 3
print 3 < 5
print 5 >= 5
print 3 <= 5`,
			expected: `True
True
True
True
True
True`,
		},
		"logical": {
			source: `print True and True
print True and False
print False or True
print not True
print not False`,
			expected: `True
False
True
False
True`,
		},
		"operator_precedence": {
			source: `print 2 + 3 * 4
print (2 + 3) * 4
print 2 * 3 + 4
print 2 * (3 + 4)`,
			expected: `14
20
10
14`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := runPythonCode(test.source)
			if err != nil {
				t.Errorf("Failed to run %s: %v", name, err)
				return
			}
			assertOutput(t, test.expected, output)
		})
	}
}

func testBuiltinFunctions(t *testing.T) {
	tests := map[string]struct {
		source   string
		expected string
	}{
		"len": {
			source: `print len("hello")
print len([1, 2, 3, 4])
print len({1: "one", 2: "two"})`,
			expected: `5
4
2`,
		},
		"type": {
			source: `print type(42)
print type(3.14)
print type("hello")
print type([1, 2, 3])
print type(True)`,
			expected: `"int"
"float"
"str"
"list"
"bool"`,
		},
		"str": {
			source: `print str(42)
print str(3.14)
print str(True)`,
			expected: `"42"
"3.14"
"True"`,
		},
		"range": {
			source: `numbers = range(5)
print len(numbers)
for i in range(3):
    print i`,
			expected: `5
0
1
2`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := runPythonCode(test.source)
			if err != nil {
				t.Errorf("Failed to run %s: %v", name, err)
				return
			}
			assertOutput(t, test.expected, output)
		})
	}
}

func testErrorCases(t *testing.T) {
	tests := map[string]string{
		"syntax_error":      `if x > 0 print "missing colon"`,
		"indentation_error": "if True:\nprint \"bad indentation\"",
		"undefined_variable": `print undefined_var`,
		"type_error":        `print "string" + 42`,
		"index_error":       `arr = [1, 2, 3]\nprint arr[10]`,
	}

	for name, source := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := runPythonCode(source)
			if err == nil {
				t.Errorf("Expected error for %s, but execution succeeded", name)
			} else {
				t.Logf("Expected error for %s: %v", name, err)
			}
		})
	}
}

// Helper functions for E2E testing

func setupE2E() error {
	// Build the tools
	cmd := exec.Command("go", "build", "-o", "py2c", "./cmd/py2c")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("go", "build", "-o", "py2vm", "./cmd/py2vm")
	return cmd.Run()
}

func cleanupE2E() {
	os.Remove("py2c")
	os.Remove("py2vm")
	
	// Remove any leftover .pyc files
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(path, ".pyc") && strings.Contains(path, "test_") {
			os.Remove(path)
		}
		return nil
	})
}

func runPythonCode(source string) (string, error) {
	// Create temporary source file
	sourceFile, err := os.CreateTemp("", "test_*.py")
	if err != nil {
		return "", err
	}
	defer os.Remove(sourceFile.Name())
	defer sourceFile.Close()

	_, err = sourceFile.WriteString(source)
	if err != nil {
		return "", err
	}
	sourceFile.Close()

	// Compile
	bytecodeFile := strings.TrimSuffix(sourceFile.Name(), ".py") + ".pyc"
	defer os.Remove(bytecodeFile)

	compileCmd := exec.Command("./py2c", "-o", bytecodeFile, sourceFile.Name())
	compileCmd.Stderr = os.Stderr
	if err := compileCmd.Run(); err != nil {
		return "", err
	}

	// Run with timeout
	runCmd := exec.Command("./py2vm", bytecodeFile)
	
	// Set a timeout to prevent infinite loops
	done := make(chan error, 1)
	var output []byte
	
	go func() {
		var err error
		output, err = runCmd.Output()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", err
		}
		return string(output), nil
	case <-time.After(5 * time.Second):
		runCmd.Process.Kill()
		return "", &timeoutError{}
	}
}

func assertOutput(t *testing.T, expected, actual string) {
	expected = strings.TrimSpace(expected)
	actual = strings.TrimSpace(actual)
	
	if expected != actual {
		t.Errorf("Output mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
	}
}

type timeoutError struct{}

func (e *timeoutError) Error() string {
	return "execution timeout"
}

// TestE2EBenchmarks provides basic performance benchmarking
func TestE2EBenchmarks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping benchmarks in short mode")
	}

	if err := setupE2E(); err != nil {
		t.Fatalf("E2E setup failed: %v", err)
	}
	defer cleanupE2E()

	benchmarks := map[string]string{
		"fibonacci": `def fib(n):
    if n <= 1:
        return n
    return fib(n-1) + fib(n-2)
print fib(10)`,
		
		"nested_loops": `total = 0
for i in range(100):
    for j in range(100):
        total = total + 1
print total`,
		
		"list_operations": `numbers = []
for i in range(1000):
    numbers = numbers + [i]
print len(numbers)`,
	}

	for name, source := range benchmarks {
		t.Run(name, func(t *testing.T) {
			start := time.Now()
			_, err := runPythonCode(source)
			duration := time.Since(start)
			
			if err != nil {
				t.Errorf("Benchmark %s failed: %v", name, err)
				return
			}
			
			t.Logf("Benchmark %s completed in %v", name, duration)
			
			// Set reasonable performance thresholds
			maxDuration := 10 * time.Second
			if duration > maxDuration {
				t.Errorf("Benchmark %s took too long: %v (max: %v)", name, duration, maxDuration)
			}
		})
	}
}