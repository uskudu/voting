package main

import (
	"voting/initializers"
	"voting/internal/db"
	"voting/internal/poll"

	"github.com/gin-gonic/gin"
)

func main() {
	initializers.LoadEnvVars()
	database, err := db.Init()
	if err != nil {
		panic(err)
	}

	pollRepo := poll.NewPollRepository(database)
	pollService := poll.NewPollService(pollRepo)
	pollHandlers := poll.NewPollHandler(pollService)

	r := gin.Default()

	r.POST("/polls", pollHandlers.PostPoll)
	r.GET("/polls", pollHandlers.GetPolls)
	r.GET("/polls/:id", pollHandlers.GetPoll)
	r.PATCH("/polls/:id", pollHandlers.PatchPoll)
	r.DELETE("/polls/:id", pollHandlers.DeletePoll)

	r.Run(":8080")
}
