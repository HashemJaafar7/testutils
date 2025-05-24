package testutils

import (
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	fuzz "github.com/google/gofuzz"
)

const (
	ColorReset = "\033[0m"

	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m" //for fail
	ColorGreen   = "\033[32m" //for success
	ColorYellow  = "\033[33m" //for actual
	ColorBlue    = "\033[34m" //for expected
	ColorMagenta = "\033[35m" //for Debug
	ColorCyan    = "\033[36m" //for Benchmark
	ColorWhite   = "\033[37m"
)

// PanicIfErr checks if the provided error is not nil, and if so, it panics with the error.
// This function is useful for quickly handling unexpected errors in situations where
// error recovery is not required or desired.
//
// Parameters:
//   - err: The error to check. If it is nil, the function does nothing.
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Rand generates a random value of any type `t` using fuzzing.
// It initializes a variable of type `t`, applies fuzzing to populate it with random data,
// and then returns the result. A small delay is introduced to ensure randomness.
//
// Type Parameters:
//   - t: The type of the value to be generated.
//
// Returns:
//
//	A randomly generated value of type `t`.
func Rand[t any]() t {
	var result t
	fuzz.New().Fuzz(&result)
	time.Sleep(1 * time.Microsecond)
	return result
}

// Stack retrieves a specific line from the current goroutine's stack trace.
// It returns the specified line as a string after removing leading/trailing whitespace and tabs.
//
// Parameters:
//   - line: The line number to retrieve from the stack trace (must be >= 6 and even)
//
// Returns:
//   - string: The requested line from the stack trace
//
// Panics:
//   - If line is less than 6
//   - If line is not an even number

func GetLine() string {
	return Stack(8)
}

func Stack(line uint16) string {
	if line < 6 {
		panic("the line should be bigger than 5")
	}
	if line%2 != 0 {
		panic("the line should be even")
	}

	stack := string(debug.Stack())
	stack = strings.Split(stack, "\n")[line]
	stack = strings.TrimSpace(strings.ReplaceAll(stack, "\t", ""))

	return stack
}

func getLineFromFile(filePath string, lineNumber int) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	if lineNumber < 1 || lineNumber > len(lines) {
		return "", fmt.Errorf("invalid line number: %d", lineNumber)
	}

	return lines[lineNumber-1], nil // Adjust for 0-based indexing
}

func parseLocation(location string) (string, int, error) {
	parts := strings.SplitN(location, ".go:", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid location format: %s", location)
	}

	filePath := parts[0] + ".go"
	filePath = strings.TrimSpace(strings.ReplaceAll(filePath, " ", ""))

	var numberString string
	for _, v := range parts[1] {
		if v == ' ' {
			break
		}
		numberString += string(v)
	}

	lineNumber, err := strconv.Atoi(numberString)
	if err != nil {
		return "", 0, fmt.Errorf("invalid line number: %s", err)
	}

	return filePath, lineNumber, nil
}

func getLineFromStack(stack string) string {
	filePath, lineNumber, err := parseLocation(stack)
	PanicIfErr(err)
	s, err := getLineFromFile(filePath, lineNumber)
	PanicIfErr(err)
	return s
}

// Debug prints debugging information to the console, including the call stack,
// the name of the variable being debugged, and its value formatted according to
// the specified format string.
//
// Parameters:
//   - format: A format string that specifies how the value should be displayed.
//   - a: The value to be debugged.
//
// The function extracts the call stack, highlights it in magenta, and retrieves
// the line of code where the Debug function was called. It then parses the line
// to extract the variable name and prints the variable name and its value in
// yellow, followed by a separator line for clarity.
func Debug(format string, a any) {
	stack := Stack(8)
	fmt.Println(ColorMagenta, stack, ColorReset)
	line := getLineFromStack(stack)

	parts := strings.SplitN(line, ", ", 2)
	variableName := parts[1][0 : len(parts[1])-1]

	fmt.Printf("%v%v:%v%"+format+"\n", ColorYellow, variableName, ColorReset, a)
	fmt.Println("________________________________________________________________________________")
}

