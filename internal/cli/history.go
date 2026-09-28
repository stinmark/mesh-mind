package cli

import (
	"fmt"
	"mesh-mind/internal/store"

	"github.com/spf13/cobra"
)

var limit int

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View recent local predictions stored in SQLite",
	RunE: func(cmd *cobra.Command, args []string) error {
		histStore := store.NewHistoryStore(DB)
		records, err := histStore.GetRecent(limit)
		if err != nil {
			return fmt.Errorf("failed to fetch history: %w", err)
		}

		if len(records) == 0 {
			fmt.Println("No local prediction history found.")
			return nil
		}

		fmt.Printf("\n📜 Last %d Local Predictions:\n", len(records))
		fmt.Println("---------------------------------------------------------------------")
		for _, r := range records {
			fmt.Printf("[%s] Text: \"%s\"\n", r.CreatedAt.Format("2006-01-02 15:04:05"), r.InputText)
			fmt.Printf("   └─ Intent: %s | Sentiment: %s | Model: %s (%.1f ms)\n\n",
				r.Intent, r.Sentiment, r.ModelUsed, r.LatencyMs)
		}
		return nil
	},
}

func init() {
	historyCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Number of history records to show")
	RootCmd.AddCommand(historyCmd)
}
