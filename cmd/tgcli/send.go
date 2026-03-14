package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"

	"github.com/virat-mankali/telegram-cli/internal/config"
	tgclient "github.com/virat-mankali/telegram-cli/internal/telegram"
)

var (
	sendTo   string
	sendText string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message via Telegram",
	Long:  `Send a text message to a Telegram user or group by username or phone number.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sendTo == "" || sendText == "" {
			return fmt.Errorf("both --to and --text are required")
		}

		ctx := cmd.Context()
		sessionPath := config.SessionPath(storeDir)

		client := tgclient.NewClient(tgclient.ClientOptions{
			SessionPath: sessionPath,
		})

		flow := tgclient.NewAuthFlow("")

		return client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				return fmt.Errorf("auth: %w", err)
			}

			api := tg.NewClient(client)
			sender := message.NewSender(api)

			// Resolve handles @username, phone numbers, and t.me links.
			target := sender.Resolve(sendTo)

			if _, err := target.Text(ctx, sendText); err != nil {
				return fmt.Errorf("send message: %w", err)
			}

			fmt.Printf("✓ Message sent to %s\n", sendTo)
			return nil
		})
	},
}

func init() {
	sendCmd.Flags().StringVar(&sendTo, "to", "", "recipient username or phone number (e.g. @user or +1234567890)")
	sendCmd.Flags().StringVar(&sendText, "text", "", "message text to send")
}
