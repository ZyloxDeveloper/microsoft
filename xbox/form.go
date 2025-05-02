package xbox

import (
	"github.com/chromedp/chromedp"
	"github.com/zyloxdeveloper/microsoft/chrome"
	"github.com/zyloxdeveloper/microsoft/types"
)

func submitRemoteConnectCode( acc *types.Account, code string) error {
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
	)

	return c
}