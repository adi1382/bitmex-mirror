package bitmex

import (
	"time"

	"golang.org/x/time/rate"

	"github.com/bitmex-mirror/auth"
	"github.com/valyala/fasthttp"
)

type Error string

func (e Error) Error() string { return string(e) }

const ErrServerOverloaded = Error("server overloaded")
const ErrBadRequest = Error("bad request")
const ErrInvalidAPIKey = Error("invalid API key")
const ErrTooManyRequests = Error("too many requests")
const ErrWSConnClosed = Error("websocket connection closed")
const ErrWSClientClosed = Error("websocket client closed")
const ErrWSVerificationTimeout = Error("context canceled/timeout: message sent but could not verify")
const ErrUnexpectedError = Error("unexpected error")
const ErrRequestExpired = Error("request expired")
const ErrClientError = Error("client error (400<=code<500)")
const ErrServerError = Error("client error (500<=code<600)")

const requestTimeout = time.Second * 10              // Used to create API Expires
const wsConfirmTimeout = time.Second * 15            // time to wait for confirmation message on websocket send
const restEndpoint = "www.bitmex.com/api/v1"         // rest endpoint for mainnet
const restTestEndpoint = "testnet.bitmex.com/api/v1" // rest endpoint for testnet
const wsEndpoint = "ws.bitmex.com"                   // websocket endpoint for mainnet
const wsTestEndpoint = "ws.testnet.bitmex.com"       // websocket endpoint for testnet
const rateLimitMinuteCap = 120                       // default capacity for one minute request limiter
const rateLimitSecondCap = 10                        // default capacity for one second request limiter
const websocketConnectionTimeout = time.Second * 15  // timeout when attempting to establish websocket connection
const pingPeriod = time.Second * 5                   // time to ping from receiving the last message over websocket
const pongTimeout = time.Second * 3                  // wait time to receive pong after sending ping

const OrderStatusNew = "New"
const OrderStatusPartiallyFilled = "PartiallyFilled"
const OrderStatusFilled = "Filled"
const OrderStatusCanceled = "Canceled"
const OrderStatusRejected = "Rejected"

const OrderSideBuy = "Buy"
const OrderSideSell = "Sell"

const OrderTypeLimit = "Limit"
const OrderTypeMarket = "Market"
const OrderTypeStop = "Stop"
const OrderTypeStopLimit = "StopLimit"
const OrderTypeStopMarket = "StopMarket"
const OrderTypeTrailingStop = "TrailingStop"

const WSDataActionPartial = "partial"
const WSDataActionUpdate = "update"
const WSDataActionInsert = "insert"
const WSDataActionDelete = "delete"

const WSTableOrder = "order"
const WSTablePosition = "position"
const WSTableMargin = "margin"

func NewBitmex(test bool) *Bitmex {
	b := Bitmex{
		bucketM: rate.NewLimiter(rate.Every(time.Minute/rateLimitMinuteCap), rateLimitMinuteCap),
		bucketS: rate.NewLimiter(rate.Every(time.Second/rateLimitSecondCap), rateLimitSecondCap),
		test:    test,
		httpC: &fasthttp.Client{
			Name:                          "",
			NoDefaultUserAgentHeader:      false,
			Dial:                          nil,
			DialDualStack:                 false,
			TLSConfig:                     nil,
			MaxConnsPerHost:               512,               // Maximum number of concurrent requests
			MaxIdleConnDuration:           time.Second * 150, // HTTP-Persistence
			MaxConnDuration:               0,
			MaxIdemponentCallAttempts:     0,
			ReadBufferSize:                0,
			WriteBufferSize:               0,
			ReadTimeout:                   0,
			WriteTimeout:                  0,
			MaxResponseBodySize:           0,
			DisableHeaderNamesNormalizing: false,
			DisablePathNormalizing:        false,
			MaxConnWaitTimeout:            time.Second * 5, // Http timeout
			RetryIf:                       nil,
		},
	}
	return &b
}

type Bitmex struct {
	bucketM *rate.Limiter
	bucketS *rate.Limiter
	test    bool
	httpC   *fasthttp.Client
}

// NewRestClient provides a variable of type *RestClient which provides a method set to call
// relevant rest APIs
func (b *Bitmex) NewRestClient(config auth.Config) *RestClient {

	c := RestClient{
		auth:    config,
		bucketM: b.bucketM,
		bucketS: b.bucketS,
		httpC:   b.httpC,
	}
	if b.test {
		c.endpoint = restTestEndpoint
	} else {
		c.endpoint = restEndpoint
	}
	return &c
}
