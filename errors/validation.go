package errors

import "fmt"

type ValidationError struct {
	kind      string
	path      string
	code      string
	message   string
	reference Error // TODO: not used yet
}

func NewValidation(p string, c string, m string) Error {
	return &ValidationError{
		kind:    "validation",
		path:    p,
		code:    c,
		message: m,
	}
}

func ValidationFromPath(p string) func(string, string) Error {
	return func(c string, m string) Error {
		return NewValidation(p, c, m)
	}
}

func (e *ValidationError) Kind() string {
	return e.kind
}

func (e *ValidationError) Path() string {
	return e.path
}

func (e *ValidationError) Message() string {
	return e.message
}

func (e *ValidationError) Code() string {
	return e.code
}

func (e *ValidationError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("%s#[%s]", e.path, e.code)
	}
	return fmt.Sprintf("%s#[%s]: %s", e.path, e.code, e.message)
}

func (e *ValidationError) String() string {
	return e.Error()
}

func (e *ValidationError) Reference() Error {
	return e.reference
}
