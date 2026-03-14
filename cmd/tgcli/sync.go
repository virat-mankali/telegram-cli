package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	syncFollow bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync messages from Telegram to local storage",
	Long: `Connect to Telegram and sync incoming messages into the local SQLite database.
Use --follow to keep the connection open and continuously capture new messages.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Syncing messages... (store: %s, follow: %v)\n", storeDir, syncFollow)
		// Phase 4: MTProto update dispatcher + SQLite insert will be wired here.
		fmt.Println("Not yet implemented. Coming in Phase 4.")
		return nil
	},
}

func init() {
	syncCmd.Flags().BoolVar(&syncFollow, "follow", false, "keep syncing continuously (like tail -f)")
}
