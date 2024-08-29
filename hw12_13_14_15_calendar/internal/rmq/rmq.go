package rmq

import (
	"fmt"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRMQ(conf config.RMQ) (*amqp.Connection, *amqp.Channel, error) {

	uri := getRMQConnectionString(conf)

	connection, err := amqp.Dial(uri)
	if err != nil {
		if connection != nil {
			connection.Close()
		}
		return nil, nil, fmt.Errorf("error dial amqp: %s", err.Error())

	}

	channel, err := connection.Channel()
	if err != nil {
		if channel != nil {
			channel.Close()
		}
		return nil, nil, fmt.Errorf("error get connection channel amqp: %s", err.Error())
	}

	if err := channel.ExchangeDeclare(
		conf.Exchange,     // name
		conf.ExchangeType, // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		true,              // noWait
		nil,               // arguments
	); err != nil {
		return nil, nil, fmt.Errorf("error ExchangeDeclare amqp:: %s", err.Error())
	}

	return connection, channel, nil

}

func getRMQConnectionString(conf config.RMQ) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.UserName, conf.Password, conf.Host, conf.Port)
}
