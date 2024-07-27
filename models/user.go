package models

import "time"

type User struct {
	ID            string `gorm:"primaryKey;size:36"`
	TestName      string `gorm:"primaryKey"`
	MMR           int
	Epoch         string `gorm:"size:36"`
	EpochPosition int
	WaitingTime   int
	Roles         string
	Sessions      []Session `gorm:"many2many:user_sessions"`
	CreatedAt     time.Time
}
