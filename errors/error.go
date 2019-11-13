package errors

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
