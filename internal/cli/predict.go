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

var targetModel string

var predictCmd = &cobra.Command{
	Use:   "predict [text]",
	Short: "Run model prediction on input text",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputText := args[0]
		cfgStore := store.NewConfigStore(DB)

		// 1. Fetch JWT token
		token, _ := cfgStore.Get("jwt_token")

		gatewayURL, _ := cfgStore.Get("gateway_url")
		if gatewayURL == "" {
			gatewayURL = "http://localhost:8080"
		}

		// 2. Prepare payload
		payload, _ := json.Marshal(map[string]string{
			"text":  inputText,
			"model": targetModel, // Passed via CLI flag ("auto", "sentiment", "intent")
		})

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
			return fmt.Errorf("failed to reach gateway: %w", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("prediction failed (%d): %s", resp.StatusCode, string(body))
		}

		fmt.Println("🧠 Model Response:")
		fmt.Println(string(body))
		return nil
	},
}

func init() {
	predictCmd.Flags().StringVarP(&targetModel, "model", "m", "auto", "Specify target model (auto, sentiment, intent)")
	RootCmd.AddCommand(predictCmd)
}
