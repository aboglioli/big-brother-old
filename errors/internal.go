package errors

import "fmt"

type InternalError struct {
	path    string
	code    string
	message string
}

func NewInternal(p string, c string, m string) Error {
	return &InternalError{
		path:    p,
		code:    c,
		message: m,
	}
}

func InternalFromPath(p string) func(c string, m string) Error {
	return func(c string, m string) Error {
		return &InternalError{
			path:    p,
			code:    c,
			message: m,
		}
	}
}

func (e *InternalError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("[!] %s#[%s]", e.path, e.code)
	}
	return fmt.Sprintf("[!] %s#[%s]: %s", e.path, e.code, e.message)
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
