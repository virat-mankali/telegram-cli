package storage

// Session persistence is handled by gotd's built-in session.FileStorage.
// It stores the MTProto auth key, DC info, and salt as JSON at the configured path.
//
// Usage:
//
//	import "github.com/gotd/td/session"
//	storage := &session.FileStorage{Path: "~/.tgcli/session.json"}
//
// The file is written with 0600 permissions (owner read/write only).
// No custom implementation is needed for tgcli.
