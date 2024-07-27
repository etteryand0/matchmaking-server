package models

import (
	"etteryand0/matchmaking/server/common"
	"time"
)

type Session struct {
	ID           string `gorm:"primaryKey;size:36"`
	TestName     string
	Matches      []Match
	Finished     bool
	LeftOutUsers []User `gorm:"many2many:user_sessions"`
	WaitingTime  int
	CreatedAt    time.Time
	FinishedAt   time.Time
}

func (s *Session) Finish() error {
	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if result := tx.Model(&Session{}).Where("id = ?", s.ID).Updates(&Session{
		Finished:   true,
		FinishedAt: time.Now(),
	}); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func (s *Session) Create() error {
	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if result := tx.Create(&s); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	return tx.Commit().Error
}

func (s *Session) AddWaitingTime(time int) error {
	tx := common.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	result := tx.Exec("UPDATE sessions SET waiting_time = waiting_time + ? WHERE id = ?", time, s.ID)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Session) GetResult() (int, error) {
	score := 0
	result := common.DB.Model(&Match{}).Select("sum(score)").Where("session_id = ?", s.ID).Scan(&score)
	if result.Error != nil {
		return 0, nil
	}
	score += s.WaitingTime

	return score, nil
}
