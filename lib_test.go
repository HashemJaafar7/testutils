package testutils

import (
	"math"
	"testing"
	"time"
)

func Test_Test(t *testing.T) {
	Test(true, false, true, 8, "#v", 1, 1)
	Test(true, true, true, 8, "#v", 1, 1)
	Test(true, false, false, 8, "#v", 1, 2)
	Test(true, true, false, 8, "#v", 1, 2)
	// Test(true, false, true,8, "#v", 1, 2)
}

func Test1(t *testing.T) {
	names := []string{"Samuel", "John", "Samuel"}
	Debug("v", names)
	Debug("x", names)

	age := 10
	Debug("v", age)
	Debug("v", Stack(8))
	Debug("v", Rand[int]())
}

func Test3(t *testing.T) {
	for i := 0; i < 1000; i++ {
		println(Rand[int]())
	}
}

func Test4(t *testing.T) {
	Benchmark(10,
		func() {
			time.Sleep(50 * time.Millisecond)
		},
		func() {
			time.Sleep(100 * time.Millisecond)
		},
		func() {
			time.Sleep(500 * time.Millisecond)
		},
		func() {
			time.Sleep(10 * time.Millisecond)
		},
	)
}

func TestTestCase(t *testing.T) {

	type input struct {
		X float64
		Y float64
	}
	type output struct {
		z   float64
		err error
	}
	tests := []struct {
		Name   string
		Input  input
		Output output
	}{
		{
			Name: "pass",
			Input: input{
				X: 1,
				Y: 5,
			},
			Output: output{
				z:   5,
				err: nil,
			},
		},
		{
			Name: "pass",
			Input: input{
				X: 6,
				Y: 5,
			},
			Output: output{
				z:   6,
				err: nil,
			},
		},
		// {
		// 	Name: "hashem",
		// 	Input: input{
		// 		X: 1,
		// 		Y: 7,
		// 	},
		// 	Output: output{
		// 		z:   5,
		// 		err: nil,
		// 	},
		// },
	}
	for _, tt := range tests {
		var act output
		act.z = math.Max(tt.Input.X, tt.Input.Y)
		act.err = nil
		TestCase("v", tt.Name, tt.Input, tt.Output, act)
	}
}
