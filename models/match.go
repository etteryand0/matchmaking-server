package models

import (
	"etteryand0/matchmaking/server/common"
	"time"
)

type Match struct {
	ID        string `gorm:"primaryKey"`
	TestName  string
	Epoch     string `gorm:"size:36"`
	Score     float64
	SessionID string
	CreatedAt time.Time
}

func CreateMatchBatch(matches *[]Match) error {
	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if result := tx.CreateInBatches(matches, 1000); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}
