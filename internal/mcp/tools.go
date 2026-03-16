package mcp

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"

	"github.com/virat-mankali/telegram-cli/internal/config"
	"github.com/virat-mankali/telegram-cli/internal/storage"
	tgclient "github.com/virat-mankali/telegram-cli/internal/telegram"
)

// RegisterTools adds all Telegram MCP tools to the server.
func RegisterTools(s *server.MCPServer, storeDir string) {
	registerSendMessage(s, storeDir)
	registerSearchMessages(s, storeDir)
	registerSyncMessages(s, storeDir)
	registerListChats(s, storeDir)
}

// ── send_message ─────────────────────────────────────────────────────────────

func registerSendMessage(s *server.MCPServer, storeDir string) {
	tool := mcp.NewTool("send_message",
		mcp.WithDescription("Send a text message or media file to a Telegram user or group. Requires prior authentication via 'tgcli auth'."),
		mcp.WithString("to",
			mcp.Required(),
			mcp.Description("Recipient: @username, +phone number, or t.me/ link"),
		),
		mcp.WithString("text",
			mcp.Description("Text message to send"),
		),
		mcp.WithString("media",
			mcp.Description("Absolute path to a local file to upload (image, audio, video, document)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		to, _ := req.RequireString("to")
		text := req.GetString("text", "")
		media := req.GetString("media", "")

		if to == "" {
			return mcp.NewToolResultError("'to' parameter is required"), nil
		}
		if text == "" && media == "" {
			return mcp.NewToolResultError("either 'text' or 'media' (or both) is required"), nil
		}

		sessionPath := config.SessionPath(storeDir)
		client := tgclient.NewClient(tgclient.ClientOptions{
			SessionPath: sessionPath,
		})

		var result string
		err := client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, tgclient.NewAuthFlow("")); err != nil {
				return fmt.Errorf("auth: %w", err)
			}

			api := tg.NewClient(client)
			u := uploader.NewUploader(api)
			sender := message.NewSender(api).WithUploader(u)
			target := sender.Resolve(to)

			if media != "" {
				upload, err := u.FromPath(ctx, media)
				if err != nil {
					return fmt.Errorf("upload %q: %w", media, err)
				}

				filename := filepath.Base(media)
				ext := strings.ToLower(filepath.Ext(media))
				mimeType := mime.TypeByExtension(ext)
				if mimeType == "" {
					mimeType = "application/octet-stream"
				}

				switch {
				case isPhoto(ext):
					if _, err := target.UploadedPhoto(ctx, upload); err != nil {
						return fmt.Errorf("send photo: %w", err)
					}
				case isAudio(ext):
					doc := message.UploadedDocument(upload).MIME(mimeType).Filename(filename).Audio()
					if _, err := target.Media(ctx, doc); err != nil {
						return fmt.Errorf("send audio: %w", err)
					}
				case isVideo(ext):
					doc := message.UploadedDocument(upload).MIME(mimeType).Filename(filename).Video()
					if _, err := target.Media(ctx, doc); err != nil {
						return fmt.Errorf("send video: %w", err)
					}
				default:
					doc := message.UploadedDocument(upload).MIME(mimeType).Filename(filename)
					if _, err := target.Media(ctx, doc); err != nil {
						return fmt.Errorf("send document: %w", err)
					}
				}

				result = fmt.Sprintf("Media sent to %s (%s)", to, filename)

				if text != "" {
					if _, err := target.Text(ctx, text); err != nil {
						return fmt.Errorf("send caption: %w", err)
					}
					result += fmt.Sprintf(" with text: %s", text)
				}
				return nil
			}

			// Text-only.
			if _, err := target.Text(ctx, text); err != nil {
				return fmt.Errorf("send message: %w", err)
			}
			result = fmt.Sprintf("Message sent to %s", to)
			return nil
		})

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}

// ── search_messages ──────────────────────────────────────────────────────────

func registerSearchMessages(s *server.MCPServer, storeDir string) {
	tool := mcp.NewTool("search_messages",
		mcp.WithDescription("Full-text search across locally synced Telegram messages. Works offline. Run 'tgcli sync --follow' first to populate the database."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Full-text search query"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Max number of results to return (default 20)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.RequireString("query")
		if query == "" {
			return mcp.NewToolResultError("'query' parameter is required"), nil
		}

		limit := 20
		if l := req.GetFloat("limit", 0); l > 0 {
			limit = int(l)
		}

		dbPath := config.DBPath(storeDir)
		db, err := storage.InitDB(dbPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("init db: %v", err)), nil
		}
		defer db.Close()

		results, err := storage.SearchMessages(db, query, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("search: %v", err)), nil
		}

		if len(results) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("No results for %q", query)), nil
		}

		var out strings.Builder
		out.WriteString(fmt.Sprintf("Found %d result(s) for %q:\n\n", len(results), query))
		for _, m := range results {
			ts := m.Timestamp.Format(time.DateTime)
			out.WriteString(fmt.Sprintf("  %s  [%s] %s\n", ts, m.Sender, m.Text))
		}
		return mcp.NewToolResultText(out.String()), nil
	})
}

