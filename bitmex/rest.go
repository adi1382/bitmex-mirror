package bitmex

import (
	"time"

	"github.com/juju/ratelimit"
	"github.com/valyala/fasthttp"
)

func NewRestConnection(test bool) *rest {
	r := rest{
		clients:   nil,
		bucket10s: ratelimit.NewBucket(time.Second*10, 10),
		bucket1m:  ratelimit.NewBucket(time.Second, 60),
		http: &fasthttp.Client{
			MaxIdleConnDuration: time.Second * 150,
		},
	}

	if test {
		r.endpoint = "https://testnet.bitmex.com/api/v1"
	} else {
		r.endpoint = "https://www.bitmex.com/api/v1"
	}

	return &r
}

type rest struct {
	clients   []*restClient
	bucket10s *ratelimit.Bucket
	bucket1m  *ratelimit.Bucket
	http      *fasthttp.Client
	endpoint  string
}

func (r *rest) AddNewClient(apiKey, secret string) *restClient {
	c := newClient(apiKey, secret, r.endpoint, r.bucket10s, r.bucket1m, r.http)
	r.clients = append(r.clients, c)
	return c
}
