package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/virat-mankali/telegram-cli/internal/config"
	tgclient "github.com/virat-mankali/telegram-cli/internal/telegram"
)

var (
	authPhone string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Telegram (interactive login)",
	Long: `Authenticate with your Telegram account via phone number and verification code.
On success, the session is saved locally so subsequent commands skip re-authentication.

Set TGCLI_APP_ID and TGCLI_APP_HASH environment variables with your own credentials
from https://my.telegram.org/apps for production use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		sessionPath := config.SessionPath(storeDir)

		client := tgclient.NewClient(tgclient.ClientOptions{
			SessionPath: sessionPath,
		})

		flow := tgclient.NewAuthFlow(authPhone)

		fmt.Println("Connecting to Telegram...")
		return client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				return fmt.Errorf("auth: %w", err)
			}

			user, err := client.Self(ctx)
			if err != nil {
				return fmt.Errorf("get self: %w", err)
			}

			fmt.Printf("✓ Authenticated as %s %s (ID: %d)\n",
				user.FirstName, user.LastName, user.ID)
			fmt.Printf("  Session saved to %s\n", sessionPath)
			return nil
		})
	},
}

func init() {
	authCmd.Flags().StringVar(&authPhone, "phone", "",
		"phone number in international format (e.g. +1234567890)")
}
