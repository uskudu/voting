package main

import (
	"log"
	"voting/initializers"
	"voting/internal/db"
	"voting/internal/middleware"
	"voting/internal/poll"
	"voting/internal/user"
	"voting/notifications/ws"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"

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

	//// rabbitmq
	//conn, err := amqp091.Dial("amqp://localhost:5672")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//ch, err := conn.Channel()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_, err = ch.QueueDeclare("VoteNotifications", true, false, false, false, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	pollRepo := poll.NewPollRepository(database)
	pollService := poll.NewPollService(pollRepo)
	pollHandlers := poll.NewPollHandler(pollService, ch)

	userRepo := user.NewUserRepository(database)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ws", ws.WebSocketHandler)

	r.POST("/polls", middleware.RequireAuth, pollHandlers.PostPoll)
	r.GET("/polls", pollHandlers.GetPolls)
	r.GET("/polls/:id", pollHandlers.GetPoll)
	r.PATCH("/polls/:id", middleware.RequireAuth, pollHandlers.PatchPoll)
	r.DELETE("/polls/:id", middleware.RequireAuth, pollHandlers.DeletePoll)
	r.POST("/polls/:id/options/:optionID/vote", middleware.RequireAuth, pollHandlers.Vote)

	r.POST("/users", userHandler.PostUser)
	r.GET("/users", userHandler.GetUsers)
	r.GET("/users/:id", userHandler.GetUser)
	r.PATCH("/users/:id", middleware.RequireAuth, userHandler.PatchUser)
	r.DELETE("/users/:id", middleware.RequireAuth, userHandler.DeleteUser)

	r.POST("/login", userHandler.Login)
	r.GET("/validate", middleware.RequireAuth, userHandler.Validate)

	r.Run(":8080")
}
