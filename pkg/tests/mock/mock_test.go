package mock

import (
	"errors"
	"testing"

	"github.com/aboglioli/big-brother/pkg/tests/assert"
)

type mockStruct struct {
	Mock
}

func TestAssert(t *testing.T) {
	m := &mockStruct{}
	m.Called("FirstMethod", "one", "two")
	m.Called("SecondMethod", "one", 2)
	m.Called("ThirdMethod", 1, "two")
	m.Assert(t,
		call{"FirstMethod", []interface{}{"one", "two"}},
		call{"SecondMethod", []interface{}{"one", 2}},
		call{"ThirdMethod", []interface{}{1, "two"}},
	)
	m.Reset()

	m.Called("WithAny", 1, Any, 3)
	m.Assert(t,
		Call("WithAny", 1, 45, 3),
	)
	m.Reset()

	m.Called("One", 1, 2)
	m.Called("Two", 1, 2)
	m.Assert(t,
		Call("One", 1, Any),
		call{"Two", []interface{}{Any, 2}},
	)
	m.Reset()

	m.Called("simple")
	m.Called("simple")
	m.Assert(t,
		Call("simple"),
		Call("simple"),
	)
	m.Reset()

	err := errors.New("Error")
	m.Called("Method1", "data", nil)
	m.Called("Method1", nil, err)
	m.Called("Method2", 1, 2, nil)
	m.Assert(t,
		Call("Method1", "data", Nil),
		Call("Method1", nil, NotNil),
		Call("Method2", 1, NotNil, Nil),
	)
	assert.Equal(t, m.CallsTo("Method1"), 2)
	assert.Equal(t, m.CallsTo("Method2"), 1)
	assert.Equal(t, m.CountCalls(), 3)
	m.Reset()
}
