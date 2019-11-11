package events

import (
	"github.com/aboglioli/big-brother/config"
	"github.com/aboglioli/big-brother/errors"
	"github.com/streadway/amqp"
)

// Message
type message struct {
	amqp.Delivery
}

func newMessage(d amqp.Delivery) Message {
	return message{d}
}

func (d message) Body() []byte {
	return d.Delivery.Body
}

func (d message) Ack() {
	d.Delivery.Ack(false)
}

// Manager
var mgr *manager

type manager struct {
	connection *amqp.Connection
	emitters   map[string]*amqp.Channel
	consumers  map[string]*amqp.Channel
}

func GetManager() (Manager, errors.Error) {
	if mgr == nil {
		mgr = &manager{
			emitters:  make(map[string]*amqp.Channel),
			consumers: make(map[string]*amqp.Channel),
		}
		_, err := mgr.Connect()
		if err != nil {
			return nil, err
		}
		go func() {
			for {
				<-mgr.connection.NotifyClose(make(chan *amqp.Error))
				mgr.Connect()
			}
		}()
	}

	return mgr, nil
}

func (m *manager) Connect() (*amqp.Connection, errors.Error) {
	if m.connection != nil {
		m.connection.Close()
	}

	conf := config.Get()
	conn, err := amqp.Dial(conf.RabbitURL)
	if err != nil {
		return nil, errors.New("infrastructure/events/manager.Connect", "FAILED_TO_CONNECT", err.Error())
	}

	m.connection = conn

	return m.connection, nil
}

func (m *manager) Disconnect() {
	if m.connection != nil {
		m.connection.Close()
		m.connection = nil
	}
}

func (m *manager) Publish(exchange string, eType string, key string, body []byte) errors.Error {
	if m.emitters[exchange] == nil {
		ch, err := m.createChannelWithExchange(exchange, eType)
		if err != nil {
			return err
		}

		m.emitters[exchange] = ch
	}

	ch := m.emitters[exchange]

	err := ch.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			Body: []byte(body),
		},
	)
	if err != nil {
		return errors.New("infrastructure/events/manager.Send", "FAILED_TO_PUBLISH_MESSAGE", err.Error())
	}

	return nil
}

func (m *manager) Consume(exchange string, eType string, key string) (<-chan Message, errors.Error) {
	errGen := errors.FromPath("infrastructure/events/manager.Consume")

	if m.emitters[exchange] == nil {
		ch, err := m.createChannelWithExchange(exchange, eType)
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
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, errGen("FAILED_TO_BIND_QUEUE", err.Error())
	}

	delivery, err := ch.Consume(
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

	msg := make(chan Message)
	go func() {
		for d := range delivery {
			msg <- newMessage(d)
		}
		close(msg)
	}()

	return msg, nil

	// return msg, nil
}

func (m *manager) createChannelWithExchange(exchange string, eType string) (*amqp.Channel, errors.Error) {
	errGen := errors.FromPath("infrastructure/events/manager.createChannelWithExchange")

	ch, err := m.connection.Channel()
	if err != nil {
		return nil, errGen("FAILED_TO_CREATE_CHANNEL", err.Error())
	}

	err = ch.ExchangeDeclare(
		exchange,
		eType,
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
