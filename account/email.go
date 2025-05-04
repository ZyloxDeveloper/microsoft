package account

import (
	"fmt"
	"regexp"
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

var codeRegex = regexp.MustCompile(`\b\d{6,8}\b`)
func extractCode(body string) string {
	return codeRegex.FindString(body)
}

func (m *AccountManager) waitForVerificationCode(codeChan <-chan string) (string, error) {
	select {
	case code := <-codeChan:
		return code, nil
	case <-time.After(2 * time.Minute):
		return "", fmt.Errorf("timeout waiting for verification email")
	}
}
