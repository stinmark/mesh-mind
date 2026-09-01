package main

import (
	"gateway-go/internal/auth"
	"gateway-go/internal/db"
	"gateway-go/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/api/register", handlers.HandleRegister)
	http.HandleFunc("/api/login", handlers.HandleLogin)
	http.HandleFunc("/api/buy-tokens", auth.AuthMiddleware(handlers.HandleBuyTokens))
	http.HandleFunc("/api/chat", auth.AuthMiddleware(handlers.HandleChat))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway running on http://0.0.0.0:%s (InferenceTarget: %s)", port, handlers.InferenceURL)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
