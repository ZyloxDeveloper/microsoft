package account

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/zyloxdeveloper/mailtracker"
)

func (m *AccountManager) setupEmailListener() <-chan string {
	codeChan := make(chan string, 1)
	var once sync.Once

	go m.mail.Start(func(email mailtracker.Email) {
		code := extractCode(email.Body)
		if code != "" {
			once.Do(func() {
				codeChan <- code
				m.mail.Stop()
			})
		}
	})

	return codeChan
}

var codeRegex = regexp.MustCompile(`(?i)(?:code|otp|verification)[^\d]{0,10}(\d{6})`)
func extractCode(body string) string {
	matches := codeRegex.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (m *AccountManager) waitForVerificationCode(codeChan <-chan string) (string, error) {
	select {
	case code := <-codeChan:
		return code, nil
	case <-time.After(2 * time.Minute):
		return "", fmt.Errorf("timeout waiting for verification email")
	}
}
