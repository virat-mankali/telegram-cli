package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/gotd/td/tg"

	"github.com/virat-mankali/telegram-cli/internal/config"
	"github.com/virat-mankali/telegram-cli/internal/storage"
	tgclient "github.com/virat-mankali/telegram-cli/internal/telegram"
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
		ctx := cmd.Context()
		sessionPath := config.SessionPath(storeDir)
		dbPath := config.DBPath(storeDir)

		db, err := storage.InitDB(dbPath)
		if err != nil {
			return fmt.Errorf("init db: %w", err)
		}
		defer db.Close()

		// Set up the update dispatcher to capture new messages.
		dispatcher := tg.NewUpdateDispatcher()
		msgCount := 0

		dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
			msg, ok := u.Message.(*tg.Message)
			if !ok || msg.Out {
				return nil
			}

			// Resolve sender name from entities.
			sender := resolveSender(e, msg)

			// Extract chat ID from the peer.
			chatID := extractChatID(msg)

			m := storage.Message{
				ID:        msg.ID,
				ChatID:    chatID,
				Sender:    sender,
				Text:      msg.Message,
				Timestamp: time.Unix(int64(msg.Date), 0),
			}

			if err := storage.InsertMessage(db, m); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "⚠ insert error: %v\n", err)
				return nil // Don't kill the sync loop on insert errors.
			}

			msgCount++
			fmt.Printf("  [%d] %s: %s\n", msg.ID, sender, truncate(msg.Message, 80))
			return nil
		})

		// Also capture channel messages.
		dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewChannelMessage) error {
			msg, ok := u.Message.(*tg.Message)
			if !ok || msg.Out {
				return nil
			}

			sender := resolveSender(e, msg)
			chatID := extractChatID(msg)

			m := storage.Message{
				ID:        msg.ID,
				ChatID:    chatID,
				Sender:    sender,
				Text:      msg.Message,
				Timestamp: time.Unix(int64(msg.Date), 0),
			}

			if err := storage.InsertMessage(db, m); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "⚠ insert error: %v\n", err)
				return nil
			}

			msgCount++
			fmt.Printf("  [%d] %s: %s\n", msg.ID, sender, truncate(msg.Message, 80))
			return nil
		})

		client := tgclient.NewClient(tgclient.ClientOptions{
			SessionPath:   sessionPath,
			UpdateHandler: dispatcher,
		})

		flow := tgclient.NewAuthFlow("")

		fmt.Println("Connecting to Telegram...")
		return client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, flow); err != nil {
				return fmt.Errorf("auth: %w", err)
			}

			user, err := client.Self(ctx)
			if err != nil {
				return fmt.Errorf("get self: %w", err)
			}
			fmt.Printf("✓ Authenticated as %s %s\n", user.FirstName, user.LastName)

			if syncFollow {
				fmt.Println("Listening for new messages (Ctrl+C to stop)...")
				<-ctx.Done()
				fmt.Printf("\n✓ Synced %d messages\n", msgCount)
				return ctx.Err()
			}

			// Without --follow, just confirm connection and exit.
			fmt.Println("✓ Connected. Use --follow to listen for incoming messages.")
			return nil
		})
	},
}

func init() {
	syncCmd.Flags().BoolVar(&syncFollow, "follow", false, "keep syncing continuously (like tail -f)")
}
