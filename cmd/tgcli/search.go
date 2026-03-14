package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/virat-mankali/telegram-cli/internal/config"
	"github.com/virat-mankali/telegram-cli/internal/storage"
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
		dbPath := config.DBPath(storeDir)

		db, err := storage.InitDB(dbPath)
		if err != nil {
			return fmt.Errorf("init db: %w", err)
		}
		defer db.Close()

		results, err := storage.SearchMessages(db, query, searchLimit)
		if err != nil {
			return fmt.Errorf("search: %w", err)
		}

		if len(results) == 0 {
			fmt.Printf("No results for %q\n", query)
			return nil
		}

		bold := color.New(color.Bold)
		dim := color.New(color.FgHiBlack)
		cyan := color.New(color.FgCyan)

		fmt.Printf("Found %d result(s) for %q:\n\n", len(results), query)
		for _, m := range results {
			ts := m.Timestamp.Format(time.DateTime)
			dim.Printf("  %s  ", ts)
			cyan.Printf("[%s] ", m.Sender)
			bold.Println(m.Text)
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "max number of results to return")
}
