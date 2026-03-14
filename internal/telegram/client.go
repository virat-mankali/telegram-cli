// Package telegram wraps gotd/td for MTProto connection management.
package telegram

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

// ClientOptions holds configuration for creating a Telegram client.
type ClientOptions struct {
	SessionPath string
	Logger      *zap.Logger
}

// AppCredentials returns app ID and hash from env vars, falling back to
// gotd test credentials for development.
func AppCredentials() (int, string) {
	appID := telegram.TestAppID
	appHash := telegram.TestAppHash

	if id := os.Getenv("TGCLI_APP_ID"); id != "" {
		if parsed, err := strconv.Atoi(id); err == nil {
			appID = parsed
		}
	}
	if hash := os.Getenv("TGCLI_APP_HASH"); hash != "" {
		appHash = hash
	}
	return appID, appHash
}

// NewClient creates a configured gotd Telegram client with flood wait
// and rate limit middleware.
func NewClient(opts ClientOptions) *telegram.Client {
	appID, appHash := AppCredentials()

	logger := opts.Logger
	if logger == nil {
		logger = zap.NewNop()
	}

	return telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: &session.FileStorage{Path: opts.SessionPath},
		Logger:         logger,
		Middlewares: []telegram.Middleware{
			floodwait.NewSimpleWaiter().WithMaxRetries(10),
			ratelimit.New(rate.Every(100*time.Millisecond), 5),
		},
	})
}
