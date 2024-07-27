package models

type Epoch struct {
	Epoch    string
	Position int `gorm:"primaryKey"`
	Interval int
	TestName string `gorm:"primaryKey"`
}
