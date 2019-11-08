package events

import (
	"github.com/aboglioli/big-brother/config"
	"github.com/aboglioli/big-brother/errors"
	"github.com/streadway/amqp"
)

var manager *Manager

type Manager struct {
	connection *amqp.Connection
	emitters   map[string]*amqp.Channel
	consumers  map[string]*amqp.Channel
}

func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{
			emitters: make(map[string]*amqp.Channel),
			consumers: make(map[string]*amqp.Channel),
		}
		manager.Connect()
	}

	return manager
}

func (m *Manager) Connect() (*amqp.Connection, errors.Error) {
	if m.connection == nil {
		conf := config.Get()
		conn, err := amqp.Dial(conf.RabbitURL)
		if err != nil {
			return nil, errors.New("infrastructure/events/manager.Connect", "FAILED_TO_CONNECT", err.Error())
		}

		m.connection = conn
	}

	return m.connection, nil
}

func (m *Manager) Disconnect() {
	if m.connection != nil {
		m.connection.Close()
		m.connection = nil
	}
}

func (m *Manager) FanoutSend(exchange string, body []byte) errors.Error {
	if m.emitters[exchange] == nil {
		ch, err := m.createChannelWithExchange(exchange)
		if err != nil {
			return err
		}

		m.emitters[exchange] = ch
	}

	ch := m.emitters[exchange]

	err := ch.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			Body: []byte(body),
		},
	)
	if err != nil {
		return errors.New("infrastructure/events/manager.FanoutSend", "FAILED_TO_PUBLISH_MESSAGE", err.Error())
	}

	return nil
}

func (m *Manager) Consume(exchange string) (<-chan amqp.Delivery, errors.Error) {
	errGen := errors.FromPath("infrastructure/events/manager.Consume")

	if m.emitters[exchange] == nil {
		ch, err := m.createChannelWithExchange(exchange)
		if err != nil {
			return nil, err
		}

		m.consumers[exchange] = ch
	}

	ch := m.consumers[exchange]

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, errGen("FAILED_TO_DECLARE_QUEUE", err.Error())
	}

	err = ch.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, errGen("FAILED_TO_BIND_QUEUE", err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errGen("FAILED_TO_CONSUME", err.Error())
	}

	return msgs, nil
}

func (m *Manager) createChannelWithExchange(exchange string) (*amqp.Channel, errors.Error) {
	errGen := errors.FromPath("infrastructure/events/manager.FanoutSend")

	ch, err := m.connection.Channel()
	if err != nil {
		return nil, errGen("FAILED_TO_CREATE_CHANNEL", err.Error())
	}

	err = ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errGen("FAILED_TO_DECLARE_EXCHANGE", err.Error())
	}

	return ch, nil
}
