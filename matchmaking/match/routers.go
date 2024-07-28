package match

import (
	"etteryand0/matchmaking/server/common"
	"etteryand0/matchmaking/server/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveMatch(c *gin.Context) {
	testName, _ := c.GetQuery("test_name")
	epoch, _ := c.GetQuery("epoch")

	var matchDatas []Match
	if err := c.ShouldBindJSON(&matchDatas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect data"})
		return
	}

	var epochModel models.Epoch
	result := common.DB.Model(&models.Epoch{}).Where("epoch = ? AND test_name = ?", epoch, testName).First(&epochModel)
	if result.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprint("DB error: ", result.Error.Error())},
		)
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Epoch does not exist"})
		return
	}

	var matchesToCommit []models.Match
	for _, matchData := range matchDatas {
		score, err := matchData.CalculateScore(testName, epochModel.Position)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		matchesToCommit = append(matchesToCommit, models.Match{
			ID:        matchData.MatchId,
			TestName:  testName,
			Epoch:     epoch,
			Score:     score,
			SessionID: c.GetString("session_id"),
		})
	}

	if err := models.CreateMatchBatch(&matchesToCommit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprint("DB error: ", err.Error()),
		})
		return
	}

	var nextEpoch models.Epoch
	result = common.DB.Where("position = ? AND test_name = ?", epochModel.Position+1, testName).Limit(1).Find(&nextEpoch)
	if result.Error != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprint("DB error: ", result.Error.Error())},
		)
		return
	}
	if result.RowsAffected == 0 {
		session := models.Session{ID: c.GetString("session_id")}
		if err := session.Finish(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint("DB error: ", err.Error()),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"new_epoch":     nil,
			"is_last_epoch": true,
		})
		return
	}

	isLast := false
	result = common.DB.Model(&nextEpoch).Select("1").Where("position = ? AND test_name = ?", epochModel.Position+2, testName).Limit(1).Find(&nextEpoch)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintln("DB error:", result.Error.Error()),
		})
		return
	}
	if result.RowsAffected == 0 {
		isLast = true
	}

	c.JSON(http.StatusOK, gin.H{
		"new_epoch":     nextEpoch.Epoch,
		"is_last_epoch": isLast,
	})
}

func GetResult(c *gin.Context) {
	var session models.Session
	result := common.DB.First(&session, "id = ?", c.GetString("session_id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprint("DB error: ", result.Error.Error()),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Session does not exist",
		})
		return
	}

	score, err := session.GetResult()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":   session.ID,
		"test_name":    session.TestName,
		"score":        score,
		"started_at":   session.CreatedAt,
		"time_elapsed": fmt.Sprint(session.FinishedAt.Sub(session.CreatedAt)),
		"finished":     session.Finished,
	})
}
