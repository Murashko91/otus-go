package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Scheduler struct {
	conf   config.Scheduler
	db     storage.Storage
	logger app.Logger
	done   atomic.Bool
}

func NewScheduler(conf config.Scheduler, db storage.Storage, logger app.Logger) Scheduler {

	return Scheduler{conf: conf, db: db, logger: logger}
}

func (s *Scheduler) Run(ctx context.Context) {

	s.done.Store(false)
	uri := getRMQConnectionString(s.conf)

	connection, err := amqp.Dial(uri)
	if err != nil {
		s.logger.Error("error dial amqp:", err.Error())
		return
	}
	defer connection.Close()

	s.logger.Info("scheduler connection established...")

	channel, err := connection.Channel()
	if err != nil {
		s.logger.Error("error get connection channel amqp:", err.Error())
		return
	}

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

	for {
		if s.done.Load() {
			break
		}

		events, err := s.db.GetEventsToSend(ctx)

		if err != nil {
			s.logger.Error("error GetEventsToSend:", err.Error())
			time.Sleep(time.Second * time.Duration(s.conf.IntervalCheck))
			continue
		}

		data, err := json.Marshal(events)

		if err != nil {
			s.logger.Error("error marshal EventsToSend:", err.Error())
			time.Sleep(time.Second * time.Duration(s.conf.IntervalCheck))
			continue
		}

		if err = channel.PublishWithContext(
			context.Background(),
			s.conf.Exchange,
			s.conf.RoutingKey,
			false,
			false,
			amqp.Publishing{
				Headers:      amqp.Table{},
				ContentType:  "application/json",
				Body:         data,
				DeliveryMode: amqp.Transient,
			},
		); err != nil {
			s.logger.Error("error publish data amqp:", err.Error())
			time.Sleep(time.Second * time.Duration(s.conf.IntervalCheck))
		}
		time.Sleep(time.Second * time.Duration(s.conf.IntervalCheck))
	}

}

func (s *Scheduler) Cancel() {

	s.done.Store(true)
}

func getRMQConnectionString(conf config.Scheduler) string {

	fmt.Println("AAAAAA")
	fmt.Println(conf)

	return fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.UserName, conf.Password, conf.Host, conf.Port)

}
