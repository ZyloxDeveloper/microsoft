package chrome

import (
	"context"

	"github.com/chromedp/chromedp"
)

func SetupChromedpContext(headless bool) (context.Context, context.CancelFunc) {
	opts := chromedp.DefaultExecAllocatorOptions[:]
	if !headless {
		opts = append(opts,
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
		)
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	return ctx, func() {
		cancelCtx()
		cancel()
	}
}
