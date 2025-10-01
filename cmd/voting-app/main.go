package main

import (
	"voting/initializers"
	"voting/internal/db"
	"voting/internal/poll"
	"voting/internal/user/crud"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "voting/docs"
)

// @title Voting API
// @version 1.0
// @description API for polls and users management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@voting.local

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	initializers.LoadEnvVars()
	database, err := db.Init()
	if err != nil {
		panic(err)
	}

	pollRepo := poll.NewPollRepository(database)
	pollService := poll.NewPollService(pollRepo)
	pollHandlers := poll.NewPollHandler(pollService)

	userRepo := crud.NewUserRepository(database)
	userService := crud.NewUserService(userRepo)
	userHandler := crud.NewUserHandler(userService)

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
