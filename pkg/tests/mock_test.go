package tests

import (
	"testing"
)

type mockStruct struct {
	Mock
}

func TestAssert(t *testing.T) {
	m := &mockStruct{}
	m.Called("FirstMethod", "one", "two")
	m.Called("SecondMethod", "one", 2)
	m.Called("ThirdMethod", 1, "two")

	res := m.Assert([]Call{
		Call{"FirstMethod", []interface{}{"one", "two"}},
		Call{"SecondMethod", []interface{}{"one", 2}},
		Call{"ThirdMethod", []interface{}{1, "two"}},
	})
	Assert(t, res)

	m.Reset()
	m.Called("WithAny", 1, Any, 3)

	res = m.Assert([]Call{
		Call{"WithAny", []interface{}{1, 45, 3}},
	})
	Assert(t, res)
}
