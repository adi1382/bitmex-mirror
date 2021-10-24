package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Config struct {
	Key    string
	secret string
}

func NewConfig(key, secret string) Config {
	config := Config{
		Key:    key,
		secret: secret,
	}

	return config
}

func (p *Config) Sign(body string) string {
	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}
