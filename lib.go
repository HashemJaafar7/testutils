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

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
func Rand[t any]() t {
	var result t
	fuzz.New().Fuzz(&result)
	time.Sleep(1 * time.Microsecond)
	return result
}

const (
	ColorReset = "\033[0m"

	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m" //for fail
	ColorGreen   = "\033[32m" //for success
	ColorYellow  = "\033[33m" //for expected
	ColorBlue    = "\033[34m" //for actual
	ColorMagenta = "\033[35m" //for Debug
	ColorCyan    = "\033[36m" //for Benchmark
	ColorWhite   = "\033[37m"
)

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

func Debug(format string, a any) {
	stack := Stack()
	fmt.Println(ColorMagenta, stack, ColorReset)
	line := getLineFromStack(stack)

	parts := strings.SplitN(line, ", ", 2)
	variableName := parts[1][0 : len(parts[1])-1]

	fmt.Printf("%v%v:%v%"+format+"\n", ColorYellow, variableName, ColorReset, a)
	fmt.Println("________________________________________________________________________________")
}

func Test[t any](print, isEqual bool, format string, actual, expected t) {
	stack := Stack()

	printStack := func(color string) {
		fmt.Println(color, stack, ColorReset)
	}
	printActual := func() {
		fmt.Printf("%v%"+format+"%v\n", ColorYellow, actual, ColorReset)
	}
	printExpected := func() {
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
