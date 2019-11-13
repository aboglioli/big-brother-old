package errors

// Error is the main interface
// There are two kinds of errors:
// - Validation: can be displayed to the user
// - Internal: contains sensitive data
type Error interface {
	Error() string
	Path() string
	Code() string
	Message() string
}
