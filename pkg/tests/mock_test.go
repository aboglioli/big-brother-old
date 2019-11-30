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
	m.Assert(t, []Call{
		Call{"FirstMethod", []interface{}{"one", "two"}},
		Call{"SecondMethod", []interface{}{"one", 2}},
		Call{"ThirdMethod", []interface{}{1, "two"}},
	})
	m.Reset()

	m.Called("WithAny", 1, Any, 3)
	m.Assert(t, []Call{
		Call{"WithAny", []interface{}{1, 45, 3}},
	})
	m.Reset()

	m.Called("One", 1, 2)
	m.Called("Two", 1, 2)
	m.Assert(t, []Call{
		Call{"One", []interface{}{1, Any}},
		Call{"Two", []interface{}{Any, 2}},
	})
	m.Reset()

	m.Called("simple")
	m.Called("simple")
	m.Assert(t, []Call{
		Call{"simple", []interface{}{}},
		Call{"simple", []interface{}{}},
	})
	m.Reset()
}
