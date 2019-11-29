package tests

const (
	Any = "tests.Anything"
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

func (m *Mock) Assert(calls []Call) bool {
	if len(m.Calls) != len(calls) {
		return false
	}

	for i, call1 := range m.Calls {
		call2 := calls[i]

		if call1.Func != call2.Func || !compareArgs(call1.Args, call2.Args) {
			return false
		}
	}

	return true
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

		if arg1 != arg2 {
			return false
		}
	}

	return true
}
