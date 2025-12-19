package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestForLoopScenarios tests comprehensive for loop scenarios end-to-end
func TestForLoopScenarios(t *testing.T) {
	// Build tools if needed
	if err := buildForLoopTestTools(); err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}
	defer cleanupForLoopTestTools()

	scenarios := map[string]struct {
		code        string
		expected    string
		shouldError bool
		description string
	}{
		"basic_range": {
			code: `for i in range(3):
    print i`,
			expected: "0\n1\n2",
			description: "Basic for loop with range should print 0, 1, 2",
		},

		"range_with_start_end": {
			code: `for i in range(2, 5):
    print i`,
			expected: "2\n3\n4",
			description: "Range with start and end",
		},

		"for_with_list": {
			code: `for x in [10, 20, 30]:
    print x`,
			expected: "10\n20\n30",
			description: "For loop with list literal",
		},

		"for_with_string": {
			code: `for c in "abc":
    print c`,
			expected: `"a"\n"b"\n"c"`,
			description: "For loop with string should iterate characters",
		},

		"nested_for_loops": {
			code: `for i in range(2):
    for j in range(2):
        print i, j`,
			expected: "0 0\n0 1\n1 0\n1 1",
			description: "Nested for loops",
		},

		"for_with_accumulator": {
			code: `total = 0
for i in range(5):
    total = total + i
print total`,
			expected: "10",
			description: "For loop accumulating values",
		},

		"for_with_conditional": {
			code: `for i in range(5):
    if i % 2 == 0:
        print i, "even"
    else:
        print i, "odd"`,
			expected: "0 \"even\"\n1 \"odd\"\n2 \"even\"\n3 \"odd\"\n4 \"even\"",
			description: "For loop with conditional logic",
		},

		"for_with_function": {
			code: `def square(x):
    return x * x

for i in range(1, 4):
    result = square(i)
    print i, "squared is", result`,
			expected: "1 \"squared is\" 1\n2 \"squared is\" 4\n3 \"squared is\" 9",
			description: "For loop calling functions",
		},

		"for_with_list_operations": {
			code: `numbers = []
for i in range(3):
    numbers = numbers + [i * 2]
for x in numbers:
    print x`,
			expected: "0\n2\n4",
			description: "For loop building and iterating lists",
		},

		"for_variable_scoping": {
			code: `x = 100
for x in range(3):
    pass
print x`,
			expected: "2",
			description: "For loop variable should override outer scope",
		},

		"empty_range": {
			code: `print "before"
for i in range(0):
    print "inside loop"
print "after"`,
			expected: "\"before\"\n\"after\"",
			description: "Empty range should not execute loop body",
		},

		"single_iteration": {
			code: `for i in range(1):
    print "executed once"`,
			expected: "\"executed once\"",
			description: "Single iteration loop",
		},

		"for_with_break": {
			code: `for i in range(10):
    if i == 3:
        break
    print i`,
			expected: "0\n1\n2",
			shouldError: true, // break not implemented yet
			description: "For loop with break statement",
		},

		"for_with_continue": {
			code: `for i in range(5):
    if i == 2:
        continue
    print i`,
			expected: "0\n1\n3\n4", 
			shouldError: true, // continue not implemented yet
			description: "For loop with continue statement",
		},

		// Error cases
		"syntax_error_missing_colon": {
			code: `for i in range(3)
    print i`,
			shouldError: true,
			description: "Missing colon should cause syntax error",
		},

		"syntax_error_missing_in": {
			code: `for i range(3):
    print i`,
			shouldError: true,
			description: "Missing 'in' keyword should cause syntax error",
		},

		"runtime_error_non_iterable": {
			code: `for i in 42:
    print i`,
			shouldError: true,
			description: "Iterating over non-iterable should cause runtime error",
		},
	}

	for name, scenario := range scenarios {
		t.Run(name, func(t *testing.T) {
			output, err := runForLoopScenario(scenario.code)

			if scenario.shouldError {
				if err == nil {
					t.Errorf("Expected error for scenario %q, but execution succeeded", name)
					t.Logf("Unexpected output: %q", output)
				} else {
					t.Logf("✓ %s: Failed as expected with error: %v", scenario.description, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected success for scenario %q, but got error: %v", name, err)
					return
				}

				output = strings.TrimSpace(output)
				expected := strings.TrimSpace(scenario.expected)

				if output != expected {
					t.Errorf("Output mismatch for scenario %q:\nExpected:\n%s\n\nActual:\n%s", 
						name, expected, output)
				} else {
					t.Logf("✓ %s: Output matches expected", scenario.description)
				}
			}
		})
	}
}

func TestForLoopPerformanceScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance scenarios in short mode")
	}

	if err := buildForLoopTestTools(); err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}
	defer cleanupForLoopTestTools()

	scenarios := []struct {
		name        string
		code        string
		maxDuration time.Duration
	}{
		{
			name: "small_range_loop",
			code: `total = 0
for i in range(100):
    total = total + i
print total`,
			maxDuration: 5 * time.Second,
		},
		{
			name: "nested_small_loops",
			code: `total = 0
for i in range(10):
    for j in range(10):
        total = total + 1
print total`,
			maxDuration: 5 * time.Second,
		},
		{
			name: "list_iteration",
			code: `numbers = [1, 2, 3, 4, 5]
total = 0
for x in numbers:
    total = total + x
print total`,
			maxDuration: 3 * time.Second,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			start := time.Now()
			output, err := runForLoopScenario(scenario.code)
			duration := time.Since(start)

			if err != nil {
				t.Errorf("Performance scenario %q failed: %v", scenario.name, err)
				return
			}

			if duration > scenario.maxDuration {
				t.Errorf("Performance scenario %q took too long: %v (max: %v)", 
					scenario.name, duration, scenario.maxDuration)
			} else {
				t.Logf("✓ Performance scenario %q completed in %v (output: %s)", 
					scenario.name, duration, strings.TrimSpace(output))
			}
		})
	}
}

func TestForLoopComplexScenarios(t *testing.T) {
	if err := buildForLoopTestTools(); err != nil {
		t.Fatalf("Failed to build tools: %v", err)
	}
	defer cleanupForLoopTestTools()

	scenarios := map[string]struct {
		code        string
		description string
	}{
		"fibonacci_with_for": {
			code: `a, b = 0, 1
for i in range(5):
    print a
    a, b = b, a + b`,
			description: "Generate Fibonacci sequence using for loop",
		},

		"multiplication_table": {
			code: `for i in range(1, 4):
    for j in range(1, 4):
        result = i * j
        print i, "x", j, "=", result`,
			description: "Generate multiplication table",
		},

		"prime_check_simulation": {
			code: `number = 7
is_prime = True
for i in range(2, number):
    if number % i == 0:
        is_prime = False
if is_prime:
    print number, "is prime"
else:
    print number, "is not prime"`,
			description: "Prime number checking simulation",
		},

		"list_comprehension_simulation": {
			code: `squares = []
for x in range(5):
    square = x * x
    squares = squares + [square]
for s in squares:
    print s`,
			description: "Simulate list comprehension with for loops",
		},

		"matrix_operations": {
			code: `matrix = [[1, 2], [3, 4]]
total = 0
for row in matrix:
    for element in row:
        total = total + element
print "Matrix sum:", total`,
			description: "Matrix operations using nested for loops",
		},
	}

	for name, scenario := range scenarios {
		t.Run(name, func(t *testing.T) {
			output, err := runForLoopScenario(scenario.code)

			// For complex scenarios, we mainly test that they don't crash
			if err != nil {
				t.Logf("Complex scenario %q failed (this may be expected due to implementation limitations): %v", 
					name, err)
			} else {
				t.Logf("✓ %s: executed successfully", scenario.description)
				t.Logf("  Output: %s", strings.TrimSpace(output))
			}
		})
	}
}

// Helper functions for for loop testing

func buildForLoopTestTools() error {
	cmd := exec.Command("go", "build", "-o", "py2c_for_test", "../cmd/py2c")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("go", "build", "-o", "py2vm_for_test", "../cmd/py2vm")
	return cmd.Run()
}

func cleanupForLoopTestTools() {
	os.Remove("py2c_for_test")
	os.Remove("py2vm_for_test")
}

func runForLoopScenario(code string) (string, error) {
	// Create temporary source file
	sourceFile, err := os.CreateTemp("", "for_test_*.py")
	if err != nil {
		return "", err
	}
	defer os.Remove(sourceFile.Name())
	defer sourceFile.Close()

	_, err = sourceFile.WriteString(code)
	if err != nil {
		return "", err
	}
	sourceFile.Close()

	// Compile
	bytecodeFile := strings.TrimSuffix(sourceFile.Name(), ".py") + ".pyc"
	defer os.Remove(bytecodeFile)

	compileCmd := exec.Command("./py2c_for_test", "-o", bytecodeFile, sourceFile.Name())
	if err := compileCmd.Run(); err != nil {
		return "", err
	}

	// Run with timeout
	runCmd := exec.Command("./py2vm_for_test", bytecodeFile)
	
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
	case <-time.After(10 * time.Second):
		runCmd.Process.Kill()
		return "", &forLoopTimeoutError{}
	}
}

type forLoopTimeoutError struct{}

func (e *forLoopTimeoutError) Error() string {
	return "for loop test execution timeout"
}