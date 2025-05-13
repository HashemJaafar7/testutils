package go_test_test

import (
	"time"

	gt "github.com/HashemJaafar7/go_test"
)

func Example() {
	// Example of error handling
	gt.PanicIfErr(nil)

	// Examples of random value generation
	randomInt := gt.Rand[int]()
	gt.Debug("v", randomInt)

	randomString := gt.Rand[string]()
	gt.Debug("v", randomString)

	// Example of stack trace
	stackTrace := gt.Stack(8)
	gt.Debug("v", stackTrace)

	// Examples of Test function
	gt.Test(true, true, true, 8, "v", 42, 42)            // Equal values
	gt.Test(true, true, false, 8, "v", "hello", "world") // Different values

	// Example with slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	gt.Test(true, true, true, 8, "v", slice1, slice2)

	// Example of struct comparison
	type Person struct {
		Name string
		Age  int
	}
	p1 := Person{"John", 30}
	p2 := Person{"John", 30}
	gt.Test(true, true, true, 8, "v", p1, p2)

	// Example of Benchmark
	gt.Benchmark(5,
		func() {
			time.Sleep(10 * time.Millisecond)
			gt.Debug("v", "Fast operation")
		},
		func() {
			time.Sleep(50 * time.Millisecond)
			gt.Debug("v", "Medium operation")
		},
		func() {
			time.Sleep(100 * time.Millisecond)
			gt.Debug("v", "Slow operation")
		},
	)
}
