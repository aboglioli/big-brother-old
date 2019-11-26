package events

import (
	"github.com/aboglioli/big-brother/errors"
)

type Message interface {
	Body() []byte
	Type() string
	Decode(dst interface{}) errors.Error
	Ack()
}

type Manager interface {
	Publish(exchange, exchangeType, key string, body interface{}) errors.Error
	Consume(exchange, exchangeType, queue, key string) (<-chan Message, errors.Error)
}

type Converter interface {
	Decode(src []byte, dst interface{}) errors.Error
	Code(src interface{}) ([]byte, errors.Error)
}
