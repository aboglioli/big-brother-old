package bench

import (
	"runtime"
	"testing"
)

type caller struct {
}

func (c *caller) fromArgs(_ string) {

}

func (c *caller) fromRuntime() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

func BenchmarkWithoutRuntimeCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := &caller{}
		c.fromArgs("path")
	}
}

func BenchmarkWithRuntimeCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := &caller{}
		c.fromRuntime()
	}
}
