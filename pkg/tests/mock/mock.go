package mock

import (
	"fmt"
	"testing"

	"github.com/aboglioli/big-brother/pkg/tests"
)

const (
	Any    = "mock.Anything"
	Nil    = "mock.Nil"
	NotNil = "mock.NotNil"
)

type call struct {
	Func string
	Args []interface{}
}

func Call(f string, args ...interface{}) call {
	return call{f, args}
}

type Mock struct {
	Calls []call
}

func (m *Mock) Called(f string, args ...interface{}) {
	m.Calls = append(m.Calls, call{f, args})
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

func (m *Mock) Assert(t *testing.T, calls ...call) {
	if len(m.Calls) != len(calls) {
		callsStr := "expected\n"
		for _, call := range calls {
			callsStr += fmt.Sprintf("- %s {%v}\n", call.Func, call.Args)
		}

		callsStr += "actual:\n"
		for _, call := range m.Calls {
			callsStr += fmt.Sprintf("- %s {%v}\n", call.Func, call.Args)
		}
		t.Fatalf("MOCK: Different number of calls\n%s%s\n", callsStr, tests.PrintStackInfo())
	}

	for i, call1 := range m.Calls {
		call2 := calls[i]

		if call1.Func != call2.Func || !compareArgs(call1.Args, call2.Args) {
			t.Fatalf("MOCK:\nexpected: %s {%v}\nactual: %s {%v}\n%s\n", call2.Func, call2.Args, call1.Func, call1.Args, tests.PrintStackInfo())
		}
	}
}

func (m *Mock) Reset() {
	m.Calls = []call{}
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
