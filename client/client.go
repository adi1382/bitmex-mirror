package client

import (
	"time"

	"github.com/bitmex-mirror/auth"
	"github.com/juju/ratelimit"
	"github.com/valyala/fasthttp"
)

func NewClient(apiKey, secret string, test bool,
	bucket10s, bucket1m *ratelimit.Bucket,
	httpClient *fasthttp.Client) *Client {

	c := Client{
		auth:      auth.NewConfig(apiKey, secret, test),
		bucket10s: bucket10s,
		bucket1m:  bucket1m,
		httpC:     httpClient,
	}

	c.httpTimeout = time.Second * 5

	return &c
}

type Client struct {
	auth        *auth.Config
	bucket10s   *ratelimit.Bucket
	bucket1m    *ratelimit.Bucket
	httpC       *fasthttp.Client
	httpTimeout time.Duration
}
