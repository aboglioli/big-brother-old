package events

type Message interface {
	Body() []byte
	Type() string
	Decode(dst interface{}) error
	Ack()
}
