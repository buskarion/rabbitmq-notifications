package main

import (
	"log"
	"time"

	"github.com/buskarion/rabbitmq-notifications/internal/consumer"
	"github.com/buskarion/rabbitmq-notifications/internal/notification"
	"github.com/buskarion/rabbitmq-notifications/internal/producer"
	"github.com/buskarion/rabbitmq-notifications/internal/rabbitmq"
)

func main() {
	// start notification service
	notificationService := &notification.Service{}

	// setup a rabbitMQ instance
	rmq, err := rabbitmq.New("notificationQueue")
	if err != nil {
		log.Fatalf("Error initializing RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// create new notification
	notifications := []*notification.ScheduledNotification{
		notificationService.CreateNotification("M_ONE", time.Now().Add(5*time.Second)),
		notificationService.CreateNotification("M_TWO", time.Now().Add(10*time.Second)),
		notificationService.CreateNotification("M_THREE", time.Now().Add(20*time.Second)),
	}

	// send notifications with producer
	for _, n := range notifications {
		err := producer.SendNotification(n, rmq, notificationService)
		if err != nil {
			log.Fatalf("Error sending notification: %v", err)
		}

		// simulates the time to send a notification to the queue
		time.Sleep(2 * time.Second)
	}

	log.Println("\nAll notifications sent. \nStarting to consume notifications...")

	// consume notifications with consumer
	err = consumer.ConsumeNotification(rmq, notificationService)
	if err != nil {
		log.Fatalf("Error consumming notifications: %v", err)
	}

	log.Println("All messages consumed.")

}
