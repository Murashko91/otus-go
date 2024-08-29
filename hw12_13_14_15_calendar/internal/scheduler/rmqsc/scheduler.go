package rmqsc

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/rmq"
	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Scheduler struct {
	conf   config.Scheduler
	db     storage.Storage
	logger app.Logger
	done   atomic.Bool
}

func NewScheduler(conf config.Scheduler, db storage.Storage, logger app.Logger) *Scheduler {
	return &Scheduler{conf: conf, db: db, logger: logger}
}

func (s *Scheduler) Run(ctx context.Context) {
	s.done.Store(false)

	connection, channel, err := rmq.SetupRMQ(s.conf.RMQ)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	defer connection.Close()
	defer channel.Close()

	s.logger.Info("scheduler connection established...")

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
			s.conf.RMQ.Exchange,
			s.conf.RMQ.RoutingKey,
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
