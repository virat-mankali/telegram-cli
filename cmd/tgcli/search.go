package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	searchLimit int
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Full-text search across synced messages",
	Long:  `Search your locally synced Telegram messages using SQLite FTS5 for fast offline results.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		fmt.Printf("Searching for %q (limit: %d, store: %s)\n", query, searchLimit, storeDir)
		// Phase 4: SQLite FTS5 query will be wired here.
		fmt.Println("Not yet implemented. Coming in Phase 4.")
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "max number of results to return")
}
