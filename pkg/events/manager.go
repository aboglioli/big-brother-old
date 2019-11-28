package events

import (
	"github.com/aboglioli/big-brother/pkg/errors"
)

type Options struct {
	Exchange     string
	ExchangeType string
	Key          string
	Queue        string
}

type Manager interface {
	Publish(body interface{}, opts *Options) errors.Error
	Consume(opts *Options) (<-chan Message, errors.Error)
}
