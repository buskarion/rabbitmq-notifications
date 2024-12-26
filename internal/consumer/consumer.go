package consumer

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/buskarion/rabbitmq-notifications/internal/notification"
	rabbitmq "github.com/buskarion/rabbitmq-notifications/internal/queue"
	"github.com/streadway/amqp"
)

func ConsumeNotification(rmq *rabbitmq.RabbitMq, service *notification.Service) error {
	// consume messages from the queue
	msgs, err := rmq.Consume(
		rmq.Queue.Name,
		"",
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Error consuming messages: %w", err)
	}

	// processing messages
	for msg := range msgs {
		var notification notification.ScheduledNotification
		err := json.Unmarshal(msg.Body, &notification)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		retryCount := 0
		maxRetries := 3

		// process notification attempts
		for retryCount < maxRetries {
			err := processNotification()
			if err == nil {
				service.UpdateStatus(&notification, "processed")
				fmt.Printf("Notification \"%s\" status updated to %s\n", notification.Message, notification.Status)
				break
			} else {
				retryCount++
				fmt.Printf("Error processing message: %s, retrying... (%d/%d)\n", notification.Message, retryCount, maxRetries)
				time.Sleep(time.Duration(1<<retryCount) * time.Second) // exponential backoff
			}
		}

		// if retries limit get reached, move msg to DLQ
		if retryCount >= maxRetries {
			fmt.Printf("Failed to process the message \"%s\" after %d retries. Moving to DLQ.\n", notification.Message, retryCount)
			err := rmq.Channel.Publish(
				"",
				rmq.DLQ.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg.Body,
				},
			)
			if err != nil {
				fmt.Println("Error moving message to DLQ:", err)
			}
		}

		fmt.Print("Continuing with next message...\n\n")
	}

	return nil
}

// function that simulate a real notification process
func processNotification() error {
	// simulate a fail processing a msg at random
	if rand.IntN(2) == 0 {
		return fmt.Errorf("Simulated failure")
	}
	return nil
}
