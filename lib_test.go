package testutils

import (
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
