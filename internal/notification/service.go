package notification

import "time"

type Service struct {
}

func (s *Service) CreateNotification(message string, sendAt time.Time) *ScheduledNotification {
	return &ScheduledNotification{
		Message: message,
		SendAt:  sendAt,
		Status:  "scheduled",
	}
}

func (s *Service) pdateStatus(notification *ScheduledNotification, status string) {
	notification.Status = status
}
