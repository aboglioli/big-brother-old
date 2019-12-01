package tests

import (
	"testing"
)

const (
	Any    = "tests.Anything"
	Nil    = "tests.Nil"
	NotNil = "tests.NotNil"
)

type Call struct {
	Func string
	Args []interface{}
}

type Mock struct {
	Calls []Call
}

func (m *Mock) Called(f string, args ...interface{}) {
	m.Calls = append(m.Calls, Call{f, args})
}

func (m *Mock) CallsTo(f string) int {
	count := 0
	for _, c := range m.Calls {
		if c.Func == f {
			count++
		}
	}
	return count
}

func (m *Mock) CountCalls() int {
	return len(m.Calls)
}

func (m *Mock) Assert(t *testing.T, calls []Call) {
	if len(m.Calls) != len(calls) {
		t.Fatalf("MOCK: Different number of calls\n%s\n", printStackInfo())
	}

	for i, call1 := range m.Calls {
		call2 := calls[i]

		if call1.Func != call2.Func || !compareArgs(call1.Args, call2.Args) {
			t.Fatalf("MOCK:\nexpected: %s {%v}\nactual: %s {%v}\n%s\n", call1.Func, call1.Args, call2.Func, call2.Args, printStackInfo())
		}
	}
}

func (m *Mock) Reset() {
	m.Calls = []Call{}
}

func compareArgs(args1 []interface{}, args2 []interface{}) bool {
	if len(args1) != len(args2) {
		return false
	}

	for i, arg1 := range args1 {
		arg2 := args2[i]
		if arg1 == Any || arg2 == Any {
			continue
		}

		if arg2 == Nil && arg1 == nil {
			continue
		}

		if arg2 == NotNil && arg1 != nil {
			continue
		}

		if arg1 != arg2 {
			return false
		}
	}

	return true
}
