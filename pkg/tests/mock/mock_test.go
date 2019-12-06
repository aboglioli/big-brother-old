package mock

import (
	"errors"
	"testing"
)

type mockStruct struct {
	Mock
}

func TestMock(t *testing.T) {
	m := &mockStruct{}

	m.Called(Call("FirstMethod", "one", "two"))
	m.Called(Call("SecondMethod", "one", 2))
	m.Called(Call("ThirdMethod", 1, "two"))
	m.Assert(t,
		call{"FirstMethod", []interface{}{"one", "two"}, []interface{}{}},
		call{"SecondMethod", []interface{}{"one", 2}, []interface{}{}},
		call{"ThirdMethod", []interface{}{1, "two"}, []interface{}{}},
	)
	m.Reset()

	m.Called(Call("WithAny", 1, Any, 3))
	m.Assert(t,
		Call("WithAny", 1, 45, 3),
	)
	m.Reset()

	m.Called(Call("One", 1, 2))
	m.Called(Call("Two", 1, 2))
	m.Assert(t,
		Call("One", 1, Any),
		call{"Two", []interface{}{Any, 2}, []interface{}{}},
	)
	m.Reset()

	m.Called(Call("simple"))
	m.Called(Call("simple"))
	m.Assert(t,
		Call("simple"),
		Call("simple"),
	)
	m.Reset()

	err := errors.New("Error")
	m.Called(Call("Method1", "data", nil))
	m.Called(Call("Method1", nil, err))
	m.Called(Call("Method2", 1, 2, nil))
	m.Assert(t,
		Call("Method1", "data", Nil),
		Call("Method1", nil, NotNil),
		Call("Method2", 1, NotNil, Nil),
	)
	m.Reset()

	m.Called(Call("Method1", "data").Return("return", nil))
	m.Called(Call("Method1", "data").Return(nil, err))
	m.Assert(t,
		Call("Method1", "data"),
		Call("Method1", "data"),
	)
	m.Assert(t,
		Call("Method1", "data").Return("return", nil),
		Call("Method1", "data").Return(nil, err),
	)
	m.Assert(t,
		Call("Method1", "data").Return("return", Nil),
		Call("Method1", NotNil).Return(Nil, NotNil),
	)
}
