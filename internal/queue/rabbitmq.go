package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	DLQ        amqp.Queue
}

func New(queueName string) (*RabbitMq, error) {
	// create a connection with RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error creating the channel: %w", err)
	}

	// Declare DLQ
	dlq, err := ch.QueueDeclare(
		queueName+"_dlq",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declaring dead-letter queue: %w", err)
	}

	// Declare main queue
	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": dlq.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error declaring main queue: %w", err)
	}

	return &RabbitMq{
		Connection: conn,
		Channel:    ch,
		Queue:      queue,
		DLQ:        dlq,
	}, nil
}

func (r *RabbitMq) Consume(
	queueName string, consumerTag string,
	autoAck, exclusive, noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(queueName, consumerTag, autoAck, exclusive, noWait, false, args)
}

func (r *RabbitMq) Publish(
	exchange, routingKey string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) error {
	return r.Channel.Publish(exchange, routingKey, mandatory, immediate, msg)
}

func (r *RabbitMq) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Printf("Error closing the channel: %v", err)
	}
	if err := r.Connection.Close(); err != nil {
		log.Printf("Error closing the Connection: %v", err)
	}
}
