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

var (
	loginUsername string
	loginPassword string
)

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Gateway and save the JWT token locally",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgStore := store.NewConfigStore(DB)

		// 1. Resolve Gateway URL from local SQLite database
		gatewayURL, _ := cfgStore.Get("gateway_url")
		if gatewayURL == "" {
			gatewayURL = "http://localhost:8080" // Default fallback
		}

		// 2. Prepare payload
		payload, err := json.Marshal(map[string]string{
			"username": loginUsername,
			"password": loginPassword,
		})
		if err != nil {
			return fmt.Errorf("failed to encode login payload: %w", err)
		}

		// 3. Send request to Gateway /api/v1/login
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Post(gatewayURL+"/api/v1/login", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("failed to connect to gateway (%s): %w", gatewayURL, err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("login failed (%d): %s", resp.StatusCode, string(body))
		}

		// 4. Parse JWT token from response
		var res LoginResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		if res.Token == "" {
			return fmt.Errorf("received empty token from gateway")
		}

		// 5. Save JWT token in local client SQLite database
		if err := cfgStore.Set("jwt_token", res.Token); err != nil {
			return fmt.Errorf("failed to store token in local database: %w", err)
		}

		fmt.Println("✔ Login successful!")
		fmt.Println("✔ JWT token saved to local SQLite database (~/.mesh-mind/config.db)")
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Username for authentication")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Password for authentication")

	_ = loginCmd.MarkFlagRequired("username")
	_ = loginCmd.MarkFlagRequired("password")

	RootCmd.AddCommand(loginCmd)
}
