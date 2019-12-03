package events

type Options struct {
	Exchange     string
	ExchangeType string
	Key          string
	Queue        string
}

type Manager interface {
	Publish(body interface{}, opts *Options) error
	Consume(opts *Options) (<-chan Message, error)
}
