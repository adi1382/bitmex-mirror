package bitmex

import (
	"github.com/bitmex-mirror/auth"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
)

type RestClient struct {
	auth     auth.Config
	endpoint string
	bucketM  *rate.Limiter
	bucketS  *rate.Limiter
	httpC    *fasthttp.Client
}
