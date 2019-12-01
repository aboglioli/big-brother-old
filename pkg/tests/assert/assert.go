package assert

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/tests"
)

// Global
func Assert(t *testing.T, cond bool, msgs ...string) {
	if !cond {
		t.Fatalf("ASSERT: %s\n%s\n", msgs, tests.PrintStackInfo())
	}
}

func NotNil(t *testing.T, v interface{}, msgs ...string) {
	if v == nil {
		t.Fatalf("NOT NIL: %s\n%s\n", msgs, tests.PrintStackInfo())
	}
}

func Equal(t *testing.T, v1 interface{}, v2 interface{}, msgs ...string) {
	if v1 != v2 {
		t.Fatalf("EQUAL: %s\n%v != %v\n%s\n", msgs, v1, v2, tests.PrintStackInfo())
	}
}

// Errors
func Ok(t *testing.T, err error, msgs ...string) {
	if err != nil {
		t.Fatalf("OK: %s\nError %s not expected\n%s\n", msgs, err.Error(), tests.PrintStackInfo())
	}
}

func Err(t *testing.T, err error, msgs ...string) {
	if err == nil {
		t.Fatalf("ERR: %s\nError expected\n%s\n", msgs, tests.PrintStackInfo())
	}
}

func ErrCode(t *testing.T, err errors.Error, code string, msgs ...string) {
	if err == nil || err.Code() != code {
		t.Fatalf("ERR_CODE: %s\nError with code %s expected\n%s\n", msgs, code, tests.PrintStackInfo())
	}
}
