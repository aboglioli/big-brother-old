package events

import (
	"github.com/aboglioli/big-brother/pkg/converter"
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/aboglioli/big-brother/pkg/tests/mock"
)

// Message
type mockMessage struct {
	converter    converter.Converter
	Exchange     string
	ExchangeType string
	Queue        string
	Key          string
	body         []byte
}

func (d mockMessage) Body() []byte {
	return d.body
}

func (d mockMessage) Type() string {
	var e events.Event
	if err := d.Decode(&e); err != nil {
		return ""
	}
	return e.Type
}

func (d mockMessage) Decode(dst interface{}) error {
	return d.converter.Decode(d.Body(), dst)
}

func (d mockMessage) Ack() {
}

// Manager
type mockManager struct {
	mock.Mock
	converter converter.Converter
	ch        chan events.Message
	buffer    []mockMessage
}

func Mock() *mockManager {
	converter := converter.DefaultConverter()
	mockMgr := &mockManager{
		converter: converter,
		ch:        make(chan events.Message),
		buffer:    make([]mockMessage, 0),
	}

	return mockMgr
}

func (m *mockManager) Publish(body interface{}, opts *events.Options) error {
	call := mock.Call("Publish", body, opts)

	b, err := m.converter.Encode(body)
	if err != nil {
		return err
	}
	msg := mockMessage{
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

	m.Mock.Called(call.Return(nil))
	return nil
}

func (m *mockManager) Consume(opts *events.Options) (<-chan events.Message, error) {
	call := mock.Call("Consume", opts)
	m.Mock.Called(call.Return(m.ch, nil))
	return m.ch, nil
}

func (m *mockManager) Messages() []mockMessage {
	return m.buffer
}

func (m *mockManager) Count() int {
	return len(m.buffer)
}

func (m *mockManager) Clean() {
	m.buffer = make([]mockMessage, 0)
}
