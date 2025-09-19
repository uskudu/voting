package main

import (
	"voting/initializers"
	"voting/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	initializers.LoadEnvVars()
	database, err := db.Init()
	if err != nil {
		panic(err)
	}
	var _ = database
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200,
			gin.H{
				"message": "pong",
			})
	})

	r.Run(":8080")
}
