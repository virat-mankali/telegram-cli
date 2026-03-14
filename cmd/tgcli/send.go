package main

import (
	"fmt"

	"github.com/spf13/cobra"
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
		fmt.Printf("Sending to %q: %q\n", sendTo, sendText)
		// Phase 4: MTProto peer resolution + MessagesSendMessage will be wired here.
		fmt.Println("Not yet implemented. Coming in Phase 4.")
		return nil
	},
}

func init() {
	sendCmd.Flags().StringVar(&sendTo, "to", "", "recipient username or phone number")
	sendCmd.Flags().StringVar(&sendText, "text", "", "message text to send")
}
