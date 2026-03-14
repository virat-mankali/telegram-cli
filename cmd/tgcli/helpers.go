package main

import (
	"fmt"

	"github.com/gotd/td/tg"
)

// resolveSender extracts a human-readable sender name from message entities.
func resolveSender(e tg.Entities, msg *tg.Message) string {
	if msg.FromID == nil {
		return "unknown"
	}

	switch p := msg.FromID.(type) {
	case *tg.PeerUser:
		if u, ok := e.Users[p.UserID]; ok {
			name := u.FirstName
			if u.LastName != "" {
				name += " " + u.LastName
			}
			if u.Username != "" {
				return fmt.Sprintf("%s (@%s)", name, u.Username)
			}
			return name
		}
		return fmt.Sprintf("user:%d", p.UserID)
	case *tg.PeerChannel:
		if ch, ok := e.Channels[p.ChannelID]; ok {
			return ch.Title
		}
		return fmt.Sprintf("channel:%d", p.ChannelID)
	case *tg.PeerChat:
		if ch, ok := e.Chats[p.ChatID]; ok {
			return ch.Title
		}
		return fmt.Sprintf("chat:%d", p.ChatID)
	default:
		return "unknown"
	}
}

// extractChatID pulls the chat/channel/user ID from the message's PeerID.
func extractChatID(msg *tg.Message) int64 {
	if msg.PeerID == nil {
		return 0
	}
	switch p := msg.PeerID.(type) {
	case *tg.PeerUser:
		return p.UserID
	case *tg.PeerChat:
		return p.ChatID
	case *tg.PeerChannel:
		return p.ChannelID
	default:
		return 0
	}
}

// truncate shortens a string to maxLen, appending "…" if truncated.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}
