package golang_test_packege

import (
	"testing"
	"time"

	"github.com/samber/lo"
)

func Test_Test(t *testing.T) {
	Test(false, true, "#v", 1, 1)
	Test(true, true, "#v", 1, 1)
	Test(false, false, "#v", 1, 2)
	Test(true, false, "#v", 1, 2)
	Test(false, true, "#v", 1, 2)
}

func Test1(t *testing.T) {
	names := lo.Uniq[string]([]string{"Samuel", "John", "Samuel"})
	Debug("v", names)
	Debug("x", names)
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

func Test7(t *testing.T) {
	age := 10
	Debug("v", age)
	Debug("v", Stack())
	Debug("v", Rand[int]())
}
