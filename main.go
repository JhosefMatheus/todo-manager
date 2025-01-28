package main

import (
	"log"
	"net/http"
	authcontroller "todo-manager/controllers/auth"
	middleware "todo-manager/middlewares"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Erro ao carregar .env %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/sign-in", authcontroller.SignIn)

	if err := http.ListenAndServe(":8080", middleware.GlobalMiddleware(mux)); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
