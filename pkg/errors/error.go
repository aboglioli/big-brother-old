package errors

import "fmt"

type ErrorKind string

const (
	VALIDATION ErrorKind = "Validation"
	INTERNAL   ErrorKind = "Internal"
)

// Error is the main interface
// There are two kinds of errors:
// - Validation: can be displayed to the user
// - Internal: contains sensitive data
type Error interface {
	Kind() ErrorKind
	Path() string
	Code() string
	Message() string

	Reference() Error

	Error() string
	String() string
}

// Custom implementation
type ErrField struct {
	path string
	code string
}

type CustomError struct {
	kind      ErrorKind
	path      string
	code      string
	message   string
	fields    []ErrField
	reference Error
}

func NewValidation() *CustomError {
	return &CustomError{
		kind:   VALIDATION,
		code:   "VALIDATION",
		fields: []ErrField{},
	}
}

func NewInternal() *CustomError {
	return &CustomError{
		kind:   INTERNAL,
		code:   "INTERNAL",
		fields: []ErrField{},
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

func (e *CustomError) AddPath(path string, code string) *CustomError {
	e.fields = append(e.fields, ErrField{
		path: path,
		code: code,
	})
	return e
}

func (e *CustomError) Kind() ErrorKind {
	return e.kind
}

func (e *CustomError) Path() string {
	return e.path
}

func (e *CustomError) Code() string {
	return e.code
}

func (e *CustomError) Message() string {
	return e.message
}

func (e *CustomError) Reference() Error {
	return e.reference
}

func (e *CustomError) Error() string {
	str := ""
	if e.message == "" {
		str = fmt.Sprintf("%s: %s", e.path, e.code)
	} else {
		str = fmt.Sprintf("%s: %s\n\t%s", e.path, e.code, e.message)
	}

	if e.reference != nil {
		str += "\n- "
		str += e.reference.Error()
	}

	return str
}

func (e *CustomError) String() string {
	return e.Error()
}
