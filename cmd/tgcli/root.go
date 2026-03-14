package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// storeDir is the base directory for all tgcli data (~/.tgcli by default).
	storeDir string
)

// defaultStoreDir returns ~/.tgcli as the default data directory.
func defaultStoreDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".tgcli"
	}
	return filepath.Join(home, ".tgcli")
}

var rootCmd = &cobra.Command{
	Use:   "tgcli",
	Short: "Telegram CLI: sync, search, send.",
	Long: `tgcli is a terminal-based Telegram client built on MTProto.
Authenticate as a real user, sync messages locally, search offline, and send from the command line.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Ensure the store directory exists before any subcommand runs.
		return os.MkdirAll(storeDir, 0o700)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&storeDir, "store", defaultStoreDir(), "data directory for tgcli")

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(searchCmd)
}

// Execute runs the root command with graceful signal handling.
func Execute() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	return rootCmd.ExecuteContext(ctx)
}
