package xbox

import (
	"github.com/zyloxdeveloper/microsoft/types"
	"golang.org/x/oauth2"
)

func XBLToken(acc *types.Account) (*oauth2.Token, error) {
	if err := registerXboxProfile(acc); err != nil {
		return nil, err
	}

	c := make(chan *oauth2.Token)
	code := beginXBLAuth(c)

	if err := submitRemoteConnectCode(acc, code); err != nil {
		return nil, err
	}

	t := <- c
	return t, nil
}