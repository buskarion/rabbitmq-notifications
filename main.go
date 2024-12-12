package main

import (
	"github.com/buskarion/rabbitmq-notifications/consumer"
	"github.com/buskarion/rabbitmq-notifications/producer"
)

func main() {
	go producer.Start()
	consumer.Start()
}
