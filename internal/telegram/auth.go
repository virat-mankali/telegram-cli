package telegram

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// TerminalAuth implements auth.UserAuthenticator for interactive CLI login.
type TerminalAuth struct {
	PhoneNumber string // pre-set via --phone flag (optional)
}

func (a TerminalAuth) Phone(_ context.Context) (string, error) {
	if a.PhoneNumber != "" {
		return a.PhoneNumber, nil
	}
	fmt.Print("Enter phone number (e.g. +1234567890): ")
	phone, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(phone), nil
}

func (TerminalAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter verification code: ")
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

func (TerminalAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // newline after hidden input
	return strings.TrimSpace(string(bytePwd)), nil
}

func (TerminalAuth) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("sign up not supported — tgcli only works with existing accounts")
}

func (TerminalAuth) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// NewAuthFlow creates an auth.Flow with the terminal authenticator.
func NewAuthFlow(phone string) auth.Flow {
	return auth.NewFlow(
		TerminalAuth{PhoneNumber: phone},
		auth.SendCodeOptions{},
	)
}
