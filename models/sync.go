package models

import (
	"encoding/json"
	"errors"
	"etteryand0/matchmaking/server/common"
	"fmt"
	"os"
	"path"

	"gorm.io/gorm/clause"
)

func SyncDB() error {
	fmt.Println("Syncing tests folder with SQLite database, please wait...")
	dirs, err := os.ReadDir("tests")
	if err != nil {
		return errors.New("error reading tests directory")
	}
	var usersToCommit []User
	var epochesToCommit []Epoch
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		testName := dir.Name()
		var intervals map[string]int
		err = common.ReadJson(path.Join("tests", testName, "intervals.json"), &intervals)
		if err != nil {
			return err
		}
		var testMap map[string]string
		err = common.ReadJson(path.Join("tests", testName, "test.json"), &testMap)
		if err != nil {
			return err
		}

		epoch := "00000000-0000-0000-0000-000000000000"
		position := 0
		lastEpoch, ok := testMap["last"]
		if !ok {
			return errors.New("last epoch data doesn't exist")
		}
		fmt.Println("- Syncing", testName)
		fmt.Println("last", lastEpoch)
		for {
			epochesToCommit = append(epochesToCommit, Epoch{
				Epoch:    epoch,
				Position: position,
				Interval: intervals[epoch],
				TestName: testName,
			})

			var users []struct {
				UserId      string   `json:"user_id"`
				MMR         int      `json:"mmr"`
				Roles       []string `json:"roles"`
				WaitingTime int      `json:"waitingTime"`
			}
			err = common.ReadJson(path.Join("tests", testName, epoch+".json"), &users)
			if err != nil {
				return err
			}
			for _, user := range users {
				byteValue, err := json.Marshal(user.Roles)
				if err != nil {
					return errors.New("error marshalling roles")
				}
				usersToCommit = append(usersToCommit, User{
					TestName:      testName,
					Epoch:         epoch,
					Roles:         string(byteValue),
					ID:            user.UserId,
					MMR:           user.MMR,
					WaitingTime:   intervals[epoch] + user.WaitingTime,
					EpochPosition: position,
				})
			}
			if lastEpoch == epoch {
				break
			}
			epoch, ok = testMap[epoch]
			position += 1
			fmt.Println("epoch", epoch)
			if !ok {
				return errors.New("invalid test data")
			}
		}
	}

	tx := common.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return errors.New(fmt.Sprintln("DB error", tx.Error.Error()))
	}

	if res := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&usersToCommit, 1000); res.Error != nil {
		return errors.New(fmt.Sprintln("DB error", res.Error.Error()))
	}
	fmt.Println("epoch commit", epochesToCommit)
	if res := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&epochesToCommit, 1000); res.Error != nil {
		return errors.New(fmt.Sprintln("DB error", res.Error.Error()))
	}
	tx.Commit()

	fmt.Println("Successfuly synced users")
	return nil
}
