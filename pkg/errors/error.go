package errors

// Error is the main interface
type Error interface {
	Path() string
	Code() string
	Message() string

	Error() string
}

type Code interface {
	Code() string
}

type Path interface {
	Path() string
}

type Message interface {
	Message() string
}
