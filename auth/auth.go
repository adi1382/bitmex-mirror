package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Config struct {
	Key    string // Key is exported
	secret string // secret is unexported and is only meant to be used for generating signature
}

// NewConfig returns new Config variable and provides a method to Sign
func NewConfig(key, secret string) Config {
	config := Config{
		Key:    key,
		secret: secret,
	}

	return config
}

// Sign a message with secret key of the Config variable
func (p Config) Sign(body string) string {
	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}
