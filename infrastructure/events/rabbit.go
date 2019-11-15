package events

import (
	"github.com/aboglioli/big-brother/config"
	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/events"
	"github.com/streadway/amqp"
)

// Message
type message struct {
	amqp.Delivery
}

func newMessage(d amqp.Delivery) events.Message {
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

func GetManager() (events.Manager, errors.Error) {
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
		return nil, errors.NewInternal().SetPath("infrastructure/events/manager.Connect").SetCode("FAILED_TO_CONNECT").SetMessage(err.Error())
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

func (m *manager) Publish(exchange, exchangeType, key string, body []byte) errors.Error {
	if m.emitters[exchange] == nil {
		ch, err := m.createChannelWithExchange(exchange, exchangeType)
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
		return errors.NewInternal().SetPath("infrastructure/events/manager.Send").SetCode("FAILED_TO_PUBLISH_MESSAGE").SetMessage(err.Error())
	}

	return nil
}

func (m *manager) Consume(exchange, exchangeType, queue, key string) (<-chan events.Message, errors.Error) {
	errGen := errors.NewInternal().SetPath("infrastructure/events/manager.Consume")

	if m.emitters[exchange] == nil {
		ch, err := m.createChannelWithExchange(exchange, exchangeType)
		if err != nil {
			return nil, err
		}

		m.consumers[exchange] = ch
	}

	ch := m.consumers[exchange]

	exclusive := true
	if queue != "" {
		exclusive = false
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return nil, errGen.SetCode("FAILED_TO_DECLARE_QUEUE").SetMessage(err.Error())
	}

	err = ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, errGen.SetCode("FAILED_TO_BIND_QUEUE").SetMessage(err.Error())
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
		return nil, errGen.SetCode("FAILED_TO_CONSUME").SetMessage(err.Error())
	}

	msg := make(chan events.Message)
	go func() {
		for d := range delivery {
			msg <- newMessage(d)
		}
		close(msg)
	}()

	return msg, nil
}

func (m *manager) createChannelWithExchange(exchange string, eType string) (*amqp.Channel, errors.Error) {
	errGen := errors.NewInternal().SetPath("infrastructure/events/manager.createChannelWithExchange")

	ch, err := m.connection.Channel()
	if err != nil {
		return nil, errGen.SetCode("FAILED_TO_CREATE_CHANNEL").SetMessage(err.Error())
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
		return nil, errGen.SetCode("FAILED_TO_DECLARE_EXCHANGE").SetMessage(err.Error())
	}

	return ch, nil
}
