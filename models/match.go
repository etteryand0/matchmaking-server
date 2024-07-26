package models

import "time"

type Match struct {
	ID        string `gorm:"primaryKey"`
	TestName  string
	Epoch     string `gorm:"size:36"`
	Score     int
	SessionID string
	CreatedAt time.Time
}
