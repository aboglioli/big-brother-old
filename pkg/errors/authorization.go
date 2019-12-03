package errors

type Unauthorized struct {
	code    string
	message string
}

func NewUnauthorized() *Unauthorized {
	return &Unauthorized{
		code:    "UNAUTHORIZED",
		message: "Unauthorized",
	}
}

func (u *Unauthorized) Code() string {
	return u.code
}

func (u *Unauthorized) Message() string {
	return u.message
}

func (u *Unauthorized) Error() string {
	return u.Code()
}

func (u *Unauthorized) Reference() error {
	return nil
}
