# ✈️ tgcli — Telegram CLI: sync, search, send.

Telegram CLI built on top of `gotd/td` (MTProto), focused on:

- Best-effort local sync of message history + continuous capture
- Fast offline full-text search (SQLite FTS5)
- Sending text messages and media files
- Single binary, no CGO, cross-platform
- MCP server for AI agents (Claude, Cursor, Kiro, etc.)

Acts as a real Telegram **user** via MTProto — not a bot. Inspired by [steipete/wacli](https://github.com/steipete/wacli).

---

## Install

### Option A: Homebrew (macOS & Linux) — recommended

```bash
brew install virat-mankali/tap/tgcli
```

### Option B: go install (requires Go)

```bash
go install github.com/virat-mankali/telegram-cli/cmd/tgcli@latest
```

Make sure `$GOPATH/bin` is in your PATH (`export PATH="$PATH:$(go env GOPATH)/bin"`).

### Option C: Build from source

```bash
git clone https://github.com/virat-mankali/telegram-cli.git
cd telegram-cli
go build -o tgcli ./cmd/tgcli/
./tgcli --help
```

---

## Quick Start

Default store directory is `~/.tgcli` (override with `--store DIR`).

```bash
# 1) Get your API credentials at https://my.telegram.org/apps
export TGCLI_APP_ID=your_api_id
export TGCLI_APP_HASH=your_api_hash

# 2) Authenticate (prompts for phone + code + optional 2FA)
tgcli auth

# 3) Sync incoming messages into local DB
tgcli sync

# 4) Keep syncing in real-time (like tail -f)
tgcli sync --follow

# 5) Search messages offline
tgcli search "meeting"
tgcli search "hello" --limit 50

# 6) Send a text message
tgcli send --to @username --text "Hey!"
tgcli send --to +1234567890 --text "Hello"

# 7) Send a file
tgcli send --to @username --media /path/to/file.pdf
tgcli send --to @username --media /path/to/photo.jpg --text "Check this out"
```

---

## Commands

### `tgcli auth`
Authenticate with your Telegram account. Prompts for phone number, verification code, and 2FA password if enabled. Session is saved to `~/.tgcli/session.json` — subsequent commands skip re-auth automatically.

```bash
tgcli auth
tgcli auth --phone +1234567890
```

### `tgcli sync`
Connect to Telegram and sync incoming messages into the local SQLite database.

```bash
tgcli sync             # connect, confirm auth, exit
tgcli sync --follow    # keep running, capture messages in real-time
```

### `tgcli search <query>`
Full-text search across all locally synced messages. Works completely offline.

```bash
tgcli search "standup"
tgcli search "invoice" --limit 100
```

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `20` | Max results to return |

### `tgcli send`
Send a text message or media file to any user, group, or channel.

```bash
tgcli send --to @username --text "Hello!"
tgcli send --to +1234567890 --media /path/to/audio.mp3
tgcli send --to @username --media /path/to/doc.pdf --text "Here's the file"
```

| Flag | Description |
|------|-------------|
| `--to` | Recipient — `@username`, `+phone`, or `t.me/` link |
| `--text` | Text message |
| `--media` | Path to a local file to upload |

#### Supported media types

| Category | Extensions | Sent as |
|----------|-----------|---------|
| Images | `.jpg` `.jpeg` `.png` `.webp` | Photo (inline preview) |
| Audio | `.mp3` `.ogg` `.flac` `.wav` `.m4a` `.aac` | Audio player |
| Video | `.mp4` `.mov` `.avi` `.mkv` `.webm` | Video player |
| Documents | anything else (`.pdf` `.zip` `.docx` …) | Generic document |

### `tgcli serve`
Start tgcli as an MCP server over stdio. AI agents can call Telegram tools natively.

```bash
tgcli serve
```

---

## MCP Server (AI Agents)

tgcli doubles as an MCP server so AI agents (Claude, Cursor, Kiro, etc.) can send messages, search, and list chats on your behalf — directly from the agent, no manual steps needed.

### Prerequisites

You need Go and tgcli installed on your machine. The easiest way is via Homebrew:

```bash
# 1. Install Go (if not already installed)
brew install go

# 2. Install tgcli
brew install virat-mankali/tap/tgcli
```

### One-time Authentication

Get your API credentials from [my.telegram.org/apps](https://my.telegram.org/apps), then authenticate once:

```bash
export TGCLI_APP_ID=your_api_id
export TGCLI_APP_HASH=your_api_hash
tgcli auth
```

This prompts for your phone number, verification code, and 2FA password (if enabled). The session is saved to `~/.tgcli/session.json` — you never need to do this again.

### Add to mcp.json

After authenticating, add tgcli to your MCP config (`~/.kiro/settings/mcp.json` for Kiro, or the equivalent for Claude/Cursor):

```json
{
  "mcpServers": {
    "tgcli": {
      "command": "tgcli",
      "args": ["serve"],
      "env": {
        "TGCLI_APP_ID": "your_api_id",
        "TGCLI_APP_HASH": "your_api_hash"
      },
      "disabled": false,
      "autoApprove": ["search_messages", "list_chats"]
    }
  }
}
```

> If you installed via `go install` instead of Homebrew, use `"command": "go"` and `"args": ["run", "github.com/virat-mankali/telegram-cli/cmd/tgcli@latest", "serve"]` instead.

That's it. Your AI agent can now interact with Telegram on your behalf.

### Available MCP Tools

| Tool | Description | Auto-approved |
|------|-------------|---------------|
| `list_chats` | List recent Telegram chats (dialogs) | ✅ |
| `search_messages` | Full-text search across locally synced messages | ✅ |
| `send_message` | Send a text message or media file to a user/group | requires approval |
| `sync_messages` | Connect to Telegram and verify the session is alive | requires approval |

Read-only tools (`list_chats`, `search_messages`) are safe to auto-approve. Keep `send_message` out of auto-approve so your agent asks before sending anything on your behalf.

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `TGCLI_APP_ID` | Your Telegram API ID from [my.telegram.org/apps](https://my.telegram.org/apps) |
| `TGCLI_APP_HASH` | Your Telegram API hash |

Falls back to `gotd` test credentials if unset — not recommended for production use.

---

## Storage

Defaults to `~/.tgcli` (override with `--store DIR`).

| File | Purpose |
|------|---------|
| `~/.tgcli/session.json` | MTProto session (0600 permissions) |
| `~/.tgcli/data.db` | SQLite database with FTS5 message index |

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
│   ├── serve.go        # tgcli serve (MCP server)
│   └── helpers.go      # Shared utilities
├── internal/
│   ├── config/         # Path helpers (~/.tgcli)
│   ├── mcp/            # MCP server + tool definitions
│   ├── telegram/       # MTProto client wrapper (gotd/td)
│   └── storage/        # SQLite + FTS5
├── .goreleaser.yaml    # Cross-platform build + Homebrew tap
├── .github/workflows/
│   └── release.yml     # Auto-release on git tag push
├── go.mod
└── README.md
```

---

## Releasing

Tag a version and push — GitHub Actions handles the rest:

```bash
git tag v1.0.0
git push origin v1.0.0
```

GoReleaser builds binaries for macOS/Linux/Windows (amd64 + arm64), creates a GitHub Release with checksums, and updates the Homebrew formula automatically.

---

## Prior Art / Credit

Heavily inspired by [steipete/wacli](https://github.com/steipete/wacli) — a WhatsApp CLI built on `whatsmeow`. Same philosophy, different protocol.

---

## License

See `LICENSE`.
