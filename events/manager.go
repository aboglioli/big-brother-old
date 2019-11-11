package events

import (
	"github.com/aboglioli/big-brother/errors"
)

type Message interface {
	Body() []byte
	Ack()
}

type Manager interface {
	Publish(exchange string, eType string, key string, body []byte) errors.Error
	Consume(exchange string, eType string, key string) (<-chan Message, errors.Error)
}
