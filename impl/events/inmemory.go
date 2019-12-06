package events

import (
	"github.com/aboglioli/big-brother/pkg/converter"
	"github.com/aboglioli/big-brother/pkg/events"
)

// Message
type inmemoryMessage struct {
	converter    converter.Converter
	Exchange     string
	ExchangeType string
	Queue        string
	Key          string
	body         []byte
}

func (d inmemoryMessage) Body() []byte {
	return d.body
}

func (d inmemoryMessage) Type() string {
	var e events.Event
	if err := d.Decode(&e); err != nil {
		return ""
	}
	return e.Type
}

func (d inmemoryMessage) Decode(dst interface{}) error {
	return d.converter.Decode(d.Body(), dst)
}

func (d inmemoryMessage) Ack() {
}

// Manager
type inmemoryManager struct {
	converter converter.Converter
	ch        chan events.Message
	buffer    []inmemoryMessage
}

func InMemory() *inmemoryManager {
	converter := converter.DefaultConverter()
	mockMgr := &inmemoryManager{
		converter: converter,
		ch:        make(chan events.Message),
		buffer:    make([]inmemoryMessage, 0),
	}

	return mockMgr
}

func (m *inmemoryManager) Publish(body interface{}, opts *events.Options) error {
	b, err := m.converter.Encode(body)
	if err != nil {
		return err
	}
	msg := inmemoryMessage{
		converter:    m.converter,
		Exchange:     opts.Exchange,
		ExchangeType: opts.ExchangeType,
		Key:          opts.Key,
		body:         b,
	}
	m.buffer = append(m.buffer, msg)

	go func() {
		m.ch <- msg
	}()
	return nil
}

func (m *inmemoryManager) Consume(opts *events.Options) (<-chan events.Message, error) {
	return m.ch, nil
}

func (m *inmemoryManager) Messages() []inmemoryMessage {
	return m.buffer
}

func (m *inmemoryManager) Count() int {
	return len(m.buffer)
}

func (m *inmemoryManager) Clean() {
	m.buffer = make([]inmemoryMessage, 0)
}
