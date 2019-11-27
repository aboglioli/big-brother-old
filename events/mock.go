package events

import (
	"github.com/aboglioli/big-brother/errors"
)

// Message
type mockMessage struct {
	converter    Converter
	exchange     string
	exchangeType string
	queue        string
	key          string
	body         []byte
}

func (d mockMessage) Body() []byte {
	return d.body
}

func (d mockMessage) Type() string {
	var e EventType
	if err := d.Decode(&e); err != nil {
		return ""
	}
	return e.Type
}

func (d mockMessage) Decode(dst interface{}) errors.Error {
	return d.converter.Decode(d.body, dst)
}

func (d mockMessage) Ack() {
}

// Manager
var mockMgr *mockManager

type mockManager struct {
	converter Converter
	ch        chan Message
	buffer    []mockMessage
}

func GetMockManager() *mockManager {
	if mockMgr == nil {
		converter := DefaultConverter()
		mockMgr = &mockManager{
			converter: converter,
			ch:        make(chan Message),
			buffer:    make([]mockMessage, 0),
		}
	}

	return mockMgr
}

func (m *mockManager) Publish(exchange, exchangeType, key string, body interface{}) errors.Error {
	b, err := m.converter.Code(body)
	if err != nil {
		return err
	}
	msg := mockMessage{
		converter:    m.converter,
		exchange:     exchange,
		exchangeType: exchangeType,
		key:          key,
		body:         b,
	}
	m.buffer = append(m.buffer, msg)

	go func() {
		m.ch <- msg
	}()

	return nil
}

func (m *mockManager) Consume(exchange, exchangeType, queue, key string) (<-chan Message, errors.Error) {
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
