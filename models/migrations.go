package models

import "etteryand0/matchmaking/server/common"

func MigrateDB() {
	common.DB.AutoMigrate(&Session{})
	common.DB.AutoMigrate(&User{})
	common.DB.AutoMigrate(&Match{})
}
