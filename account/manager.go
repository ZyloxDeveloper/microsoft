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
)

type AccountManager struct {
	config   *types.AccountConfig
	mail     *mailtracker.Tracker
	accounts []*types.Account
}

func New(cfg *types.AccountConfig) *AccountManager {
	return &AccountManager{
		config: cfg,
		mail:   mailtracker.NewTracker(cfg.MailConfig),
	}
}

func (m *AccountManager) NewAccount() (*types.Account, error) {
	ctx, cancel := chrome.SetupChromedpContext()
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
	email := fmt.Sprintf("%s@%s", emailPrefix, m.config.EmailDomain)
	password := randomString(12, passwordCharset)

	return &types.Account{
		Email:    email,
		Password: password,
		First:    "John",
		Last:     "Doe",
		Birthday: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}, nil
}

func extractCode(body string) string {
	re := regexp.MustCompile(`\b\d{4,8}\b`)
	return re.FindString(body)
}

func randomString(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		builder.WriteByte(charset[rand.Intn(len(charset))])
	}
	return builder.String()
}