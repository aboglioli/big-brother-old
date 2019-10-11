package errors

import "fmt"

type Error interface {
	Error() string
	Path() string
	Code() string
	Message() string
}

type CustomError struct {
	path    string
	code    string
	message string
}

func New(p string, c string, m string) Error {
	return &CustomError{
		path:    p,
		code:    c,
		message: m,
	}
}

func FromPath(p string) func(c string, m string) Error {
	return func(c string, m string) Error {
		return New(p, c, m)
	}
}

func (e *CustomError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("[%s]<%s>", e.path, e.code)
	}
	return fmt.Sprintf("[%s]<%s>: %s", e.path, e.code, e.message)
}

func (e *CustomError) Path() string {
	return e.path
}

func (e *CustomError) Message() string {
	return e.message
}

func (e *CustomError) Code() string {
	return e.code
}
