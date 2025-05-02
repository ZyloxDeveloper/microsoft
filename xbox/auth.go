package xbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type XBLAuthConnect struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

type XBLAuthPoll struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	UserID           string `json:"user_id"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
}

func beginXBLAuth(s chan(*oauth2.Token)) string {
	d, err := startXBLAuth()
	if err != nil {
		return err.Error()
	}

	go startXBLPolling(d, s)

	return d.UserCode
}

func startXBLPolling(d *XBLAuthConnect, s chan(*oauth2.Token)) {
	ticker := time.NewTicker(time.Second * time.Duration(d.Interval))
	defer ticker.Stop()

	for range ticker.C {
		t, err := pollXBLAuth(d.DeviceCode)
		if err != nil {
			panic(err)
		}
		if t != nil {
			s <- t
			return
		}
	}
}

func startXBLAuth() (*XBLAuthConnect, error) {
	resp, err := http.PostForm("https://login.live.com/oauth20_connect.srf", url.Values{
		"client_id":     {"0000000048183522"},
		"scope":         {"service::user.auth.xboxlive.com::MBI_SSL"},
		"response_type": {"device_code"},
	})
	if err != nil {
		return nil, fmt.Errorf("POST https://login.live.com/oauth20_connect.srf: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("POST https://login.live.com/oauth20_connect.srf: %v", resp.Status)
	}
	data := new(XBLAuthConnect)
	return data, json.NewDecoder(resp.Body).Decode(data)
}

func pollXBLAuth(deviceCode string) (t *oauth2.Token, err error) {
	resp, err := http.PostForm(microsoft.LiveConnectEndpoint.TokenURL, url.Values{
		"client_id":   {"0000000048183522"},
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {deviceCode},
	})
	if err != nil {
		return nil, fmt.Errorf("POST https://login.live.com/oauth20_token.srf: %w", err)
	}
	poll := new(XBLAuthPoll)
	if err := json.NewDecoder(resp.Body).Decode(poll); err != nil {
		return nil, fmt.Errorf("POST https://login.live.com/oauth20_token.srf: json decode: %w", err)
	}
	_ = resp.Body.Close()
	if poll.Error == "authorization_pending" {
		return nil, nil
	} else if poll.Error == "" {
		return &oauth2.Token{
			AccessToken:  poll.AccessToken,
			TokenType:    poll.TokenType,
			RefreshToken: poll.RefreshToken,
			Expiry:       time.Now().Add(time.Duration(poll.ExpiresIn) * time.Second),
		}, nil
	}
	return nil, fmt.Errorf("%v: %v", poll.Error, poll.ErrorDescription)
}