package storage

import "time"

// User represents a Telegram user
type User struct {
	TelegramID    int64
	CheckInterval int // in minutes
	CreatedAt     time.Time
}

// MonitoredStreet represents a street being monitored by a user
type MonitoredStreet struct {
	ID         int64
	TelegramID int64
	StreetName string
}

// SeenCall represents a call that has been seen by a user
type SeenCall struct {
	TelegramID int64
	IncidentID string
	SeenAt     time.Time
}

// AlertCall represents a call to be sent as an alert
type AlertCall struct {
	IncidentID string
	CallType   string
	Nature     string
	Address    string
	StartTime  time.Time
	Agency     string
}
