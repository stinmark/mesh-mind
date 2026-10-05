package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type PredictRequest struct {
	Text  string `json:"text"`
	Model string `json:"model,omitempty"`
}

// Struct sent to Python worker
type WorkerRequest struct {
	Text string `json:"text"`
}

// Struct returned by Python worker
type WorkerResponse struct {
	Text       string  `json:"text"`
	ModelUsed  string  `json:"model_used"`
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
	LatencyMs  float64 `json:"latency_ms"`
}

type PredictResponse struct {
	Text      string  `json:"text"`
	ModelUsed string  `json:"model_used"`
	Intent    string  `json:"intent,omitempty"`
	Sentiment string  `json:"sentiment,omitempty"`
	LatencyMs float64 `json:"latency_ms"`
}

func PredictHandler(w http.ResponseWriter, r *http.Request) {
	var req PredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Text) == "" {
		http.Error(w, "Text field is required", http.StatusBadRequest)
		return
	}

	targetModel := strings.ToLower(req.Model)
	if targetModel == "" || targetModel == "auto" {
		targetModel = routeAuto(req.Text)
	}

	// Dispatch to actual ONNX microservice worker
	res, err := dispatchToWorker(targetModel, req.Text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Model worker error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func routeAuto(text string) string {
	lowerText := strings.ToLower(text)
	if strings.Contains(lowerText, "how") || strings.Contains(lowerText, "why") || strings.Contains(lowerText, "?") {
		return "intent"
	}
	return "sentiment"
}

func dispatchToWorker(model string, text string) (PredictResponse, error) {
	start := time.Now()

	// Determine endpoint based on model selection
	workerURL := os.Getenv("SENTIMENT_WORKER_URL")
	if model == "intent" {
		workerURL = os.Getenv("INTENT_WORKER_URL")
	}

	// Fallback safety if env vars aren't set
	if workerURL == "" {
		if model == "intent" {
			workerURL = "http://model-intent:5000/predict"
		} else {
			workerURL = "http://model-sentiment:5000/predict"
		}
	}

	// Prepare JSON payload for Python FastAPI worker
	payload, err := json.Marshal(WorkerRequest{Text: text})
	if err != nil {
		return PredictResponse{}, err
	}

	// Make HTTP POST request to Python worker
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(workerURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return PredictResponse{}, fmt.Errorf("failed to reach worker at %s: %w", workerURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PredictResponse{}, fmt.Errorf("worker returned status %d", resp.StatusCode)
	}

	var workerRes WorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&workerRes); err != nil {
		return PredictResponse{}, err
	}

	totalLatency := float64(time.Since(start).Microseconds()) / 1000.0

	// Format response cleanly based on model type
	finalRes := PredictResponse{
		Text:      workerRes.Text,
		ModelUsed: workerRes.ModelUsed,
		LatencyMs: totalLatency,
	}

	if model == "intent" {
		finalRes.Intent = workerRes.Label
	} else {
		finalRes.Sentiment = workerRes.Label
	}

	return finalRes, nil
}
