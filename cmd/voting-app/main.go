package main

import (
	"voting/initializers"
	"voting/internal/db"
	"voting/internal/poll"
	"voting/internal/user"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "voting/docs"
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

	userRepo := user.NewUserRepository(database)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/polls", pollHandlers.PostPoll)
	r.GET("/polls", pollHandlers.GetPolls)
	r.GET("/polls/:id", pollHandlers.GetPoll)
	r.PATCH("/polls/:id", pollHandlers.PatchPoll)
	r.DELETE("/polls/:id", pollHandlers.DeletePoll)

	r.POST("/users", userHandler.PostUser)
	r.GET("/users", userHandler.GetUsers)
	r.GET("/users/:id", userHandler.GetUser)
	r.PATCH("/users/:id", userHandler.PatchUser)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	r.Run(":8080")
}
