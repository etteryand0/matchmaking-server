package main

import (
	"etteryand0/matchmaking/server/match"
	"etteryand0/matchmaking/server/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	r.GET("/matchmaking/users", users.GetWaitingUsers)

	r.POST("/matchmaking/match", match.LogMatch)

	r.Run()
}
