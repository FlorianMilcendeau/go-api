package main

import (
	"api/controllers"
	"api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()

	router := gin.Default()

	public := router.Group("/api")

	public.POST("/register", controllers.Register)

	router.Run(":8080")
}
