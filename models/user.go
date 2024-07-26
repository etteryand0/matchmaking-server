package models

import "time"

type User struct {
	ID          string `gorm:"primaryKey;size:36"`
	TestName    string
	Epoch       string `gorm:"size:36"`
	WaitingTime uint
	Roles       string
	CreatedAt   time.Time
}
