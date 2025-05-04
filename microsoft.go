package microsoft

import (
	"github.com/zyloxdeveloper/microsoft/account"
	"github.com/zyloxdeveloper/microsoft/types"
	"golang.org/x/oauth2"
)

type (
	Account       = types.Account
	AccountConfig = types.AccountConfig
	AccountManager = account.AccountManager
)

// NewManager returns a new instance of the account manager using the given config.
func NewManager(cfg *AccountConfig) *AccountManager {
	return account.New(cfg)
}

// CreateAccount generates a Microsoft account using the given manager.
func CreateAccount(m *AccountManager) (*Account, error) {
	return m.NewAccount()
}

// CreateXboxAccount generates a Microsoft account and returns the Xbox auth token.
func CreateXboxAccount(m *AccountManager) (*Account, *oauth2.Token, error) {
	return m.NewXboxAccount()
}
