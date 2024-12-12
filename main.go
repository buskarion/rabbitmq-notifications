package main

import (
	"log"
	"time"

	"github.com/buskarion/rabbitmq-notifications/consumer"
	"github.com/buskarion/rabbitmq-notifications/producer"
)

func main() {
	if err := producer.Start(); err != nil {
		log.Fatalf("Error starting the producer: %v", err)
	}

	time.Sleep(30 * time.Second)

	if err := consumer.Start(); err != nil {
		log.Fatalf("Error starting the consumer: %v", err)
	}

}
