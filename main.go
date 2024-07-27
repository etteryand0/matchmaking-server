package main

import (
	"etteryand0/matchmaking/server/common"
	"etteryand0/matchmaking/server/matchmaking/match"
	"etteryand0/matchmaking/server/matchmaking/users"
	"etteryand0/matchmaking/server/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := common.ConnectDatabase(); err != nil {
		fmt.Println("Database connection failed", err.Error())
		return
	}
	if err := models.MigrateDB(); err != nil {
		fmt.Println("Models migration failed:", err.Error())
		return
	}
	if err := models.SyncDB(); err != nil {
		fmt.Println("Error while syncing users:", err.Error())
		return
	}
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/matchmaking/users", users.GetWaitingUsers)

	protected := r.Group("/matchmaking", match.AuthMiddleware())

	protected.POST("/match", match.SaveMatch)
	protected.POST("/result", match.GetResult)

	r.Run()
}
