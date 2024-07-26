package users

import (
	"encoding/json"
	"etteryand0/matchmaking/server/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetWaitingUsers(c *gin.Context) {
	testName, foundTestNameParam := c.GetQuery("test_name")
	epoch, foundEpochParam := c.GetQuery("epoch")

	if !foundTestNameParam || !foundEpochParam {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	testFile, err := os.Open(path.Join("tests", testName, epoch+".json"))
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

	var users []User
	if err := json.Unmarshal(byteValue, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while unmarshalling test json file"})
		return
	}

	if epoch == "00000000-0000-0000-0000-000000000000" {
		UUID := uuid.New().String()
		session := models.Session{
			ID:       UUID,
			TestName: testName,
			Finished: false,
		}
		if err := session.Create(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintln("DB error:", err.Error()),
			})
			return
		}
		c.SetCookie("session", UUID, 60*10, "/", "localhost", false, true)
	}
	c.JSON(http.StatusOK, users)
}
