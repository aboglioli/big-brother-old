package events

import "github.com/aboglioli/big-brother/pkg/errors"

type Message interface {
	Body() []byte
	Type() string
	Decode(dst interface{}) errors.Error
	Ack()
}
