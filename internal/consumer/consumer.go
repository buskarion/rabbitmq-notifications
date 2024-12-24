package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buskarion/rabbitmq-notifications/internal/notification"
	"github.com/buskarion/rabbitmq-notifications/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func Start() error {
	// Create a rabbitmq instance
	rmq, err := rabbitmq.New("notificationQueue", "notificationDLQ")
	if err != nil {
		return err
	}
	defer rmq.Close()

	// Consuming queue messages
	msgs, err := rmq.Channel.Consume(
		rmq.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Error consumming messages: %w", err)
	}

	// Loop to proccess messages
	fmt.Println("\nConsumer is waiting for messages...")
	for msg := range msgs {
		var notification notification.ScheduledNotification
		err := json.Unmarshal(msg.Body, &notification)
		if err != nil {
			fmt.Println("error unmarshalling message:", err)
			continue
		}

		retryCount := 0
		maxRetries := 3

		// retry logic
		for retryCount < maxRetries {
			if time.Now().After(notification.SendAt) {
				// Process the notification
				fmt.Printf("Message received: %s\n", msg.Body)
			} else {
				retryCount++
				fmt.Printf("Retrying message \"%s\" (%d/%d)\n", notification.Message, retryCount, maxRetries)
				time.Sleep(time.Duration(retryCount) * time.Second) // exponential backoff
			}
		}

		// if we exhausted retries, move the message to notificationDLQ
		if retryCount >= maxRetries {
			fmt.Printf("Message \"%s\" failed after %d retries. Moving to DLQ.\n", notification.Message, maxRetries)
			err := rmq.Channel.Publish(
				"",
				rmq.DeadLetterQueue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg.Body,
				},
			)
			if err != nil {
				fmt.Println("Error moving message do DLQ:", err)
			}
		}

		if time.Now().After(notification.SendAt) {
			fmt.Printf("Message received: %s\n", msg.Body)
		} else {
			fmt.Printf("Message \"%s\" is scheduled to be processed at %s, skipping.\n", notification.Message, notification.SendAt)
		}
	}

	return nil
}
