package account

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/zyloxdeveloper/mailtracker"
	"github.com/zyloxdeveloper/microsoft/chrome"
	"github.com/zyloxdeveloper/microsoft/types"
	"github.com/zyloxdeveloper/microsoft/xbox"
	"golang.org/x/oauth2"
)

type AccountManager struct {
	config   *types.AccountConfig
	mail     *mailtracker.Tracker
	accounts []*types.Account
}

func New(cfg *types.AccountConfig) *AccountManager {
	return &AccountManager{
		config: cfg,
		mail:   mailtracker.NewTracker(cfg.Mail),
	}
}

func (m *AccountManager) NewAccount() (*types.Account, error) {
	ctx, cancel := chrome.SetupChromedpContext(false)
	defer cancel()

	acc, err := m.randomAccount()
	if err != nil {
		return nil, err
	}

	codeChan := m.setupEmailListener()

	if err := m.fillSignupForm(ctx, acc); err != nil {
		return nil, err
	}

	code, err := m.waitForVerificationCode(codeChan)
	if err != nil {
		return nil, err
	}

	if err := m.submitVerificationCode(ctx, code); err != nil {
		return nil, err
	}

	if err := m.waitForManualCaptchaSolve(ctx); err != nil {
		return nil, err
	}

	return acc, nil
}

func (m *AccountManager) NewXboxAccount() (*types.Account, *oauth2.Token, error) {
	acc, err := m.NewAccount()
	if err != nil {
		return nil, nil, err 
	}

	tok, err := xbox.XBLToken(acc)
	if err != nil {
		return acc, nil, err
	}

	if tok == nil {
		return acc, nil, fmt.Errorf("invalid token")
	}

	return acc, tok, nil
}

func (m *AccountManager) submitVerificationCode(ctx context.Context, code string) error {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(`#VerificationCode`, chromedp.ByID),
		chromedp.SendKeys(`#VerificationCode`, code),
		chromedp.Click(`#nextButton`, chromedp.ByID),
	)
	if err != nil {
		return fmt.Errorf("verification code submit error: %w", err)
	}
	return nil
}

func (m *AccountManager) randomAccount() (*types.Account, error) {
	const emailCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*?"

	emailPrefix := randomString(10, emailCharset)
	email := fmt.Sprintf("%s@%s", emailPrefix, m.config.Domain)
	password := randomString(12, passwordCharset)

	return &types.Account{
		Email:    email,
		Password: password,
		First:    "John",
		Last:     "Doe",
		Birthday: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}, nil
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var codeRegex = regexp.MustCompile(`\b\d{4,8}\b`)

func randomString(length int, charset string) string {
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		builder.WriteByte(charset[seededRand.Intn(len(charset))])
	}
	return builder.String()
}

func extractCode(body string) string {
	return codeRegex.FindString(body)
}