package events

import (
	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/streadway/amqp"
)

// Message
type rabbitMessage struct {
	amqp.Delivery
	converter events.Converter
}

func newMessage(d amqp.Delivery, c events.Converter) events.Message {
	return rabbitMessage{d, c}
}

func (d rabbitMessage) Body() []byte {
	return d.Delivery.Body
}

func (d rabbitMessage) Type() string {
	var e events.Event
	if err := d.Decode(&e); err != nil {
		return ""
	}
	return e.Type
}

func (d rabbitMessage) Decode(dst interface{}) errors.Error {
	return d.converter.Decode(d.Body(), dst)
}

func (d rabbitMessage) Ack() {
	d.Delivery.Ack(false)
}

// Manager
var mgr *manager

type manager struct {
	connection *amqp.Connection
	emitters   map[string]*amqp.Channel
	consumers  map[string]*amqp.Channel
	converter  events.Converter
}

func GetManager() (events.Manager, errors.Error) {
	if mgr == nil {
		converter := events.DefaultConverter()
		mgr = &manager{
			emitters:  make(map[string]*amqp.Channel),
			consumers: make(map[string]*amqp.Channel),
			converter: converter,
		}
		_, err := mgr.connect()
		if err != nil {
			return nil, err
		}
		go func() {
			for {
				<-mgr.connection.NotifyClose(make(chan *amqp.Error))
				mgr.connect()
			}
		}()
	}

	return mgr, nil
}

func (m *manager) connect() (*amqp.Connection, errors.Error) {
	if m.connection != nil {
		m.connection.Close()
	}

	conf := config.Get()
	conn, err := amqp.Dial(conf.RabbitURL)
	if err != nil {
		return nil, errors.NewInternal("FAILED_TO_CONNECT").SetPath("infrastructure/events/manager.Connect").SetMessage(err.Error())
	}

	m.connection = conn

	return m.connection, nil
}

func (m *manager) disconnect() {
	if m.connection != nil {
		m.connection.Close()
		m.connection = nil
	}
}

func (m *manager) Publish(body interface{}, opts *events.Options) errors.Error {
	if m.emitters[opts.Exchange] == nil {
		ch, err := m.createChannelWithExchange(opts.Exchange, opts.ExchangeType)
		if err != nil {
			return err
		}

		m.emitters[opts.Exchange] = ch
	}

	ch := m.emitters[opts.Exchange]

	b, err := m.converter.Encode(body)
	if err != nil {
		return err
	}

	if err := ch.Publish(
		opts.Exchange,
		opts.Key,
		false,
		false,
		amqp.Publishing{
			Body: b,
		},
	); err != nil {
		return errors.NewInternal("FAILED_TO_PUBLISH_MESSAGE").SetPath("infrastructure/events/manager.Publish").SetMessage(err.Error())
	}

	return nil
}

func (m *manager) Consume(opts *events.Options) (<-chan events.Message, errors.Error) {
	path := "infrastructure/events/manager.Consume"

	if m.emitters[opts.Exchange] == nil {
		ch, err := m.createChannelWithExchange(opts.Exchange, opts.ExchangeType)
		if err != nil {
			return nil, err
		}

		m.consumers[opts.Exchange] = ch
	}

	ch := m.consumers[opts.Exchange]

	exclusive := true
	if opts.Queue != "" {
		exclusive = false
	}

	q, err := ch.QueueDeclare(
		opts.Queue,
		false,
		false,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.NewInternal("FAILED_TO_DECLARE_QUEUE").SetMessage(err.Error())
	}

	err = ch.QueueBind(
		q.Name,
		opts.Key,
		opts.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.NewInternal("FAILED_TO_BIND_QUEUE").SetPath(path).SetMessage(err.Error())
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
		return nil, errors.NewInternal("FAILED_TO_CONSUME").SetPath(path).SetMessage(err.Error())
	}

	msg := make(chan events.Message)
	go func() {
		for d := range delivery {
			msg <- newMessage(d, m.converter)
		}
		close(msg)
	}()

	return msg, nil
}

func (m *manager) createChannelWithExchange(exchange string, eType string) (*amqp.Channel, errors.Error) {
	path := "infrastructure/events/manager.createChannelWithExchange"

	ch, err := m.connection.Channel()
	if err != nil {
		return nil, errors.NewInternal("FAILED_TO_CREATE_CHANNEL").SetPath(path).SetMessage(err.Error())
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
		return nil, errors.NewInternal("FAILED_TO_DECLARE_EXCHANGE").SetPath(path).SetMessage(err.Error())
	}

	return ch, nil
}
