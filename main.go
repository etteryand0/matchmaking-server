package main

import (
	"etteryand0/matchmaking/server/common"
	"etteryand0/matchmaking/server/match"
	"etteryand0/matchmaking/server/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	common.ConnectDatabase()
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	r.GET("/matchmaking/users", users.GetWaitingUsers)

	protected := r.Group("/matchmaking", match.AuthMiddleware())

	protected.POST("/match", match.LogMatch)
	protected.POST("/result", match.GetResult)

	r.Run()
}
