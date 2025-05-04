package microsoft

import (
	"time"

	"github.com/zyloxdeveloper/mailtracker"
	"github.com/zyloxdeveloper/microsoft/account"
	"github.com/zyloxdeveloper/microsoft/types"
	"golang.org/x/oauth2"
)

// NewConfig returns a new AccountConfig with domain and mail settings.
func NewConfig(domain string, mailCfg mailtracker.TrackerConfig) *types.MicrosoftConfig {
	return &types.MicrosoftConfig{
		Domain: domain,
		Mail:   mailCfg,
	}
}

// NewMailConfig returns a pre-filled mailtracker.TrackerConfig with default intervals,
// using the provided IMAP server, email address, and password.
// https://github.com/ZyloxDeveloper/mailtracker
func NewMailConfig(imapserver, email, password string) mailtracker.TrackerConfig {
	return mailtracker.TrackerConfig{
		IMAPServer: imapserver,
		EmailAddress: email,
		EmailPassword: password,
		CheckInterval: time.Second,
		CacheInterval: time.Second*10,
	}
}

// NewManager returns a new instance of the account manager using the given config.
func NewManager(cfg *types.MicrosoftConfig) *account.AccountManager {
	return account.New(cfg)
}

// CreateAccount generates a Microsoft account using the given manager.
func CreateAccount(m *account.AccountManager) (*types.Account, error) {
	return m.NewAccount()
}

// CreateXboxAccount generates a Microsoft account and returns the Xbox auth token.
func CreateXboxAccount(m *account.AccountManager) (*types.Account, *oauth2.Token, error) {
	return m.NewXboxAccount()
}
