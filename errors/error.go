package errors

import "fmt"

// Error is the main interface
// There are two kinds of errors:
// - Validation: can be displayed to the user
// - Internal: contains sensitive data
type Error interface {
	Kind() string
	Path() string
	Code() string
	Message() string

	Reference() Error

	Error() string
	String() string
}

// Custom implementation
type CustomError struct {
	kind      string
	path      string
	code      string
	message   string
	reference Error // TODO: not used yet
}

func NewValidation() *CustomError {
	return &CustomError{
		kind: "validation",
	}
}

func NewInternal() *CustomError {
	return &CustomError{
		kind: "internal",
	}
}

func (e *CustomError) SetPath(p string) *CustomError {
	e.path = p
	return e
}

func (e *CustomError) SetCode(c string) *CustomError {
	e.code = c
	return e
}

func (e *CustomError) SetMessage(m string) *CustomError {
	e.message = m
	return e
}

func (e *CustomError) SetReference(ref Error) *CustomError {
	e.reference = ref
	return e
}

func (e *CustomError) Kind() string {
	return e.kind
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

func (e *CustomError) Error() string {
	str := ""
	if e.message == "" {
		str = fmt.Sprintf("%s: %s", e.path, e.code)
	} else {
		str = fmt.Sprintf("%s: %s\n\t%s", e.path, e.code, e.message)
	}

	if e.reference != nil {
		str += "\nStack:\n- "
		str += e.reference.Error()
	}

	return str
}

func (e *CustomError) String() string {
	return e.Error()
}

func (e *CustomError) Reference() Error {
	return e.reference
}
