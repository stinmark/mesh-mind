package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"mesh-mind/internal/store"

	"github.com/spf13/cobra"
)

var (
	regUsername string
	regPassword string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new user account on the Mesh-Mind Gateway",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgStore := store.NewConfigStore(DB)
		gatewayURL, _ := cfgStore.Get("gateway_url")
		if gatewayURL == "" {
			gatewayURL = "http://localhost:8080"
		}

		payload, _ := json.Marshal(map[string]string{
			"username": regUsername,
			"password": regPassword,
		})

		resp, err := http.Post(gatewayURL+"/api/v1/register", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return fmt.Errorf("registration failed: %w", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("registration error (%d): %s", resp.StatusCode, string(body))
		}

		fmt.Println("✔ Account created successfully! You can now log in using:")
		fmt.Printf("   mesh-mind login -u %s -p <your-password>\n", regUsername)
		return nil
	},
}

func init() {
	registerCmd.Flags().StringVarP(&regUsername, "username", "u", "", "Username")
	registerCmd.Flags().StringVarP(&regPassword, "password", "p", "", "Password")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")

	RootCmd.AddCommand(registerCmd)
}
