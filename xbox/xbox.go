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

	c := make(chan *oauth2.Token, 0)
	errChan := make(chan error, 0)

	code := beginXBLAuth(c)

	go func() {
		if err := submitRemoteConnectCode(acc, code); err != nil {
			errChan <- err
			return
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return nil, err
	case token := <-c:
		if token == nil {
			return nil, fmt.Errorf("token corrupted")
		}
		return token, nil
	case <-time.After(2 * time.Minute):
		return nil, fmt.Errorf("timeout waiting for Xbox token or error")
	}
}