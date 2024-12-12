package notification

import "time"

type ScheduledNotification struct {
	Message string    `json:"message"`
	SendAt  time.Time `json:"send_at"`
}
