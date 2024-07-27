package users

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

		if epoch != "00000000-0000-0000-0000-000000000000" {
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
			if session.TestName != testName {
				c.JSON(http.StatusForbidden, gin.H{"error": "Wrong test name for this session"})
				c.Abort()
				return
			}
			c.Set("session_id", session_id)
		}

		c.Next()
	}
}
