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

	r.POST("/poll", pollHandlers.PostPoll)
	r.GET("/polls", pollHandlers.GetPolls)
	r.GET("/poll/:id", pollHandlers.GetPoll)
	r.PATCH("/poll/:id", pollHandlers.PatchPoll)
	r.DELETE("/poll/:id", pollHandlers.DeletePoll)

	r.Run(":8080")
}
