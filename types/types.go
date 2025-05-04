package types

import (
	"time"

	"github.com/zyloxdeveloper/mailtracker"
)

type Account struct {
	Email    string
	Password string
	First    string
	Last     string
	Birthday time.Time
}

type AccountConfig struct {
	Mail   mailtracker.TrackerConfig
	Domain string
}
