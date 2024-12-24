package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/buskarion/rabbitmq-notifications/internal/notification"
	"github.com/buskarion/rabbitmq-notifications/internal/rabbitmq"
)

func ConsumeNotification(rmq *rabbitmq.RabbitMq) error {
	// consume messages from the queue
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

		// process notification
		fmt.Printf("Processing message: %s\n", notification.Message)
		notification.Status = "processed"
		fmt.Printf("Notification \"%s\"  status updated to %s\n", notification.Message, notification.Status)
	}

	return nil
}
