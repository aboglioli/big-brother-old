package tests

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
)

func printStackInfo() string {
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

// Global
func Assert(t *testing.T, cond bool, msgs ...string) {
	if !cond {
		t.Errorf("ASSERT: %s\n%s\n", msgs, printStackInfo())
	}
}

func Equal(t *testing.T, v1 interface{}, v2 interface{}, msgs ...string) {
	if v1 != v2 {
		t.Errorf("EQUALS: %s\n%v != %v\n%s\n", msgs, v1, v2, printStackInfo())
	}
}

// Errors
func Ok(t *testing.T, err error, msgs ...string) {
	if err != nil {
		t.Errorf("OK: %s\nError %s not expected\n%s\n", msgs, err.Error(), printStackInfo())
	}
}

func Err(t *testing.T, err error, msgs ...string) {
	if err == nil {
		t.Errorf("ERR: %s\nError expected\n%s\n", msgs, printStackInfo())
	}
}

func ErrCode(t *testing.T, err errors.Error, code string, msgs ...string) {
	if err == nil || err.Code() != code {
		t.Errorf("ERR_CODE: %s\nError with code %s expected\n%s\n", msgs, code, printStackInfo())
	}
}
