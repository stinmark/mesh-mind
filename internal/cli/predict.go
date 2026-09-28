package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mesh-mind/internal/store"

	"github.com/spf13/cobra"
)

type APIResponse struct {
	Intent     string  `json:"intent"`
	Sentiment  string  `json:"sentiment"`
	ModelUsed  string  `json:"model_used"`
	LatencyMs  float64 `json:"latency_ms"`
}

var predictCmd = &cobra.Command{
	Use:   "predict [text]",
	Short: "Execute model pipeline on input text and save result locally",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputText := args[0]
		cfgStore := store.NewConfigStore(DB)

		// 1. Resolve Gateway URL & Token from SQLite
		gatewayURL, _ := cfgStore.Get("gateway_url")
		if gatewayURL == "" {
			gatewayURL = "http://localhost:8080" // Default fallback
		}

		token, _ := cfgStore.Get("jwt_token")

		// 2. Call Gateway API
		payload, _ := json.Marshal(map[string]string{"text": inputText})
		req, err := http.NewRequest("POST", gatewayURL+"/api/v1/predict", bytes.NewBuffer(payload))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to reach gateway (%s): %w", gatewayURL, err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("gateway error (%d): %s", resp.StatusCode, string(body))
		}

		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		// 3. Print Results
		fmt.Printf("\n🧠 Model Output:\n")
		fmt.Printf("   Intent:     %s\n", apiResp.Intent)
		fmt.Printf("   Sentiment:  %s\n", apiResp.Sentiment)
		fmt.Printf("   Model Used: %s\n", apiResp.ModelUsed)
		fmt.Printf("   Latency:    %.2f ms\n\n", apiResp.LatencyMs)

		// 4. Save Record to SQLite History
		histStore := store.NewHistoryStore(DB)
		if err := histStore.AddRecord(inputText, apiResp.Intent, apiResp.Sentiment, apiResp.ModelUsed, apiResp.LatencyMs); err != nil {
			fmt.Printf("⚠️ Warning: Failed to save record to local history: %v\n", err)
		} else {
			fmt.Println("Saved to local SQLite history.")
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(predictCmd)
}
