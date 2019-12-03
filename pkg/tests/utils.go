package tests

import (
	"fmt"
	"runtime"
	"strings"
)

func PrintStackInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("Couldn't get the caller information")
	}

	fn := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fn, "/")
	fn = parts[len(parts)-1]

	parts = strings.Split(file, "/")
	file = parts[len(parts)-1]

	return fmt.Sprintf("in %s (%s, line %d)", fn, file, line)
}
