package consumer

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func Start() {
	// RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error trying to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Create communication channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %s", err)
	}
	defer ch.Close()

	// Declare a queue where the messages will be received
	q, err := ch.QueueDeclare(
		"notificationQueue", // queue name
		true,                // durable: the queue will remain after reboot
		false,               // autoDelete: the queue will not be deleted
		false,               // exclusive: the queue can be shared
		false,               // noWait: we will not wait for a RabbitMQ response
		nil,                 // additional arguments
	)
	if err != nil {
		log.Fatalf("error declaring queue: %s", err)
	}

	// Consume queue messages
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer do not need a name
		true,   // autoAck: auto confirm the message reception
		false,  // exclusive: the consumer can be shared
		false,  // noLocal: do not consume messages sent by the own channel
		false,  // noWait: do not wait for response
		nil,    // additional arguments
	)
	if err != nil {
		log.Fatalf("error consumming messages: %s", err)
	}

	// Loop to proccess messages
	for d := range msgs {
		fmt.Printf("Message received: %s\n", d.Body)
	}
}
