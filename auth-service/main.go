package main

import (
	"nkvi/auth-service/middlewares"
	"nkvi/auth-service/models"
	authControllers "nkvi/auth-service/src/auth"

	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDataBase()

	r := gin.Default()

	public := r.Group("/api")

	public.POST("/register", authControllers.Register)
	public.POST("/login", authControllers.Login)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", authControllers.CurrentUser)

	r.Run(":8080")

}
