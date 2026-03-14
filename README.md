# tgcli — Telegram CLI

A blazing-fast, terminal-based Telegram client built on MTProto. Authenticate as a real user, sync messages locally, search offline, and send messages or media directly from the command line.

Inspired by [steipete/wacli](https://github.com/steipete/wacli) (WhatsApp CLI).

## Getting Started

### 1. Get Your Telegram API Credentials

1. Go to https://my.telegram.org/apps
2. Log in with your phone number
3. Create a new application (any name works, e.g., "tgcli")
4. Copy your `api_id` and `api_hash`

### 2. Set Environment Variables

```bash
export TGCLI_APP_ID=your_api_id
export TGCLI_APP_HASH=your_api_hash
```

### 3. Authenticate

```bash
go run ./cmd/tgcli/ auth
```

You'll be prompted for:
1. **Phone number** (e.g., `+1234567890`)
2. **Verification code** (sent via SMS or Telegram app)
3. **2FA password** (if enabled — input is hidden)

Session is saved to `~/.tgcli/session.json`. Subsequent commands skip re-authentication automatically.

### 4. Build the Binary

```bash
go build -o tgcli ./cmd/tgcli/
./tgcli --help
```

---

## Commands

### `tgcli auth`
Authenticate with your Telegram account.

```bash
tgcli auth
tgcli auth --phone +1234567890
```

---

### `tgcli sync`
Connect to Telegram and sync incoming messages into the local SQLite database.

```bash
tgcli sync                  # connect and exit
tgcli sync --follow         # keep running, capture messages in real-time (like tail -f)
```

---

### `tgcli search <query>`
Full-text search across all locally synced messages using SQLite FTS5. Works completely offline.

```bash
tgcli search "meeting"
tgcli search "hello" --limit 50
```

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `20` | Max number of results to return |

---

### `tgcli send`
Send a text message or media file to any Telegram user, group, or channel.

```bash
# Text message
tgcli send --to @username --text "Hello!"
tgcli send --to +1234567890 --text "Hey there"

# Send media
tgcli send --to @username --media /path/to/file.png
tgcli send --to +1234567890 --media /path/to/audio.mp3

# Send media with a follow-up text
tgcli send --to @username --media /path/to/doc.pdf --text "Here's the file"
```

| Flag | Description |
|------|-------------|
| `--to` | Recipient — `@username`, phone number (`+1234567890`), or `t.me/` link |
| `--text` | Text message to send |
| `--media` | Path to a local file to upload and send |

#### Supported Media Types

| Category | Extensions | Sent as |
|----------|-----------|---------|
| **Images** | `.jpg`, `.jpeg`, `.png`, `.webp` | Photo (inline preview) |
| **Audio** | `.mp3`, `.ogg`, `.flac`, `.wav`, `.m4a`, `.aac` | Audio player |
| **Video** | `.mp4`, `.mov`, `.avi`, `.mkv`, `.webm` | Video player |
| **Documents** | `.pdf`, `.zip`, `.docx`, `.xlsx`, `.txt`, and any other file type | Generic document |

Any file type not in the image/audio/video lists is sent as a generic document (e.g. PDFs, archives, spreadsheets, code files, etc.).

---

## Project Structure

```
tgcli/
├── cmd/tgcli/
│   ├── main.go         # Entry point
│   ├── root.go         # Cobra root command, --store flag, signal handling
│   ├── auth.go         # tgcli auth
│   ├── sync.go         # tgcli sync
│   ├── send.go         # tgcli send
│   ├── search.go       # tgcli search
│   └── helpers.go      # Shared utilities (sender resolution, truncation)
├── internal/
│   ├── config/         # Path helpers (~/.tgcli)
│   ├── telegram/       # MTProto client wrapper
│   │   ├── client.go
│   │   └── auth.go
│   └── storage/        # SQLite + FTS5
│       ├── db.go
│       ├── session.go
│       └── messages.go
├── go.mod
└── README.md
```

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/gotd/td` | Telegram MTProto API (core) |
| `github.com/gotd/contrib` | Flood wait & rate limit middleware |
| `github.com/spf13/cobra` | CLI framework |
| `modernc.org/sqlite` | Pure-Go SQLite driver (no CGO) |
| `github.com/fatih/color` | Colorized terminal output |
| `golang.org/x/term` | Secure password input (no echo) |

---

## Implementation Phases

- **Phase 1** ✅ — Project setup & CLI routing (Cobra)
- **Phase 2** ✅ — Telegram client & authentication (gotd/td)
- **Phase 3** ✅ — Local storage & SQLite FTS5
- **Phase 4** ✅ — Core features: sync, search, send + media upload
- **Phase 5** — Build & distribution (GoReleaser + Homebrew tap)

---

## Notes

- **Session file**: `~/.tgcli/session.json` — 0600 permissions (owner only)
- **Database**: `~/.tgcli/data.db` — SQLite with FTS5 for offline search
- **Data directory**: Override with `--store /path/to/dir`
- **Flood wait**: Automatically retried up to 10 times via middleware
- **Test credentials**: Falls back to gotd test credentials if env vars are not set — not recommended for production use

## License

See LICENSE file.