// ── sync_messages ────────────────────────────────────────────────────────────

func registerSyncMessages(s *server.MCPServer, storeDir string) {
	tool := mcp.NewTool("sync_messages",
		mcp.WithDescription("Connect to Telegram and sync incoming messages into the local database. Returns after confirming the connection is alive. For continuous sync, use the CLI: tgcli sync --follow."),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sessionPath := config.SessionPath(storeDir)
		dbPath := config.DBPath(storeDir)

		db, err := storage.InitDB(dbPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("init db: %v", err)), nil
		}
		defer db.Close()

		client := tgclient.NewClient(tgclient.ClientOptions{
			SessionPath: sessionPath,
		})

		var result string
		err = client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, tgclient.NewAuthFlow("")); err != nil {
				return fmt.Errorf("auth: %w", err)
			}

			user, err := client.Self(ctx)
			if err != nil {
				return fmt.Errorf("get self: %w", err)
			}

			result = fmt.Sprintf("Connected as %s %s (ID: %d). Sync is alive. For continuous real-time sync, run: tgcli sync --follow",
				user.FirstName, user.LastName, user.ID)
			return nil
		})

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}

// ── list_chats ───────────────────────────────────────────────────────────────

func registerListChats(s *server.MCPServer, storeDir string) {
	tool := mcp.NewTool("list_chats",
		mcp.WithDescription("List recent Telegram chats (dialogs) the authenticated user is part of."),
		mcp.WithNumber("limit",
			mcp.Description("Max number of chats to return (default 20)"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := 20
		if l := req.GetFloat("limit", 0); l > 0 {
			limit = int(l)
		}

		sessionPath := config.SessionPath(storeDir)
		client := tgclient.NewClient(tgclient.ClientOptions{
			SessionPath: sessionPath,
		})

		var result string
		err := client.Run(ctx, func(ctx context.Context) error {
			if err := client.Auth().IfNecessary(ctx, tgclient.NewAuthFlow("")); err != nil {
				return fmt.Errorf("auth: %w", err)
			}

			api := tg.NewClient(client)
			dialogs, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
				OffsetPeer: &tg.InputPeerEmpty{},
				Limit:      limit,
			})
			if err != nil {
				return fmt.Errorf("get dialogs: %w", err)
			}

			var out strings.Builder
			switch d := dialogs.(type) {
			case *tg.MessagesDialogs:
				result = formatDialogs(d.Dialogs, d.Chats, d.Users)
			case *tg.MessagesDialogsSlice:
				out.WriteString(fmt.Sprintf("Total chats: %d (showing up to %d)\n\n", d.Count, limit))
				out.WriteString(formatDialogs(d.Dialogs, d.Chats, d.Users))
				result = out.String()
			default:
				result = "No chats found."
			}
			return nil
		})

		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	})
}

// formatDialogs builds a human-readable list of chats from dialog data.
func formatDialogs(dialogs []tg.DialogClass, chats []tg.ChatClass, users []tg.UserClass) string {
	// Build lookup maps.
	userMap := make(map[int64]string)
	for _, u := range users {
		if user, ok := u.(*tg.User); ok {
			name := user.FirstName
			if user.LastName != "" {
				name += " " + user.LastName
			}
			if user.Username != "" {
				name += " (@" + user.Username + ")"
			}
			userMap[user.ID] = name
		}
	}

	chatMap := make(map[int64]string)
	for _, c := range chats {
		switch ch := c.(type) {
		case *tg.Chat:
			chatMap[ch.ID] = ch.Title
		case *tg.Channel:
			chatMap[ch.ID] = ch.Title
		}
	}

	var out strings.Builder
	for _, d := range dialogs {
		dialog, ok := d.(*tg.Dialog)
		if !ok {
			continue
		}
		switch p := dialog.Peer.(type) {
		case *tg.PeerUser:
			if name, ok := userMap[p.UserID]; ok {
				out.WriteString(fmt.Sprintf("  User: %s (ID: %d)\n", name, p.UserID))
			}
		case *tg.PeerChat:
			if name, ok := chatMap[p.ChatID]; ok {
				out.WriteString(fmt.Sprintf("  Group: %s (ID: %d)\n", name, p.ChatID))
			}
		case *tg.PeerChannel:
			if name, ok := chatMap[p.ChannelID]; ok {
				out.WriteString(fmt.Sprintf("  Channel: %s (ID: %d)\n", name, p.ChannelID))
			}
		}
	}

	if out.Len() == 0 {
		return "No chats found."
	}
	return out.String()
}

// File type helpers (same logic as CLI send command).

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
