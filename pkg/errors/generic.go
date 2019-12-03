package errors

import "fmt"

type Generic struct {
	path      string
	code      string
	message   string
	reference error
}

func New(code string) *Generic {
	return &Generic{
		code: code,
	}
}

func (c *Generic) SetPath(path string) *Generic {
	c.path = path
	return c
}

func (c *Generic) Path() string {
	return c.path
}

func (c *Generic) Code() string {
	return c.code
}

func (c *Generic) SetMessage(msg string, args ...interface{}) *Generic {
	c.message = fmt.Sprintf(msg, args...)
	return c
}

func (c *Generic) Message() string {
	return c.message
}

func (c *Generic) SetRef(err error) *Generic {
	c.reference = err
	return c
}

func (c *Generic) Reference() error {
	return c.reference
}

func (g *Generic) Error() string {
	str := g.code
	if g.message != "" {
		str += fmt.Sprintf(": %s", g.message)
	}
	if g.path != "" {
		str += fmt.Sprintf("\n\t- path: %s", g.path)
	}
	return str
}
