package account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/zyloxdeveloper/microsoft/types"
)

func (m *AccountManager) fillSignupForm(ctx context.Context, acc *types.Account) error {
	return chromedp.Run(ctx,
		chromedp.Navigate("https://signup.live.com/?lic=1"),
		chromedp.WaitVisible(`#usernameInput`, chromedp.ByID),
		chromedp.SendKeys(`#usernameInput`, acc.Email),
		chromedp.Click(`#nextButton`, chromedp.ByID),
		chromedp.WaitVisible(`#Password`, chromedp.ByID),
		chromedp.SendKeys(`#Password`, acc.Password),
		chromedp.Click(`#nextButton`, chromedp.ByID),
		chromedp.WaitVisible(`#firstNameInput`, chromedp.ByID),
		chromedp.SendKeys(`#firstNameInput`, acc.First),
		chromedp.SendKeys(`#lastNameInput`, acc.Last),
		chromedp.Click(`#nextButton`, chromedp.ByID),
		chromedp.WaitVisible(`#BirthMonth`, chromedp.ByID),
		chromedp.SetValue(`#BirthMonth`, fmt.Sprintf("%d", acc.Birthday.Month())),
		chromedp.SetValue(`#BirthDay`, fmt.Sprintf("%d", acc.Birthday.Day())),
		chromedp.SendKeys(`#BirthYear`, fmt.Sprintf("%d", acc.Birthday.Year())),
		chromedp.Click(`#nextButton`, chromedp.ByID),
	)
}

func (m *AccountManager) waitForManualCaptchaSolve(ctx context.Context) error {
	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				var url string
				if err := chromedp.Location(&url).Do(ctx); err != nil {
					return err
				}
				if strings.Contains(url, "privacynotice.account.microsoft.com") {
					return nil
				}
				time.Sleep(1 * time.Second)
			}
		}),
	)
}

