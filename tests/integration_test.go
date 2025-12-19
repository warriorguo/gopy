package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestEndToEndCompilation tests the complete compilation pipeline
func TestEndToEndCompilation(t *testing.T) {
	// Build the tools first
	err := buildTools()
	if err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}

	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "hello_world",
			source:   `print "Hello, World!"`,
			expected: `"Hello, World!"`,
		},
		{
			name:     "arithmetic",
			source:   `print 2 + 3 * 4`,
			expected: `14`,
		},
		{
			name: "variables",
			source: `x = 42
print "x =", x`,
			expected: `"x =" 42`,
		},
		{
			name: "for_loop",
			source: `for i in range(3):
    print "i =", i`,
			expected: `"i =" 0
"i =" 1
"i =" 2`,
		},
		{
			name: "if_statement",
			source: `x = 10
if x > 5:
    print "big"`,
			expected: `"big"`,
		},
		{
			name: "function_definition",
			source: `def double(x):
    return x * 2

result = double(5)
print result`,
			expected: `10`,
		},
		{
			name:     "list_operations",
			source:   `print len([1, 2, 3, 4, 5])`,
			expected: `5`,
		},
		{
			name: "string_concatenation",
			source: `name = "GoPy"
print "Hello, " + name + "!"`,
			expected: `"Hello, GoPy!"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := compileAndRunSource(test.source)
			if err != nil {
				t.Errorf("Failed to compile and run %s: %v", test.name, err)
				return
			}

			output = strings.TrimSpace(output)
			expected := strings.TrimSpace(test.expected)

			if output != expected {
				t.Errorf("For test %s:\nExpected: %q\nGot:      %q", test.name, expected, output)
			}
		})
	}
}

// TestExampleFiles tests the provided example files
func TestExampleFiles(t *testing.T) {
	err := buildTools()
	if err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}

	exampleDir := "../examples"
	examples := []string{"hello.py", "fibonacci.py", "calculator.py"}

	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			examplePath := filepath.Join(exampleDir, example)
			
			// Check if file exists
			if _, err := os.Stat(examplePath); os.IsNotExist(err) {
				t.Skipf("Example file %s not found", example)
				return
			}

			// Try to compile the example
			output, err := compileAndRunFile(examplePath)
			if err != nil {
				t.Errorf("Failed to compile and run example %s: %v", example, err)
				return
			}

			// For now, just check that it produces some output without errors
			if output == "" {
				t.Logf("Example %s produced no output (this may be expected)", example)
			} else {
				t.Logf("Example %s output: %s", example, strings.TrimSpace(output))
			}
		})
	}
}

// TestCompilerFlags tests various compiler flags and options
func TestCompilerFlags(t *testing.T) {
	err := buildTools()
	if err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}

	source := `print "test"`
	sourceFile := createTempFile(t, source)
	defer os.Remove(sourceFile)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "verbose_flag",
			args: []string{"-v", sourceFile},
		},
		{
			name: "disassemble_flag",
			args: []string{"-d", sourceFile},
		},
		{
			name: "custom_output",
			args: []string{"-o", "custom.pyc", sourceFile},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := exec.Command("./py2c", test.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("Command failed with args %v: %v\nOutput: %s", test.args, err, output)
				return
			}

			// Clean up custom output file if created
			if test.name == "custom_output" {
				os.Remove("custom.pyc")
			}

			t.Logf("Command with args %v succeeded. Output: %s", test.args, string(output))
		})
	}
}

// TestErrorHandling tests how the system handles various error conditions
func TestErrorHandling(t *testing.T) {
	err := buildTools()
	if err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}

	tests := []struct {
		name   string
		source string
		isError bool
	}{
		{
			name:   "syntax_error",
			source: `if x > 0 print "missing colon"`,
			isError: true,
		},
		{
			name:   "indentation_error",
			source: `if x > 0:
print "bad indentation"`,
			isError: true,
		},
		{
			name:   "valid_syntax",
			source: `print "this should work"`,
			isError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := compileAndRunSource(test.source)
			
			if test.isError && err == nil {
				t.Errorf("Expected error for %s, but compilation succeeded", test.name)
			}
			
			if !test.isError && err != nil {
				t.Errorf("Unexpected error for %s: %v", test.name, err)
			}
		})
	}
}

// TestPerformance runs basic performance tests
func TestPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	err := buildTools()
	if err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}

	// Test compilation performance with a larger program
	source := `
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

for i in range(10):
    result = fibonacci(i)
    print "fib(", i, ") =", result
`

	start := time.Now()
	_, err = compileAndRunSource(source)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Performance test failed: %v", err)
		return
	}

	t.Logf("Large program compilation and execution took: %v", duration)
	
	// Arbitrary performance threshold (adjust as needed)
	if duration > 10000000000 { // 10 seconds
		t.Errorf("Performance test took too long: %v", duration)
	}
}

// Helper functions

func buildTools() error {
	// Build py2c
	cmd := exec.Command("go", "build", "-o", "py2c", "./cmd/py2c")
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Build py2vm
	cmd = exec.Command("go", "build", "-o", "py2vm", "./cmd/py2vm")
	return cmd.Run()
}

func compileAndRunSource(source string) (string, error) {
	sourceFile := createTempFile(nil, source)
	defer os.Remove(sourceFile)
	
	return compileAndRunFile(sourceFile)
}

func compileAndRunFile(sourceFile string) (string, error) {
	// Compile
	bytecodeFile := strings.TrimSuffix(sourceFile, filepath.Ext(sourceFile)) + ".pyc"
	defer os.Remove(bytecodeFile)

	cmd := exec.Command("./py2c", "-o", bytecodeFile, sourceFile)
	compileOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Run
	cmd = exec.Command("./py2vm", bytecodeFile)
	runOutput, err := cmd.Output()
	if err != nil {
		return "", err
	}

	_ = compileOutput // Suppress unused variable warning
	return string(runOutput), nil
}

func createTempFile(t *testing.T, content string) string {
	tempFile, err := os.CreateTemp("", "gopy_test_*.py")
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		panic(err)
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(content)
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to write temp file: %v", err)
		}
		panic(err)
	}

	return tempFile.Name()
}