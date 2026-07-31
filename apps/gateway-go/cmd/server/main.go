package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	application := app.NewApp()

	http.HandleFunc("/healthz", application.HealthCheck)
	http.HandleFunc("/api/chat", application.HandleChat)

	fmt.Printf("Gateway service starting on port %s...\n", application.Cfg.Port)
	if err := http.ListenAndServe(":"+application.Cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
