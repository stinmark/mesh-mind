package cli

import (
	"fmt"
	"mesh-mind/internal/store"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration stored in local SQLite database",
}

var setUrlCmd = &cobra.Command{
	Use:   "set-url [url]",
	Short: "Set default Gateway API URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgStore := store.NewConfigStore(DB)
		if err := cfgStore.Set("gateway_url", args[0]); err != nil {
			return err
		}
		fmt.Printf("✔ Default Gateway URL set to: %s\n", args[0])
		return nil
	},
}

var setTokenCmd = &cobra.Command{
	Use:   "set-token [jwt-token]",
	Short: "Save JWT authentication token locally",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgStore := store.NewConfigStore(DB)
		if err := cfgStore.Set("jwt_token", args[0]); err != nil {
			return err
		}
		fmt.Println("✔ JWT Token securely stored in local database.")
		return nil
	},
}

func init() {
	configCmd.AddCommand(setUrlCmd)
	configCmd.AddCommand(setTokenCmd)
	RootCmd.AddCommand(configCmd)
}
