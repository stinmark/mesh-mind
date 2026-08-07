// package internal

// handlers, DB queries, gRPC clients

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Port         string
	InferenceURL string
}

type App struct {
	Cfg        Config
	HTTPClient *http.Client
}

func NewApp() *App {
	inferenceURL := os.Getenv("INFERENCE_SERVICE_URL")
	if inferenceURL == "" {
		inferenceURL = "http://localhost:8000"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &App{
		Cfg: Config{
			Port:         port,
			InferenceURL: inferenceURL,
		},
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request/Response structs matching the Python API
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages     []ChatMessage `json:"messages"`
	MaxNewTokens int           `json:"max_new_tokens,omitempty"`
	Temperature  float64       `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
}

// Handler for incoming web/frontend requests
func (a *App) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON body: %v", err), http.StatusBadRequest)
		return
	}

	// Forward request to Python inference service
	jsonPayload, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to serialize request", http.StatusInternalServerError)
		return
	}

	targetURL := fmt.Sprintf("%s/v1/chat/completions", a.Cfg.InferenceURL)
	resp, err := a.HTTPClient.Post(targetURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reach inference service: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from inference service", http.StatusInternalServerError)
		return
	}

	// Forward response back to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"gateway ok"}`))
}
