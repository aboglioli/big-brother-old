package errors

import "fmt"

type InternalError struct {
	kind      string
	path      string
	code      string
	message   string
	reference Error
}

func NewInternal() *InternalError {
	return &InternalError{
		kind: "internal",
	}
}

func (e *InternalError) SetPath(p string) *InternalError {
	e.path = p
	return e
}

func (e *InternalError) SetCode(c string) *InternalError {
	e.code = c
	return e
}

func (e *InternalError) SetMessage(m string) *InternalError {
	e.message = m
	return e
}

func (e *InternalError) SetReference(ref Error) *InternalError {
	e.reference = ref
	return e
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
