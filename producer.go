package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Passo 1: estabeler a conexão com o RabbitMQ através do package amqp (protocolo de comunicação de mensagens: Advanced Message Protocol Queuing)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to RabbitMQ!")

	// Passo 2: criar um canal de comunicação com o servidor do RabbitMQ
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	fmt.Println("Channel successfully created!")

	// Passo 3: Declarar a fila
	q, err := ch.QueueDeclare(
		"notificationQueue", // Queue name
		true,                // Durability
		false,               // Exclusivity
		false,               // Auto-delete
		false,               // Do not wait for confirmations
		nil,                 // Aditional arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	fmt.Printf("Queue declarated successfully: %s\n", q.Name)

	// Passo 4: Publicar uma mensagem na fila
	body := "Notification test"
	err = ch.Publish(
		"",     // Exchange (default)
		q.Name, // Destination queue
		false,  // Do not use exclusive Exchange
		false,  // Do not use auto exchange
		amqp.Publishing{
			ContentType: "text/plain", // Content Type
			Body:        []byte(body), // Body Message
		},
	)
	if err != nil {
		log.Fatalf("failed to publish a message: %v", err)
	}

	fmt.Printf("Message sent: %s\n", body)
}
