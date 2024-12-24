package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMq struct {
	Connection      *amqp.Connection
	Channel         *amqp.Channel
	Queue           amqp.Queue
	DeadLetterQueue amqp.Queue
}

func New(queueName, dlqName string) (*RabbitMq, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error creating the channel: %w", err)
	}

	// declare dead-letter queue
	dlq, err := ch.QueueDeclare(
		dlqName,
		true,  // durable
		false, // delete when unsued
		false, // exclusive
		false, // no-wait
		nil,   // arguments
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
			"x-dead-letter-routing-key": dlq.Name,     // send failed messages to DLQ
			"x-message-ttl":             int32(60000), // time-to-live = 1 minute
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error declaring the queue: %w", err)
	}

	return &RabbitMq{
		Connection:      conn,
		Channel:         ch,
		Queue:           queue,
		DeadLetterQueue: dlq,
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
