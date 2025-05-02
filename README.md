# Microsoft Account Token Generator

Example:
```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zyloxdeveloper/mailtracker"
	"github.com/zyloxdeveloper/microsoft/account"
	"github.com/zyloxdeveloper/microsoft/types"
	"github.com/zyloxdeveloper/microsoft/xbox"
)

func main() {
	mailConfig := mailtracker.TrackerConfig{
		IMAPServer:    "imap.yourmail.com:993",
		EmailAddress:  "your@email.com",
		EmailPassword: "your_app_password",
		CheckInterval: time.Second,
	}

	cfg := &types.AccountConfig{
		MailConfig:  mailConfig,
		EmailDomain: "yourdomain.org",
	}

	manager := account.New(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	acc, err := manager.NewAccount(ctx)
	if err != nil {
		panic(err)
	}

	tok, err := xbox.XBLToken(acc)
	if err != nil {
		panic(err)
	}

	fmt.Println(tok)
}

```