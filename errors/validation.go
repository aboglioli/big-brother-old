package errors

import "fmt"

type ValidationError struct {
	path    string
	code    string
	message string
}

func NewValidation(p string, c string, m string) Error {
	return &ValidationError{
		path:    p,
		code:    c,
		message: m,
	}
}

func ValidationFromPath(p string) func(c string, m string) Error {
	return func(c string, m string) Error {
		return NewValidation(p, c, m)
	}
}

func (e *ValidationError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("%s#[%s]", e.path, e.code)
	}
	return fmt.Sprintf("%s#[%s]: %s", e.path, e.code, e.message)
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
