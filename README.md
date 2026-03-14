# tgcli — Telegram CLI

A blazing-fast, terminal-based Telegram client built on MTProto. Authenticate as a real user, sync messages locally, search offline, and send from the command line.

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

Or pass them inline:
```bash
TGCLI_APP_ID=12345678 TGCLI_APP_HASH=abcdef... tgcli auth
```

### 3. Authenticate

```bash
go run ./cmd/tgcli/ auth
```

You'll be prompted for:
1. **Phone number** (e.g., `+1234567890`)
2. **Verification code** (Telegram sends it via SMS or app notification)
3. **2FA password** (if you have 2FA enabled; input is hidden)

On success, your session is saved to `~/.tgcli/session.json` and you won't need to re-authenticate.

### 4. Build the Binary

```bash
go build -o ./dist/tgcli ./cmd/tgcli/
./dist/tgcli --help
```

## Commands

- `tgcli auth` — Authenticate with Telegram
- `tgcli sync` — Sync messages from Telegram (Phase 4)
- `tgcli search <query>` — Full-text search synced messages (Phase 4)
- `tgcli send --to <user> --text <msg>` — Send a message (Phase 4)

## Development

### Project Structure

```
tgcli/
├── cmd/tgcli/              # CLI entry points
│   ├── main.go
│   ├── root.go
│   ├── auth.go
│   ├── sync.go
│   ├── send.go
│   └── search.go
├── internal/
│   ├── config/             # Config & paths
│   ├── telegram/           # MTProto client wrapper
│   │   ├── client.go
│   │   └── auth.go
│   └── storage/            # SQLite & session storage
│       ├── db.go
│       ├── session.go
│       └── messages.go
├── go.mod
└── README.md
```

### Dependencies

- `github.com/gotd/td` — Telegram MTProto API
- `github.com/gotd/contrib` — Middleware (flood wait, rate limit)
- `github.com/spf13/cobra` — CLI framework
- `golang.org/x/term` — Secure password input

### Running Tests

```bash
go build ./cmd/tgcli/
go run ./cmd/tgcli/ --help
```

## Implementation Phases

- **Phase 1** ✓ — Project setup & CLI routing (Cobra)
- **Phase 2** ✓ — Telegram client & authentication (gotd/td)
- **Phase 3** — Local storage & SQLite FTS (coming soon)
- **Phase 4** — Core features: sync, search, send (coming soon)
- **Phase 5** — Build & distribution (GoReleaser, Homebrew)

## Notes

- **Session file**: Stored at `~/.tgcli/session.json` with 0600 permissions (owner read/write only)
- **Data directory**: Override with `--store /path/to/dir`
- **Test credentials**: The code falls back to gotd's test credentials if `TGCLI_APP_ID`/`TGCLI_APP_HASH` are not set, but these are rate-limited and not recommended for production
- **Flood wait**: Automatically handled by middleware (retries up to 10 times)

## License

See LICENSE file.
