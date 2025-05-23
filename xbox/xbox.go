package xbox

import (
	"fmt"
	"time"

	"github.com/zyloxdeveloper/microsoft/types"
	"golang.org/x/oauth2"
)

func XBLToken(acc *types.Account) (*oauth2.Token, error) {
	if err := registerXboxProfile(acc); err != nil {
		return nil, err
	}

	d, err := startXBLAuth()
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 0)
	go func() {
		if err := submitRemoteConnectCode(acc, d.UserCode); err != nil {
			errChan <- err
			return
		}
	}()

	tokChan := make(chan *oauth2.Token, 0)
	go func() {
		tok, err := startXBLPolling(d)
		if err != nil {
			errChan <- err
			return
		}
		if tok == nil {
			return
		}
		tokChan <- tok
	}()

	select {
	case token := <-tokChan:
		return token, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(2 * time.Minute):
		return nil, fmt.Errorf("timeout waiting for Xbox token or error")
	}
}
