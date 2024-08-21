package main

import (
	"medods-test/db"
	"medods-test/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db.Init()
	server := gin.Default()

	server.GET("/tokens/:id", handlers.GetTokens)
	server.POST("/refresh", handlers.Refresh)

	server.Run("localhost:8080")
}
