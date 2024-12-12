package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/buskarion/rabbitmq-notifications/notification"
	"github.com/buskarion/rabbitmq-notifications/rabbitmq"
)

func Start() error {
	// Create a rabbitmq instance
	rmq, err := rabbitmq.New("notificationQueue")
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

		if time.Now().After(notification.SendAt) {
			fmt.Printf("Message received: %s\n", msg.Body)
		} else {
			fmt.Printf("Message \"%s\" is scheduled to be processed at %s, skipping.\n", notification.Message, notification.SendAt)
		}
	}

	return nil
}
