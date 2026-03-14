package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// Message represents a single Telegram message stored locally.
type Message struct {
	ID        int
	ChatID    int64
	Sender    string
	Text      string
	Timestamp time.Time
}

// InsertMessage stores a message in the FTS5 table.
func InsertMessage(db *sql.DB, m Message) error {
	_, err := db.Exec(
		`INSERT INTO messages (id, chat_id, sender, text, timestamp) VALUES (?, ?, ?, ?, ?)`,
		m.ID, m.ChatID, m.Sender, m.Text, m.Timestamp.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

// SearchMessages queries the FTS5 index and returns matching messages.
func SearchMessages(db *sql.DB, query string, limit int) ([]Message, error) {
	rows, err := db.Query(
		`SELECT id, chat_id, sender, text, timestamp FROM messages WHERE text MATCH ? ORDER BY rank LIMIT ?`,
		query, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("search messages: %w", err)
	}
	defer rows.Close()

	var results []Message
	for rows.Next() {
		var m Message
		var ts int64
		if err := rows.Scan(&m.ID, &m.ChatID, &m.Sender, &m.Text, &ts); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		m.Timestamp = time.Unix(ts, 0)
		results = append(results, m)
	}
	return results, rows.Err()
}
