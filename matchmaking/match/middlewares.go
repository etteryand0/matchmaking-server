package match

import (
	"etteryand0/matchmaking/server/common"
	"etteryand0/matchmaking/server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		testName, foundTestNameParam := c.GetQuery("test_name")
		epoch, foundEpochParam := c.GetQuery("epoch")

		if !foundTestNameParam || !foundEpochParam {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
			c.Abort()
			return
		}
		if epoch == "last" {
			c.JSON(http.StatusBadRequest, gin.H{"Nostradamus": "No... no... no..."})
			c.Abort()
			return
		}

		session_id, err := c.Cookie("session")
		if err != nil {
			// Cookie verification failed
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden with no cookie"})
			c.Abort()
			return
		}
		var session models.Session
		result := common.DB.Select("test_name").First(&session, "id = ?", session_id)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Session does not exist"})
			c.Abort()
			return
		}

		if testName != session.TestName {
			c.JSON(http.StatusForbidden, gin.H{"error": "Wrong test name for this session"})
			c.Abort()
			return
		}
		c.Set("session_id", session_id)
		c.Set("session_test_name", session.TestName)
		c.Next()
	}
}
