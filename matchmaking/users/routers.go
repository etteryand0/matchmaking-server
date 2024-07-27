package users

import (
	"encoding/json"
	"etteryand0/matchmaking/server/common"
	"etteryand0/matchmaking/server/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetWaitingUsers(c *gin.Context) {
	testName, _ := c.GetQuery("test_name")
	epoch, _ := c.GetQuery("epoch")

	var users []models.User
	result := common.DB.Where("test_name = ? AND epoch = ?", testName, epoch).Select(
		"id",
		"mmr",
		"waiting_time",
		"roles",
	).Find(&users)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, "File not found")
		return
	}
	totalEpochWaitingTime := 0
	var response []User
	for _, user := range users {
		totalEpochWaitingTime += user.WaitingTime
		var roles []string
		err := json.Unmarshal([]byte(user.Roles), &roles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "couldn't unmarshal roles json",
			})
			return
		}
		response = append(response, User{
			UserId:      user.ID,
			MMR:         user.MMR,
			Roles:       roles,
			WaitingTime: user.WaitingTime,
		})
	}

	if epoch == "00000000-0000-0000-0000-000000000000" {
		UUID := uuid.New().String()
		session := models.Session{
			ID:          UUID,
			TestName:    testName,
			Finished:    false,
			WaitingTime: totalEpochWaitingTime,
		}
		if err := session.Create(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintln("DB error:", err.Error()),
			})
			return
		}
		c.SetCookie("session", UUID, 60*10, "/", "localhost", false, true)
	} else {
		session := models.Session{
			ID: c.GetString("session_id"),
		}
		session.AddWaitingTime(totalEpochWaitingTime)
	}

	c.JSON(http.StatusOK, response)
}
