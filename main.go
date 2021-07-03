package main

import (
	"api/controllers"
	"api/middlewares"
	"api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()

	router := gin.Default()

	public := router.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := router.Group("/api/admin")

	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)

	router.Run(":8080")
}
