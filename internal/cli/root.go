package cli

import (
	"database/sql"
	"fmt"
	"os"

	"mesh-mind/internal/config"

	"github.com/spf13/cobra"
)

var (
	DB *sql.DB
)

var RootCmd = &cobra.Command{
	Use:   "mesh-mind",
	Short: "mesh-mind CLI - Route predictions and manage MLOps tasks",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		DB, err = config.InitDB()
		if err != nil {
			return fmt.Errorf("SQLite storage initialization failed: %w", err)
		}
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
