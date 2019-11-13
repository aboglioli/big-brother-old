package errors

import "fmt"

type InternalError struct {
	kind      string
	path      string
	code      string
	message   string
	reference Error
}

func NewInternal(p string, c string, m string) Error {
	return &InternalError{
		kind:    "internal",
		path:    p,
		code:    c,
		message: m,
	}
}

func InternalFromPath(p string) func(string, string) Error {
	return func(c string, m string) Error {
		return NewInternal(p, c, m)
	}
}

func (e *InternalError) Kind() string {
	return e.kind
}

func (e *InternalError) Path() string {
	return e.path
}

func (e *InternalError) Message() string {
	return e.message
}

func (e *InternalError) Code() string {
	return e.code
}

func (e *InternalError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("[!] %s#[%s]", e.path, e.code)
	}
	return fmt.Sprintf("[!] %s#[%s]: %s", e.path, e.code, e.message)
}

func (e *InternalError) String() string {
	return e.Error()
}

func (e *InternalError) Reference() Error {
	return e.reference
}
