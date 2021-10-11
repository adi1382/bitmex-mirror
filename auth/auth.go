package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Config struct {
	Endpoint string
	Key      string
	secret   string
}

func NewConfig(key, secret string, test bool) Config {
	config := Config{
		Key:    key,
		secret: secret,
	}

	if test {
		config.Endpoint = "https://testnet.bitmex.com/api/v1"
	} else {
		config.Endpoint = "https://www.bitmex.com/api/v1"
	}

	return config
}

func (p *Config) Sign(body string) string {
	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}
