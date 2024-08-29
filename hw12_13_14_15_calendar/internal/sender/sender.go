package sender

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	conf    config.Sender
	db      storage.Storage
	logger  app.Logger
	done    atomic.Bool
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewSender(conf config.Sender, db storage.Storage, logger app.Logger) Sender {
	return Sender{conf: conf, db: db, logger: logger}
}

func (s *Sender) Run(ctx context.Context) {
	s.done.Store(false)
	uri := getRMQConnectionString(s.conf)

	connection, err := amqp.Dial(uri)
	s.conn = connection
	if err != nil {
		s.logger.Error("error dial amqp:", err.Error())
		return
	}
	defer connection.Close()

	s.logger.Info("sender connection established...")

	channel, err := connection.Channel()
	if err != nil {
		s.logger.Error("error get connection channel amqp:", err.Error())
		return
	}
	s.channel = channel

	if err := channel.ExchangeDeclare(
		s.conf.Exchange,     // name
		s.conf.ExchangeType, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		true,                // noWait
		nil,                 // arguments
	); err != nil {
		s.logger.Error("error ExchangeDeclare amqp:", err.Error())
		return
	}

	queue, err := channel.QueueDeclare(
		s.conf.Queue, // name of the queue
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		s.logger.Error("error QueueDeclare amqp:", err.Error())
		return
	}

	if err = channel.QueueBind(
		queue.Name,      // name of the queue
		s.conf.BindKey,  // bindingKey
		s.conf.Exchange, // sourceExchange
		false,           // noWait
		nil,             // arguments
	); err != nil {
		s.logger.Error("error Queue Bind amqp:", err.Error())
	}

	deliveries, err := channel.Consume(
		queue.Name,         // name
		s.conf.ConsumerTag, // consumerTag,
		false,              // noAck
		false,              // exclusive
		false,              // noLocal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		s.logger.Error("error Consume amqp:", err.Error())
	}

	handleDeliveries(s, deliveries)
	ctx.Done()
}

func handleDeliveries(s *Sender, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		s.logger.Info(fmt.Printf("Got delivery,[%v] %q", d.DeliveryTag, d.Body))

		errAck := d.Ack(false)
		if errAck != nil {
			s.logger.Warn(fmt.Sprintf("errAck delivery,[%s]", errAck))
		}
	}
}

func (s *Sender) Cancel() {
	// will close() the deliveries channel
	if err := s.channel.Cancel(s.conf.ConsumerTag, true); err != nil {
		s.logger.Error("Consumer cancel failed:", err.Error())
	}

	if err := s.conn.Close(); err != nil {
		s.logger.Error("AMQP connection close error:", err.Error())
	}
}

func getRMQConnectionString(conf config.Sender) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.UserName, conf.Password, conf.Host, conf.Port)
}
