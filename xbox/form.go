package xbox

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/zyloxdeveloper/microsoft/chrome"
	"github.com/zyloxdeveloper/microsoft/types"
)

func registerXboxProfile(acc *types.Account) error {
	ctx, cancel := chrome.SetupChromedpContext()
	defer cancel()

	return chromedp.Run(ctx,
		chromedp.Navigate("https://www.xbox.com/en-US/auth/msa?action=logIn&returnUrl=https%3A%2F%2Fwww.xbox.com%2Fen-US&ru=https%3A%2F%2Fwww.xbox.com%2Fen-US"),

		chromedp.WaitVisible(`#usernameEntry`, chromedp.ByID),
		chromedp.SendKeys(`#usernameEntry`, acc.Email, chromedp.ByID),
		chromedp.Click(`[data-testid="primaryButton"]`, chromedp.ByQuery),

		chromedp.WaitVisible(`#passwordEntry`, chromedp.ByID),
		chromedp.SendKeys(`#passwordEntry`, acc.Password, chromedp.ByID),
		chromedp.Click(`[data-testid="primaryButton"]`, chromedp.ByQuery),

		chromedp.WaitVisible(`[data-testid="secondaryButton"]`, chromedp.ByQuery),
		chromedp.Click(`[data-testid="secondaryButton"]`, chromedp.ByQuery),

		chromedp.WaitVisible(`#create-account-gamertag-suggestion-3`, chromedp.ByID),
		chromedp.Sleep(3*time.Second),

		chromedp.WaitVisible(`#inline-continue-control`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 400; i++ {
				switch(i) {
					case 0: chromedp.Click(`#create-account-gamertag-suggestion-1`, chromedp.ByID); break;
					case 100: chromedp.Click(`#create-account-gamertag-suggestion-2`, chromedp.ByID); break;
					case 200: chromedp.Click(`#create-account-gamertag-suggestion-3`, chromedp.ByID); break;
					case 300: chromedp.Click(`#create-account-gamertag-suggestion-4`, chromedp.ByID); break;
				}
				var disabled bool
				err := chromedp.Evaluate(`document.getElementById("inline-continue-control")?.disabled === true`, &disabled).Do(ctx)
				if err != nil {
					return err
				}
				if !disabled {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			return chromedp.Click(`#inline-continue-control`, chromedp.ByID).Do(ctx)
		}),

		chromedp.WaitVisible(`#inline-continue-control`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 100; i++ {
				var disabled bool
				err := chromedp.Evaluate(`document.getElementById("inline-continue-control")?.disabled === true`, &disabled).Do(ctx)
				if err != nil {
					return err
				}
				if !disabled {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			return chromedp.Click(`#inline-continue-control`, chromedp.ByID).Do(ctx)
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 100; i++ {
				var currentURL string
				if err := chromedp.Location(&currentURL).Do(ctx); err != nil {
					return err
				}
				if currentURL == "https://www.xbox.com/en-US" {
					return nil
				}
				time.Sleep(100 * time.Millisecond)
			}
			return fmt.Errorf("timeout on redirect to xbox")
		}),
	)
}


func submitRemoteConnectCode(acc *types.Account, code string) error {
	ctx, close := chrome.SetupChromedpContext()
	defer close()

	c := chromedp.Run(ctx,
		chromedp.Navigate("https://login.live.com/oauth20_remoteconnect.srf"),
		chromedp.WaitVisible(`#otc`, chromedp.ByID),
		chromedp.SendKeys(`#otc`, code, chromedp.ByID),
		chromedp.Click(`#idSIButton9`, chromedp.ByID),

		chromedp.WaitVisible(`#i0116`, chromedp.ByID),
		chromedp.SendKeys(`#i0116`, acc.Email, chromedp.ByID),
		chromedp.Click(`#idSIButton9`, chromedp.ByID),

		chromedp.WaitVisible(`#i0118`, chromedp.ByID),
		chromedp.SendKeys(`#i0118`, acc.Password, chromedp.ByID),
		chromedp.Click(`#idSIButton9`, chromedp.ByID),

		chromedp.WaitVisible(`#declineButton`, chromedp.ByID),
		chromedp.Click(`#declineButton`, chromedp.ByID),

		chromedp.Sleep(time.Second),
	)

	return c
}
