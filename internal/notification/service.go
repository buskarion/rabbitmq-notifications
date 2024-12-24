package notification

import "time"

func CreateNotification(message string, sendAt time.Time) *ScheduledNotification {
	return &ScheduledNotification{
		Message: message,
		SendAt:  sendAt,
		Status:  "scheduled",
	}
}

func UpdateStatus(notification *ScheduledNotification, status string) {
	notification.Status = status
}
