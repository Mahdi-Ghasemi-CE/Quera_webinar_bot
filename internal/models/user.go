package models

import "time"

type User struct {
	ID          uint  `gorm:"primaryKey"`
	TelegramID  int64 `gorm:"uniqueIndex"`
	CreatedAt   time.Time
	FirstName   string
	LastName    string
	PhoneNumber string
}
