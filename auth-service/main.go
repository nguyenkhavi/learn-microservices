package main

import (
	"log"
	"nkvi/auth-service/middlewares"
	"nkvi/auth-service/models"

	authControllers "nkvi/auth-service/src/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file %s", err.Error())
	}
	models.Init()
	r := gin.Default()

	public := r.Group("/api")

	public.POST("/register", authControllers.Register)
	public.POST("/login", authControllers.Login)
	public.POST("/refresh-token", authControllers.Refresh)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", authControllers.CurrentUser)

	r.Run(":8080")

}
