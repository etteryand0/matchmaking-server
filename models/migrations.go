package models

import (
	"etteryand0/matchmaking/server/common"
	"fmt"
)

func MigrateDB() error {
	fmt.Println("Migrating database models")

	if err := common.DB.AutoMigrate(&Session{}); err != nil {
		return err
	}
	fmt.Println("- Migrated Session")

	if err := common.DB.AutoMigrate(&User{}); err != nil {
		return err
	}
	fmt.Println("- Migrated User")

	if err := common.DB.AutoMigrate(&Match{}); err != nil {
		return err
	}
	fmt.Println("- Migrated Match")

	if err := common.DB.AutoMigrate(&Epoch{}); err != nil {
		return err
	}
	fmt.Println("- Migrated Epoch")
	fmt.Println("Successfuly migrated models")
	return nil
}
