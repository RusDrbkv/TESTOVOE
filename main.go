package main

import (
	"TESTOVOE/controllers"
	"TESTOVOE/events"
	"TESTOVOE/models"
	auth "TESTOVOE/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	router := gin.Default()

	models.ConnectDatabase()

	api := router.Group("/api")
	api.Use(auth.UserIdentity)

	{
		api.POST("/tasks", controllers.CreateTask)
		api.GET("/tasks", controllers.FindTasks)
		api.GET("/tasks/:id", controllers.FindTask)
		api.PATCH("/tasks/:id", controllers.UpdateTask)
		api.DELETE("/tasks/:id", controllers.DeleteTask)
	}
	router.POST("/signup", controllers.CreateUser)
	router.POST("/login", controllers.LoginUser)

	go events.KafkaConsumer()
	go events.KafkaProducer()

	// time.AfterFunc(3*time.Second, events.KafkaProducer)

	router.Run(":8080")

}
