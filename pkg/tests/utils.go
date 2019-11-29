package tests

import (
	"testing"

	"github.com/aboglioli/big-brother/pkg/errors"
)

// Global
func Assert(t *testing.T, cond bool, msg string) {
	if !cond {
		t.Errorf("ASSERT: %s\n", msg)
	}
}

func Equals(t *testing.T, v1 interface{}, v2 interface{}, msg string) {
	if v1 != v2 {
		t.Errorf("EQUALS: %s\n%v != %v\n", msg, v1, v2)
	}
}

func IsNil(t *testing.T, v interface{}, msg string) {
	if v != nil {
		t.Errorf("IS NOT NIL: %s\n", msg)
	}
}

func IsNotNil(t *testing.T, v interface{}, msg string) {
	if v == nil {
		t.Errorf("IS NIL: %s\n", msg)
	}
}

// Errors
func Ok(t *testing.T, err error, msg string) {
	if err != nil {
		t.Errorf("ERROR %s NOT EXPECTED: %s\n", err.Error(), msg)
	}
}

func Err(t *testing.T, err error, msg string) {
	if err == nil {
		t.Errorf("ERROR EXPECTED: %s\n", msg)
	}
}

func ErrCode(t *testing.T, err errors.Error, code string, msg string) {
	if err == nil || err.Code() != code {
		t.Errorf("CODE %s EXPECTED: %s\n", code, msg)
	}
}
