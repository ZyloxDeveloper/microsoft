package account

import (
	"fmt"
	"time"

	"github.com/zyloxdeveloper/mailtracker"
)

func (m *AccountManager) setupEmailListener() <-chan string {
	codeChan := make(chan string, 1)
	go m.mail.Start(func(email mailtracker.Email) {
		code := extractCode(email.Body)
		if code != "" {
			codeChan <- code
			m.mail.Stop()
		}
	})
	return codeChan
}

func (m *AccountManager) waitForVerificationCode(codeChan <-chan string) (string, error) {
	select {
	case code := <-codeChan:
		return code, nil
	case <-time.After(2 * time.Minute):
		return "", fmt.Errorf("timeout waiting for verification email")
	}
}
