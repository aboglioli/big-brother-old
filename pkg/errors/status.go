package errors

import "fmt"

type Status struct {
	Generic
	statusCode int
}

func NewStatus(code string) *Status {
	return &Status{
		Generic: *New(code),
	}
}

func (s *Status) SetStatus(code int) *Status {
	s.statusCode = code
	return s
}

func (s *Status) Status() int {
	return s.statusCode
}

func (s *Status) SetPath(path string) *Status {
	s.Generic.SetPath(path)
	return s
}

func (s *Status) SetMessage(msg string, args ...interface{}) *Status {
	s.Generic.SetMessage(msg, args...)
	return s
}

func (s *Status) SetRef(err error) *Status {
	s.Generic.SetRef(err)
	return s
}

func (s *Status) Error() string {
	g := s.Generic.Error()
	if s.statusCode > 0 {
		return fmt.Sprintf("%s\n\t- Status code: %d", g, s.statusCode)
	}

	return g
}
