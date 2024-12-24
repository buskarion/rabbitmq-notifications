package producer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buskarion/rabbitmq-notifications/internal/notification"
	"github.com/buskarion/rabbitmq-notifications/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func Start() error {
	// Create rabbitmq instance
	rmq, err := rabbitmq.New("notificationQueue", "notificationDLQ")
	if err != nil {
		return err
	}
	defer rmq.Close()

	// create scheduled rabbitmq-notifications
	notifications := []notification.ScheduledNotification{
		{
			Message: "Message 1",
			SendAt:  time.Now().Add(10 * time.Second),
		},
		{
			Message: "Message 2",
			SendAt:  time.Now().Add(60 * time.Second),
		},
	}

	// proccess notifications to send to the queue
	for _, n := range notifications {
		sendToQueue(n, rmq.Channel, rmq.Queue.Name)
	}

	return nil
}

func sendToQueue(notification notification.ScheduledNotification, ch *amqp.Channel, queueName string) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error marshalling notification: %w", err)
	}
	fmt.Printf("Sending message %s to the queue\n", notification.Message)

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("Error sending message to queue: %w", err)
	}

	fmt.Printf("Message \"%s\" sent to the queue \"%s\" at %s\n", body, queueName, time.Now())
	return nil
}
