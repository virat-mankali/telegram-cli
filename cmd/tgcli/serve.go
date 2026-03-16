package main

import (
	"github.com/spf13/cobra"

	mcpserver "github.com/virat-mankali/telegram-cli/internal/mcp"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server for AI agents (stdio)",
	Long: `Start tgcli as an MCP server over stdio.

AI agents (Claude, Cursor, Kiro, etc.) can call Telegram tools natively.
Requires prior authentication via 'tgcli auth'.

Add this to your mcp.json:

  "tgcli": {
    "command": "/path/to/tgcli",
    "args": ["serve"],
    "env": {
      "TGCLI_APP_ID": "your_api_id",
      "TGCLI_APP_HASH": "your_api_hash"
    },
    "disabled": false,
    "autoApprove": ["search_messages", "list_chats"]
  }

Available tools: send_message, search_messages, sync_messages, list_chats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcpserver.Serve(storeDir)
	},
}

func init() {
	// serveCmd is registered in root.go
}
