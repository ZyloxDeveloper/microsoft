package chrome

import (
	"context"

	"github.com/chromedp/chromedp"
)

func SetupChromedpContext() (context.Context, context.CancelFunc) {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
	)...)
	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	return ctx, func() {
		cancelCtx()
		cancel()
	}
}