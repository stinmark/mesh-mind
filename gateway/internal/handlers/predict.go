package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type PredictRequest struct {
	Text  string `json:"text"`
	Model string `json:"model,omitempty"` // e.g., "auto", "sentiment", "intent", or "classifier"
}

type PredictResponse struct {
	Text      string `json:"text"`
	ModelUsed string `json:"model_used"`
	Intent    string `json:"intent,omitempty"`
	Sentiment string `json:"sentiment,omitempty"`
	LatencyMs float64 `json:"latency_ms"`
}

func PredictHandler(w http.ResponseWriter, r *http.Request) {
	var req PredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Text) == "" {
		http.Error(w, "Text field is required", http.StatusBadRequest)
		return
	}

	targetModel := strings.ToLower(req.Model)

	// Routing Logic: Auto-resolve if empty or set to "auto"
	if targetModel == "" || targetModel == "auto" {
		targetModel = routeAuto(req.Text)
	}

	// Route to specific model worker...
	res := dispatchToWorker(targetModel, req.Text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Simple rule-based or lightweight ML auto-router
func routeAuto(text string) string {
	lowerText := strings.ToLower(text)

	// If user is asking a question or technical query -> intent classifier
	if strings.Contains(lowerText, "how") || strings.Contains(lowerText, "why") || strings.Contains(lowerText, "?") {
		return "intent"
	}

	// Default fallback model
	return "sentiment"
}

func dispatchToWorker(model string, text string) PredictResponse {
	// Send to Python ONNX backend for the selected model
	return PredictResponse{
		Text:      text,
		ModelUsed: "onnx-" + model + "-v1",
		Intent:    "general_query",
		Sentiment: "positive",
		LatencyMs: 12.4,
	}
}
