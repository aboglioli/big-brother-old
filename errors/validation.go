package errors

import "fmt"

type ValidationError struct {
	kind      string
	path      string
	code      string
	message   string
	reference Error // TODO: not used yet
}

func NewValidation() *ValidationError {
	return &ValidationError{
		kind: "validation",
	}
}

func (e *ValidationError) SetPath(p string) *ValidationError {
	e.path = p
	return e
}

func (e *ValidationError) SetCode(c string) *ValidationError {
	e.code = c
	return e
}

func (e *ValidationError) SetMessage(m string) *ValidationError {
	e.message = m
	return e
}

func (e *ValidationError) SetReference(ref Error) *ValidationError {
	e.reference = ref
	return e
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
		return fmt.Sprintf("%s: %s", e.path, e.code)
	}
	return fmt.Sprintf("%s: %s\n\t%s", e.path, e.code, e.message)
}

func (e *ValidationError) String() string {
	return e.Error()
}

func (e *ValidationError) Reference() Error {
	return e.reference
}
