package bitmex

import (
	"time"

	"github.com/bitmex-mirror/auth"
	"github.com/juju/ratelimit"
	"github.com/valyala/fasthttp"
)

// testClient - only used for writing tests
func testClient(apiKey, secret string, test bool) *restClient {

	var endpoint string
	if test {
		endpoint = "https://testnet.bitmex.com/api/v1"
	} else {
		endpoint = "https://www.bitmex.com/api/v1"
	}

	c := newClient(
		apiKey, secret, endpoint,
		ratelimit.NewBucket(time.Second*10, 10), ratelimit.NewBucket(time.Second, 60),
		&fasthttp.Client{MaxIdleConnDuration: 150 * time.Second})
	return c
}

func newClient(apiKey, secret, endpoint string,
	bucket10s, bucket1m *ratelimit.Bucket,
	httpClient *fasthttp.Client) *restClient {

	c := restClient{
		auth:      auth.NewConfig(apiKey, secret),
		endpoint:  endpoint,
		bucket10s: bucket10s,
		bucket1m:  bucket1m,
		httpC:     httpClient,
	}

	c.httpTimeout = time.Second * 5

	return &c
}

type restClient struct {
	auth        auth.Config
	endpoint    string
	bucket10s   *ratelimit.Bucket
	bucket1m    *ratelimit.Bucket
	httpC       *fasthttp.Client
	httpTimeout time.Duration
}
