package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"mesh-mind/gateway/internal/auth"
	"mesh-mind/gateway/internal/db"
	"mesh-mind/gateway/internal/handlers"
)

func main() {
	// 1. Initialize Server-Side SQLite Database
	dbPath := getEnv("DB_PATH", "./data/gateway.db")
	database, err := db.InitGatewayDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize gateway DB: %v", err)
	}
	defer database.Close()

	// 2. Instantiate Handlers
	authHandler := &auth.AuthHandler{DB: database}

	// 3. Setup HTTP Routes
	mux := http.NewServeMux()

	// Public Auth Endpoints
	mux.HandleFunc("POST /api/v1/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/login", authHandler.Login)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Protected Prediction Route (Secured with JWT Middleware)
	mux.HandleFunc("POST /api/v1/predict", auth.AuthenticateMiddleware(handlers.PredictHandler))

	// 4. Start HTTP Server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMiddleware(mux),
	}

	// 5. Graceful Shutdown
	go func() {
		log.Printf("🚀 Mesh-Mind Gateway listening on http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Stopping Gateway gracefully...")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
