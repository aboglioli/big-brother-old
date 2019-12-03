package errors

type Code interface {
	Code() string
}

type Path interface {
	Path() string
}

type Message interface {
	Message() string
}
