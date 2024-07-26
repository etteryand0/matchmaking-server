package models

import (
	"etteryand0/matchmaking/server/common"
	"time"
)

type Session struct {
	ID        string `gorm:"primaryKey;size:36"`
	TestName  string
	Matches   []Match
	Finished  bool
	CreatedAt time.Time
}

func (s *Session) Finish() error {
	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if result := tx.Model(&Session{}).Where("id = ?", s.ID).Update("finished", true); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}

func (s *Session) Create() error {
	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if result := tx.Create(&s); result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}
