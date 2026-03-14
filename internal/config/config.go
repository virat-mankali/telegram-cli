// Package config handles tgcli environment and default directory setup.
package config

import (
	"os"
	"path/filepath"
)

const (
	// DefaultDirName is the hidden directory under $HOME.
	DefaultDirName = ".tgcli"
	// DBFileName is the SQLite database file.
	DBFileName = "data.db"
	// SessionFileName stores the MTProto session.
	SessionFileName = "session.json"
)

// StoreDir returns the resolved store directory, ensuring it exists.
func StoreDir(override string) (string, error) {
	dir := override
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, DefaultDirName)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

// DBPath returns the full path to the SQLite database.
func DBPath(storeDir string) string {
	return filepath.Join(storeDir, DBFileName)
}

// SessionPath returns the full path to the session file.
func SessionPath(storeDir string) string {
	return filepath.Join(storeDir, SessionFileName)
}
