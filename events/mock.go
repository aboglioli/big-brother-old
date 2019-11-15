package events

import (
	"github.com/aboglioli/big-brother/errors"
)

// Message
type mockMessage struct {
	exchange     string
	exchangeType string
	queue        string
	key          string
	body         []byte
}

func (d mockMessage) Body() []byte {
	return []byte(d.body)
}

func (d mockMessage) Ack() {
}

// Manager
var mockMgr *mockManager

type mockManager struct {
	buffer []mockMessage
	ch     chan Message
}

func GetMockManager() *mockManager {
	if mockMgr == nil {
		mockMgr = &mockManager{
			buffer: make([]mockMessage, 0),
			ch:     make(chan Message),
		}
	}

	return mockMgr
}

func (m *mockManager) Publish(exchange, exchangeType, key string, body []byte) errors.Error {
	msg := mockMessage{exchange, exchangeType, key, "", body}
	m.buffer = append(m.buffer, msg)

	go func() {
		m.ch <- msg
	}()

	return nil
}

func (m *mockManager) Consume(exchange, exchangeType, queue, key string) (<-chan Message, errors.Error) {
	return m.ch, nil
}
