package match

import (
	"etteryand0/matchmaking/server/common"
	"etteryand0/matchmaking/server/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/perimeterx/marshmallow"
)

func LogMatch(c *gin.Context) {
	testName, foundTestNameParam := c.GetQuery("test_name")
	epoch, foundEpochParam := c.GetQuery("epoch")

	if !foundTestNameParam || !foundEpochParam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}
	if epoch == "last" {
		c.JSON(http.StatusBadRequest, gin.H{"Nostradamus": "No... no... no..."})
		return
	}
	if sessionTestName := c.GetString("session_test_name"); testName != sessionTestName {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wrong test name for this session"})
		return
	}

	testFile, err := os.Open(path.Join("tests", testName, "test.json"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	defer testFile.Close()

	byteValue, err := io.ReadAll(testFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while reading epoch file"})
		return
	}

	testCollection := TestCollection{}
	result, err := marshmallow.Unmarshal(byteValue, &testCollection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while unmarshalling test json file"})
		return
	}

	nextEpoch, ok := result[epoch]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Epoch does not exist"})
		return
	}

	if testCollection.Last == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Last epoch data does not exist"})
		return
	}

	if testCollection.Last == nextEpoch {
		session := models.Session{ID: c.GetString("session_id")}
		if err := session.Finish(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintln("DB error:", err.Error()),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"epoch":         nextEpoch,
		"is_last_epoch": testCollection.Last == nextEpoch,
	})
}

func GetResult(c *gin.Context) {
	var session models.Session
	result := common.DB.First(&session, "id = ?", c.GetString("session_id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintln("DB error:", result.Error.Error()),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Session does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"test_name":  session.TestName,
		"started_at": session.CreatedAt,
		"finished":   session.Finished,
	})
}
