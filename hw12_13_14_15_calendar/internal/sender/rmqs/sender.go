package rmqs

import (
	"context"
	"fmt"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/rmq"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	conf    config.Sender
	db      storage.Storage
	logger  app.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewSender(conf config.Sender, db storage.Storage, logger app.Logger) Sender {
	return Sender{conf: conf, db: db, logger: logger}
}

func (s *Sender) Run(ctx context.Context) {
	connection, channel, err := rmq.SetupRMQ(s.conf.RMQ)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	defer connection.Close()
	defer channel.Close()
	s.logger.Info("sender connection established...")

	s.channel = channel
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
		queue.Name,            // name of the queue
		s.conf.RMQ.RoutingKey, // bindingKey
		s.conf.RMQ.Exchange,   // sourceExchange
		false,                 // noWait
		nil,                   // arguments
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
		s.logger.Info(fmt.Sprintf("Got delivery,[%v] %q", d.DeliveryTag, d.Body))

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
