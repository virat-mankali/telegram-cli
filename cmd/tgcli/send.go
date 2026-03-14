package main

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"

	"github.com/virat-mankali/telegram-cli/internal/config"
	tgclient "github.com/virat-mankali/telegram-cli/internal/telegram"
)

var (
	sendTo    string
	sendText  string
	sendMedia string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message or media via Telegram",
	Long: `Send a text message, image, audio, video, or any file to a Telegram user or group.
Supports @username, phone numbers, and t.me links as recipients.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sendTo == "" {
			return fmt.Errorf("--to is required")
		}
		if sendText == "" && sendMedia == "" {
			return fmt.Errorf("either --text or --media (or both) is required")
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
			u := uploader.NewUploader(api)
			sender := message.NewSender(api).WithUploader(u)
			target := sender.Resolve(sendTo)

			// Media send path.
			if sendMedia != "" {
				fmt.Printf("Uploading %s...\n", filepath.Base(sendMedia))

				upload, err := u.FromPath(ctx, sendMedia)
				if err != nil {
					return fmt.Errorf("upload %q: %w", sendMedia, err)
				}

				filename := filepath.Base(sendMedia)
				ext := strings.ToLower(filepath.Ext(sendMedia))
				mimeType := mime.TypeByExtension(ext)
				if mimeType == "" {
					mimeType = "application/octet-stream"
				}

				// Route by file type: photo, audio, or generic document.
				switch {
				case isPhoto(ext):
					if _, err := target.UploadedPhoto(ctx, upload); err != nil {
						return fmt.Errorf("send photo: %w", err)
					}

				case isAudio(ext):
					doc := message.UploadedDocument(upload).
						MIME(mimeType).
						Filename(filename).
						Audio()
					if _, err := target.Media(ctx, doc); err != nil {
						return fmt.Errorf("send audio: %w", err)
					}

				case isVideo(ext):
					doc := message.UploadedDocument(upload).
						MIME(mimeType).
						Filename(filename).
						Video()
					if _, err := target.Media(ctx, doc); err != nil {
						return fmt.Errorf("send video: %w", err)
					}

				default:
					// Generic document (PDF, ZIP, etc.)
					doc := message.UploadedDocument(upload).
						MIME(mimeType).
						Filename(filename)
					if _, err := target.Media(ctx, doc); err != nil {
						return fmt.Errorf("send document: %w", err)
					}
				}

				fmt.Printf("✓ Media sent to %s (%s)\n", sendTo, filename)

				// If text is also provided, send it as a follow-up.
				if sendText != "" {
					if _, err := target.Text(ctx, sendText); err != nil {
						return fmt.Errorf("send caption: %w", err)
					}
				}
				return nil
			}

			// Text-only message.
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
	sendCmd.Flags().StringVar(&sendMedia, "media", "", "path to file to send (image, audio, video, PDF, etc.)")
}

func isPhoto(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	}
	return false
}

func isAudio(ext string) bool {
	switch ext {
	case ".mp3", ".ogg", ".flac", ".wav", ".m4a", ".aac":
		return true
	}
	return false
}

func isVideo(ext string) bool {
	switch ext {
	case ".mp4", ".mov", ".avi", ".mkv", ".webm":
		return true
	}
	return false
}
