package common

import (
	"fmt"

	"gorm.io/driver/postgres"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() error {
	fmt.Println("Connecting to database")
	// db, err := gorm.Open(sqlite.Open("server.db"), &gorm.Config{})
	dsn := "host=db user=server password=server dbname=server port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)

	DB = db
	return nil
}
