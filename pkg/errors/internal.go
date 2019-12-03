package errors

type Internal struct {
	Generic
}

func NewInternal(code string) *Internal {
	return &Internal{*New(code)}
}

func (i *Internal) SetPath(path string) *Internal {
	i.Generic.SetPath(path)
	return i
}

func (i *Internal) SetMessage(msg string, args ...interface{}) *Internal {
	i.Generic.SetMessage(msg, args...)
	return i
}

func (i *Internal) SetRef(err Error) *Internal {
	i.Generic.SetRef(err)
	return i
}
