package main

import (
	"log"
	authcontroller "todo-manager/controllers/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Erro ao carregar .env %v", err)
	}

	router := gin.Default()

	router.POST("/auth/sign-in", authcontroller.SignIn)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
