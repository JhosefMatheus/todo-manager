package main

import (
	"log"
	"net/http"
	auth_controller "todo-manager/controllers/auth"
	token_controller "todo-manager/controllers/token"
	middleware "todo-manager/middlewares"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Erro ao carregar .env %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/sign-in", auth_controller.SignIn)
	mux.HandleFunc("/token/verify", token_controller.Verify)

	if err := http.ListenAndServe(":8080", middleware.GlobalMiddleware(mux)); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
