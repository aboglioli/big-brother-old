package events

import "github.com/aboglioli/big-brother/errors"

type Message interface {
	Body() []byte
	Type() string
	Decode(dst interface{}) errors.Error
	Ack()
}
