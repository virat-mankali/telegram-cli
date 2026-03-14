✈️ tgcli - Telegram Command Line Interface Master Architecture Plan

1. Project Overview

tgcli is a blazing-fast, terminal-based client for Telegram, heavily inspired by steipete/wacli. It allows users to authenticate, sync messages locally, perform lightning-fast offline full-text searches, and send messages/media directly from the Unix terminal.

Primary Objectives:

Provide a native, Unix-like CLI experience.

Act as a real Telegram user via MTProto (not a bot).

Store session data and synced messages locally in SQLite for rapid offline access.

Be easily distributable via a single binary and Homebrew.

2. Tech Stack & Dependencies

Language: Go (Golang) - chosen for its speed and single-binary compilation.

CLI Framework: github.com/spf13/cobra - the industry standard for Go CLI routing.

Telegram API Engine: github.com/gotd/td - A pure Go implementation of the Telegram MTProto API.

Local Storage & FTS: modernc.org/sqlite (Pure Go SQLite to avoid CGO headaches) or mattn/go-sqlite3 with sqlite_fts5 build tags.

Distribution: GoReleaser (for cross-compilation and Homebrew tap automation).

3. Project Directory Structure

The AI agent must strictly adhere to this standard Go project layout:

tgcli/
├── cmd/
│   └── tgcli/                  # Application entry points
│       ├── main.go             # Boots up the Cobra root command
│       ├── root.go             # Defines the base `tgcli` command
│       ├── auth.go             # Logic for `tgcli auth`
│       ├── sync.go             # Logic for `tgcli sync`
│       ├── send.go             # Logic for `tgcli send`
│       └── search.go           # Logic for `tgcli search`
├── internal/
│   ├── telegram/               # Wrappers around `gotd/td`
│   │   ├── client.go           # MTProto connection management
│   │   └── auth.go             # Terminal-based authentication flow
│   ├── storage/                # Local SQLite DB management
│   │   ├── db.go               # Connection & Schema migrations
│   │   ├── session.go          # Saving/Loading MTProto session strings
│   │   └── messages.go         # Saving messages for Full-Text Search
│   └── config/                 # Environment and default dir handling (`~/.tgcli`)
├── .goreleaser.yaml            # Config for automated releases
├── go.mod
└── go.sum


4. Step-by-Step Implementation Phases (For Kiro Agent)

Phase 1: Project Setup & CLI Routing (Cobra)

Initialize the go module: go mod init github.com/virat-mankali/telegram-cli

Install Cobra: go get github.com/spf13/cobra@latest

Set up cmd/tgcli/main.go to simply call Execute() from root.go.

In root.go, define the base command and persistent flags (e.g., --store to override the default ~/.tgcli directory).

Create empty boilerplate files (auth.go, sync.go, send.go) and register them as subcommands to root.go using rootCmd.AddCommand().

Phase 2: Telegram Client & Authentication (gotd/td)

Install gotd: go get github.com/gotd/td

In internal/config/, establish the default data directory (~/.tgcli/).

In internal/telegram/client.go, initialize the gotd client using App ID and App Hash (can be hardcoded to generic open-source ones or fetched via ENV).

Implement tgcli auth:

Use auth.NewFlow from github.com/gotd/td/telegram/auth.

Implement a terminal prompt to ask the user for their phone number, the verification code sent by Telegram, and optionally their 2FA password.

Save the successful session string to internal/storage/session.go so subsequent commands don't require re-authentication.

Phase 3: Local Storage & SQLite FTS

Choose a SQLite driver (e.g., github.com/mattn/go-sqlite3).

In internal/storage/db.go, initialize a local database file at ~/.tgcli/data.db.

Create two main tables:

session: To securely store the gotd session token.

messages: A table utilizing SQLite's FTS5 extension. Columns should include id, chat_id, sender, text, timestamp.

Phase 4: Core Features (Sync, Search, Send)

tgcli sync:

Connect to MTProto.

Use tg.NewUpdateDispatcher() from gotd to listen for UpdateNewMessage.

When a new message arrives, insert it into the SQLite messages table.

tgcli search <query>:

Query the SQLite FTS5 table: SELECT * FROM messages WHERE text MATCH ?.

Format the output nicely in the terminal (consider using a library like github.com/fatih/color).

tgcli send --to <username/phone> --text <msg>:

Connect to MTProto.

Resolve the target peer (convert the username/phone into an MTProto InputPeer).

Dispatch the message using client.API().MessagesSendMessage().

Phase 5: Build & Distribution (GoReleaser)

Create a basic .goreleaser.yaml file in the root directory.

Configure it to build binaries for darwin (macOS), linux, and windows (amd64 and arm64).

Configure the brews section to point to a secondary GitHub repository named homebrew-tap (e.g., github.com/<your-username>/homebrew-tap).

Set up a GitHub Actions workflow to trigger GoReleaser automatically whenever a new Git tag (e.g., v1.0.0) is pushed.

5. Instructions for Kiro / AI Agent

Agent Directive: You are acting as a Senior Staff Go Engineer. Your task is to implement the phases above sequentially.

Rule 1: Write clean, modular Go code. Do not put all logic in main.go. Strictly respect the cmd/ and internal/ structure.

Rule 2: Use graceful context cancellation (signal.NotifyContext(context.Background(), os.Interrupt)).

Rule 3: Handle MTProto flood waits and rate limits gracefully using gotd built-in middlewares.

Rule 4: When instructed to start, begin with Phase 1 and provide the complete code for go.mod, main.go, and root.go. Wait for the user to confirm before moving to Phase 2.