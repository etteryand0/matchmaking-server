package common

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("server.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintln("Database connection failed", err.Error()))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintln("Database connection failed", err.Error()))
	}

	sqlDB.SetMaxIdleConns(10)

	DB = db
}
