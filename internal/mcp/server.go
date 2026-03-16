// Package mcp exposes tgcli functionality as an MCP server over stdio.
package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// Serve starts the MCP server over stdio. It registers all Telegram tools
// and blocks until the connection is closed.
func Serve(storeDir string) error {
	s := server.NewMCPServer(
		"tgcli",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	RegisterTools(s, storeDir)

	return server.ServeStdio(s)
}
