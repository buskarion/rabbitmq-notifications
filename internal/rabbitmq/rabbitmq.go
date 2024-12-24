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
}

func New(queueName string) (*RabbitMq, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error creating the channel: %w", err)
	}

	// Declare main queue
	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declaring the queue: %w", err)
	}

	return &RabbitMq{
		Connection: conn,
		Channel:    ch,
		Queue:      queue,
	}, nil
}

func (r *RabbitMq) Close() {
	if err := r.Channel.Close(); err != nil {
		log.Printf("Error closing the channel: %v", err)
	}
	if err := r.Connection.Close(); err != nil {
		log.Printf("Error closing the Connection: %v", err)
	}
}
