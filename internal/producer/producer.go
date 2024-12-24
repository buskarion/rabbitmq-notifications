package producer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buskarion/rabbitmq-notifications/internal/notification"
	"github.com/buskarion/rabbitmq-notifications/internal/rabbitmq"
	"github.com/streadway/amqp"
)

func SendNotification(
	notification *notification.ScheduledNotification,
	rmq *rabbitmq.RabbitMq,
) error {
	// verify if the notification is on time be sent
	if time.Now().Before(notification.SendAt) {
		fmt.Printf("Notification for \"%s\" is scheduled to be sent at: %s\n", notification.Message, notification.SendAt)
		time.Sleep(time.Until(notification.SendAt))
	}

	notification.Status = "sent"
	fmt.Printf("Notification \"%s\" status updated to sent at: %s\n", notification.Message, notification.SendAt)

	// Publish the message on the queue
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error marshilling notification: %w", err)
	}

	err = rmq.Channel.Publish(
		"",
		rmq.Queue.Name,
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

	fmt.Printf("Message \"%s\" sent to the queue \"%s\" at %s\n", body, rmq.Queue.Name, time.Now())
	return nil
}
