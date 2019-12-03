package assert

import (
	"reflect"
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
	value := reflect.ValueOf(v)
	if value.IsNil() {
		t.Fatalf("NIL: %s\n%s\n", msgs, tests.PrintStackInfo())
	}
}

func Nil(t *testing.T, v interface{}, msgs ...string) {
	value := reflect.ValueOf(v)
	if !value.IsNil() {
		t.Fatalf("NOT NIL: %s\n%v should be nil\n%s\n", msgs, v, tests.PrintStackInfo())
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

func ErrCode(t *testing.T, err error, code string, msgs ...string) {
	if err == nil {
		t.Fatalf("ERR: %s\nexpected: error with code %s\nactual: nil error\n%s\n", msgs, code, tests.PrintStackInfo())
		return
	}

	c, ok := err.(errors.Code)
	if !ok {
		t.Fatalf("ERR: %s\nexpected: error with code %s\nactual: not an error with Code\n%s\n", msgs, code, tests.PrintStackInfo())
		return
	}

	if c.Code() != code {
		t.Fatalf("ERR: %s\nexpected: error with code %s\nactual: %s\n%s\n", msgs, code, c.Code(), tests.PrintStackInfo())
	}
}

func ErrMessage(t *testing.T, err error, msg string, msgs ...string) {
	if err == nil {
		t.Fatalf("ERR: %s\nexpected: error with message %s\nactual: nil error\n%s\n", msgs, msg, tests.PrintStackInfo())
		return
	}

	m, ok := err.(errors.Message)
	if !ok {
		t.Fatalf("ERR: %s\nexpected: error with message %s\nactual: not an error with Message\n%s\n", msgs, msg, tests.PrintStackInfo())
		return
	}

	if m.Message() != msg {
		t.Fatalf("ERR: %s\nexpected: error with message %s\nactual: %s\n%s\n", msgs, msg, m.Message(), tests.PrintStackInfo())
	}
}

func ErrValidation(t *testing.T, err error, path string, code string, msgs ...string) {
	if err == nil {
		t.Fatalf("ERR: %s\nexpected: Validation error\nactual: nil error\n%s\n", msgs, tests.PrintStackInfo())
		return
	}

	v, ok := err.(*errors.Validation)
	if !ok {
		t.Fatalf("ERR: %s\nexpected: Validation error\nactual: not a Validation error\n%s\n", msgs, tests.PrintStackInfo())
		return
	}

	fieldExists := false
	for _, f := range v.Fields() {
		if f.Path == path && f.Code == code {
			fieldExists = true
			break
		}
	}

	if !fieldExists {
		t.Fatalf("ERR: %s\nexpected: Validation error with field {%s, %s}\nactual: fields %s\n%s\n", msgs, path, code, v.Fields(), tests.PrintStackInfo())
	}
}
