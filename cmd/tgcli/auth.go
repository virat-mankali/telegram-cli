package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Telegram (interactive login)",
	Long: `Authenticate with your Telegram account via phone number and verification code.
On success, the session is saved locally so subsequent commands skip re-authentication.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Authenticating... (store: %s)\n", storeDir)
		// Phase 2: MTProto auth flow via gotd/td will be wired here.
		fmt.Println("Not yet implemented. Coming in Phase 2.")
		return nil
	},
}
