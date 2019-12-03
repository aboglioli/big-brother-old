package errors

type Status struct {
	Generic
	statusCode int
}

func NewStatus(code string) *Status {
	return &Status{
		Generic: *New(code),
	}
}

func (s *Status) SetStatus(code int) {
	s.statusCode = code
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

func (s *Status) SetRef(err Error) *Status {
	s.Generic.SetRef(err)
	return s
}
