package go_test

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

// Stack returns a string representation of the caller's stack trace at a specific depth.
// It extracts and formats the 8th line of the stack trace, which typically corresponds
// to the caller's location in the code. The returned string is trimmed of whitespace
// and tab characters for readability.
func Stack() string {
	stack := string(debug.Stack())
	stack = strings.Split(stack, "\n")[8]
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
	stack := Stack()
	fmt.Println(ColorMagenta, stack, ColorReset)
	line := getLineFromStack(stack)

	parts := strings.SplitN(line, ", ", 2)
	variableName := parts[1][0 : len(parts[1])-1]

	fmt.Printf("%v%v:%v%"+format+"\n", ColorYellow, variableName, ColorReset, a)
	fmt.Println("________________________________________________________________________________")
}

// Test is a generic testing utility function that compares two values and optionally prints
// debugging information. It uses reflection to determine equality between the actual and
// expected values.
//
// Type Parameters:
//   - t: The type of the values being compared.
//
// Parameters:
//   - print (bool): If true, the function prints debugging information to the console.
//   - isEqual (bool): Indicates whether the actual and expected values are expected to be equal.
//   - format (string): A format string used for printing the actual and expected values.
//   - actual (t): The actual value to be tested.
//   - expected (t): The expected value to compare against.
//
// Behavior:
//   - If the comparison result (using reflect.DeepEqual) does not match the isEqual parameter,
//     the function prints an error message, the stack trace, and the actual/expected values,
//     then exits the program with a status code of 1.
//   - If the comparison result matches the isEqual parameter, and the print parameter is true,
//     the function prints a success message along with the actual/expected values.
//
// Notes:
//   - The function uses a helper function Stack() to retrieve the stack trace.
//   - Color-coded output is used for better readability, with colors defined by constants
//     such as ColorRed, ColorGreen, ColorYellow, and ColorBlue.
//   - The function terminates the program on failure using os.Exit(1).
func Test[t any](print, isEqual bool, format string, actual, expected t) {
	stack := Stack()

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

		os.Exit(1)
	}
	printStack(ColorGreen)

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

	fmt.Println(ColorCyan, Stack(), ColorReset)

	for _, v := range list {
		fmt.Printf("block index %v: it takes %v and %v for each loop\n", v.blockIndex, v.duration, v.duration/time.Duration(loops))
	}

	fmt.Println("________________________________________________________________________________")
}
