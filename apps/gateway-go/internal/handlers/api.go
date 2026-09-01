package handlers

import (
	"bytes"
	"encoding/json"
	"gateway-go/internal/auth"
	"gateway-go/internal/db"
	"io"
	"net/http"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	InferenceURL string
	mu           sync.Mutex
)

func init() {
	InferenceURL = os.Getenv("INFERENCE_URL")
	if InferenceURL == "" {
		InferenceURL = "http://localhost:8000"
	}
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	_, err := db.DB.Exec("INSERT INTO users (username, password_hash, tokens) VALUES (?, ?, 10)", req.Username, string(hash))
	if err != nil {
		http.Error(w, `{"error":"Username taken"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created with 10 free tokens!"})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var id, tokens int
	var hash string
	err := db.DB.QueryRow("SELECT id, password_hash, tokens FROM users WHERE username = ?", req.Username).
		Scan(&id, &hash, &tokens)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateJWT(req.Username)
	json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "tokens": tokens, "username": req.Username})
}

func HandleBuyTokens(w http.ResponseWriter, r *http.Request, username string) {
	mu.Lock()
	defer mu.Unlock()

	_, err := db.DB.Exec("UPDATE users SET tokens = tokens + 20 WHERE username = ?", username)
	if err != nil {
		http.Error(w, `{"error":"Failed to add tokens"}`, http.StatusInternalServerError)
		return
	}

	var newBalance int
	db.DB.QueryRow("SELECT tokens FROM users WHERE username = ?", username).Scan(&newBalance)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Added 20 tokens!", "tokens": newBalance})
}

func HandleChat(w http.ResponseWriter, r *http.Request, username string) {
	mu.Lock()
	var tokens int
	err := db.DB.QueryRow("SELECT tokens FROM users WHERE username = ?", username).Scan(&tokens)
	if err != nil || tokens <= 0 {
		mu.Unlock()
		http.Error(w, `{"error":"Insufficient tokens! Buy more."}`, http.StatusPaymentRequired)
		return
	}

	db.DB.Exec("UPDATE users SET tokens = tokens - 1 WHERE username = ?", username)
	mu.Unlock()

	var chatReq struct {
		Prompt string `json:"prompt"`
	}
	json.NewDecoder(r.Body).Decode(&chatReq)

	payload, _ := json.Marshal(map[string]interface{}{
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful, concise AI assistant."},
			{"role": "user", "content": chatReq.Prompt},
		},
		"max_new_tokens": 150,
		"temperature":    0.7,
	})
	resp, err := http.Post(InferenceURL+"/v1/chat/completions", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, `{"error":"Inference service offline"}`, http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var pyResp struct {
		Response string `json:"response"`
	}
	json.Unmarshal(body, &pyResp)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": pyResp.Response,
		"tokens":   tokens - 1,
	})
}
