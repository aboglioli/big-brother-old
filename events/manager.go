package events

import (
	"github.com/aboglioli/big-brother/errors"
)

type Manager interface {
	Publish(exchange, exchangeType, key string, body interface{}) errors.Error
	Consume(exchange, exchangeType, queue, key string) (<-chan Message, errors.Error)
}
