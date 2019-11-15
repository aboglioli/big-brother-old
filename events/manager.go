package events

import (
	"github.com/aboglioli/big-brother/errors"
)

type Message interface {
	Body() []byte
	Ack()
}

type Manager interface {
	Publish(exchange, exchangeType, key string, body []byte) errors.Error
	Consume(exchange, exchangeType, queue, key string) (<-chan Message, errors.Error)
}
