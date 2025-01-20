package main

import (
	authcontroller "todo-manager/controllers/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/auth/sign-in", authcontroller.SignIn)

	router.Run("localhost:8080")
}
