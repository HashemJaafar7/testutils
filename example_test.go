package testutils_test

import (
	"time"

	"github.com/HashemJaafar7/testutils"
)

func Example() {
	// Example of error handling
	testutils.PanicIfErr(nil)

	// Examples of random value generation
	randomInt := testutils.Rand[int]()
	testutils.Debug("v", randomInt)

	randomString := testutils.Rand[string]()
	testutils.Debug("v", randomString)

	// Example of stack trace
	stackTrace := testutils.Stack(8)
	testutils.Debug("v", stackTrace)

	// Examples of Test function
	testutils.Test(true, true, true, 8, "v", 42, 42)            // Equal values
	testutils.Test(true, true, false, 8, "v", "hello", "world") // Different values

	// Example with slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	testutils.Test(true, true, true, 8, "v", slice1, slice2)

	// Example of struct comparison
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{"John", 30}
	p2 := Person{"John", 30}
	testutils.Test(true, true, true, 8, "v", p1, p2)

	// Example of Benchmark
	testutils.Benchmark(5,
		func() {
			time.Sleep(10 * time.Millisecond)
			testutils.Debug("v", "Fast operation")
		},
		func() {
			time.Sleep(50 * time.Millisecond)
			testutils.Debug("v", "Medium operation")
		},
		func() {
			time.Sleep(100 * time.Millisecond)
			testutils.Debug("v", "Slow operation")
		},
	)
}
