# TestUtils - Go Testing Utilities Library

A comprehensive testing utilities library for Go that provides colorized output, easy debugging, benchmarking, and random test data generation.

## Features

- **Colorized Test Output**: Clear visual distinction between success, failure, and debug information
- **Simple Test Framework**: Easy-to-use testing function with customizable behavior
- **Benchmark Tool**: Measure and compare execution times of multiple code blocks
- **Random Data Generation**: Generate random test data for any type
- **Stack Trace Utilities**: Access and format stack traces for better debugging
- **Panic Handler**: Simplified error handling with panic option

## Installation

```bash
go get github.com/HashemJaafar7/testutils
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/HashemJaafar7/testutils"
)

func main() {
	// Simple test
	actual := 42
	expected := 42
	testutils.Test(false, true, true, 8, "v", actual, expected)

	// Debug output
	value := "debug this"
	testutils.Debug("v", value)

	// Benchmark comparison
	testutils.Benchmark(1000,
		func() { /* First operation */ },
		func() { /* Second operation */ },
	)
}
```

## API Reference

### Test Function

```go
func Test[t any](isPanic, print, isEqual bool, line uint16, format string, actual, expected t)
```

Parameters:

- `isPanic`: Exit program if test fails
- `print`: Print results regardless of pass/fail
- `isEqual`: Expected equality relationship
- `line`: Stack trace line number (must be > 8)
- `format`: Printf format string for output
- `actual`: Value being tested
- `expected`: Value to test against

Example:

```go
// Test for equality with output
testutils.Test(false, true, true, 8, "%v", result, expectedResult)

// Test for inequality with panic on failure
testutils.Test(true, false, false, 8, "%v", value1, value2)
```

### Debug Function

```go
func Debug(format string, a any)
```

Prints formatted debug information with stack trace and variable name.

Example:

```go
complexValue := calculateSomething()
testutils.Debug("+v", complexValue)
```

### Benchmark Function

```go
func Benchmark(loops uint, codesBlocks ...func())
```

Measures execution time of multiple code blocks.

Example:

```go
testutils.Benchmark(1000,
    func() { method1() },
    func() { method2() },
    func() { method3() },
)
```

### Random Data Generation

```go
func Rand[t any]() t
```

Generates random values of any type using fuzzing.

Example:

```go
randomInt := testutils.Rand[int]()
randomString := testutils.Rand[string]()
randomStruct := testutils.Rand[MyStruct]()
```

### Stack Trace Utility

```go
func Stack(line uint16) string
```

Retrieves specific line from stack trace.

Example:

```go
stackLine := testutils.Stack(6)
```

### Panic If Error

```go
func PanicIfErr(err error)
```

Panics if error is not nil.

Example:

```go
testutils.PanicIfErr(err)
```

## Color Constants

The library provides color constants for output formatting:

```go
const (
    ColorReset   = "\033[0m"
    ColorRed     = "\033[31m" // Failure
    ColorGreen   = "\033[32m" // Success
    ColorYellow  = "\033[33m" // Actual value
    ColorBlue    = "\033[34m" // Expected value
    ColorMagenta = "\033[35m" // Debug
    ColorCyan    = "\033[36m" // Benchmark
)
```

## Best Practices

1. **Test Line Numbers**

   - Use appropriate line numbers (> 8) for stack traces
   - Keep line numbers consistent within test suites

2. **Benchmarking**

   - Use sufficient loop counts for accurate measurements
   - Keep benchmark functions focused and isolated
   - Compare similar operations in the same benchmark

3. **Debug Output**

   - Use descriptive format strings
   - Debug complex structures with `+v`
   - Keep debug statements organized

4. **Random Testing**
   - Use Rand() with appropriate types
   - Validate generated data meets requirements
   - Consider edge cases

## Example Use Cases

### Complex Testing Scenario

```go
func TestComplexOperation() {
    input := testutils.Rand[MyStruct]()
    result := ComplexOperation(input)

    // Test with output and panic on failure
    testutils.Test(true, true, true, 8, "+v",
        result,
        expectedResult,
    )
}
```

### Performance Comparison

```go
func BenchmarkAlgorithms() {
    data := generateTestData()

    testutils.Benchmark(1000,
        func() { algorithm1(data) },
        func() { algorithm2(data) },
        func() { algorithm3(data) },
    )
}
```

### Debugging Complex Structures

```go
func ProcessData() {
    result := complexCalculation()
    testutils.Debug("+v", result)

    // Continue processing...
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)
