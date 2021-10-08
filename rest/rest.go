package rest

import (
	"time"

	"github.com/bitmex-mirror/client"
	"github.com/juju/ratelimit"
	"github.com/valyala/fasthttp"
)

func NewRestObject() *Rest {
	r := Rest{
		clients:   nil,
		bucket10s: ratelimit.NewBucket(time.Second*10, 10),
		bucket1m:  ratelimit.NewBucket(time.Minute, 60),
		http:      new(fasthttp.Client),
	}

	return &r
}

type Rest struct {
	clients   []*client.Client
	bucket10s *ratelimit.Bucket
	bucket1m  *ratelimit.Bucket
	http      *fasthttp.Client
}

func (r *Rest) AddNewClient(apiKey, secret string, test bool) *client.Client {
	c := client.NewClient(apiKey, secret, test, r.bucket10s, r.bucket1m, r.http)
	r.clients = append(r.clients, c)
	return c
}