// Test is a generic testing function that compares two values and provides formatted output.
//
// Parameters:
//   - t: Any type parameter for the values being compared
//   - isPanic: If true, exits program when test fails
//   - print: If true, prints test results regardless of pass/fail
//   - isEqual: Expected equality relationship between actual and expected
//   - line: Line number for stack trace (must be > 8)
//   - format: Printf format string for value output
//   - actual: The value being tested
//   - expected: The value to test against
//
// The function:
//   - Validates line number is > 8
//   - Compares actual vs expected using reflect.DeepEqual
//   - Prints colored stack traces and formatted values
//   - Can exit program on test failure if isPanic is true
//   - Supports testing for both equality and inequality based on isEqual flag
//
// Example usage:
//
//	Test(true, false, true, 10, "%v", actual, expected) // Test equality with panic on failure
//	Test(false, true, false, 8, "%d", val1, val2) // Test inequality with output
func Test[t any](isPanic, print, isEqual bool, line uint16, format string, actual, expected t) {
	if line < 8 {
		panic("the line should be bigger than 8")
	}
	stack := Stack(line)

	printStack := func(color string) {
		fmt.Println(color, stack, ColorReset)
	}
	printActual := func() {
		fmt.Printf("%vActual:%v\n", ColorYellow, ColorReset)
		fmt.Printf("%v%"+format+"%v\n", ColorYellow, actual, ColorReset)
	}
	printExpected := func() {
		fmt.Printf("%vExpected:%v\n", ColorBlue, ColorReset)
		fmt.Printf("%v%"+format+"%v\n", ColorBlue, expected, ColorReset)
	}

	if !reflect.DeepEqual(actual, expected) == isEqual {
		printStack(ColorRed)

		if isEqual {
			fmt.Println("this should equal to each other")
		} else {
			fmt.Println("this should not equal to each other")
		}

		printActual()
		printExpected()
		fmt.Println("________________________________________________________________________________")

		if isPanic {
			os.Exit(1)
		}
	} else {
		printStack(ColorGreen)
	}

	if print {
		if isEqual {
			fmt.Println("this is equal to each other")
		} else {
			fmt.Println("this is not equal to each other")
		}

		printActual()
		if !isEqual {
			printExpected()
		}
		fmt.Println("________________________________________________________________________________")
	}
}

// Benchmark measures the execution time of one or more code blocks over a specified number of loops
// and prints the results in ascending order of execution time.
//
// Parameters:
//   - loops: The number of times each code block should be executed.
//   - codesBlocks: A variadic parameter representing one or more functions (code blocks) to benchmark.
//
// Behavior:
//   - Each code block is executed the specified number of times.
//   - The total execution time for each code block is measured and stored.
//   - The results are sorted by execution time in ascending order.
//   - The function prints the stack trace, followed by the execution time for each code block
//     and the average time per loop.
//
// Example:
//
//	Benchmark(1000, func() {
//	    // Code block 1
//	}, func() {
//	    // Code block 2
//	})
//
// Output:
//
//	The function outputs the benchmark results to the console, including:
//	  - The block index.
//	  - The total execution time for the block.
//	  - The average execution time per loop.
//
// Notes:
//   - The function uses the `time` package to measure execution time.
//   - The `sort` package is used to sort the results by execution time.
func Benchmark(loops uint, codesBlocks ...func()) {
	type element struct {
		blockIndex int
		duration   time.Duration
	}

	var list []element

	for blockIndex, codeBlock := range codesBlocks {
		start := time.Now()

		for i := uint(0); i < loops; i++ {
			codeBlock()
		}

		duration := time.Since(time.Time(start))

		list = append(list, element{
			blockIndex: blockIndex,
			duration:   duration,
		})

	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].duration < list[j].duration
	})

	fmt.Println(ColorCyan, Stack(8), ColorReset)

	for _, v := range list {
		fmt.Printf("block index %v: it takes %v and %v for each loop\n", v.blockIndex, v.duration, v.duration/time.Duration(loops))
	}

	fmt.Println("________________________________________________________________________________")
}

func TestCase[t any](format string, line string, input any, expected, actual t) {
	printName := func() {
		fmt.Printf("%vline: %v%v\n\n", ColorRed, line, ColorReset)
	}
	printInput := func() {
		fmt.Printf("%vInput:%v\n", ColorMagenta, ColorReset)
		fmt.Printf("%v%"+format+"%v\n", ColorMagenta, input, ColorReset)
	}
	printActual := func() {
		fmt.Printf("%vActual:%v\n", ColorYellow, ColorReset)
		fmt.Printf("%v%"+format+"%v\n", ColorYellow, actual, ColorReset)
	}
	printExpected := func() {
		fmt.Printf("%vExpected:%v\n", ColorBlue, ColorReset)
		fmt.Printf("%v%"+format+"%v\n", ColorBlue, expected, ColorReset)
	}

	if !reflect.DeepEqual(actual, expected) {
		printName()
		printInput()
		printActual()
		printExpected()
		fmt.Println("________________________________________________________________________________")
		os.Exit(1)
	}
}
