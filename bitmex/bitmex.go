package bitmex

import (
	"context"
	"github.com/bitmex-mirror/auth"
	"github.com/pkg/errors"
	"log"
	"time"

	"golang.org/x/time/rate"

	"github.com/valyala/fasthttp"
)

const RequestTimeout = time.Second * 10              // Used to create API Expires
const wsConfirmTimeout = time.Second * 15            // time to wait for confirmation message on websocket send
const restEndpoint = "www.bitmex.com/api/v1"         // rest endpoint for mainnet
const restTestEndpoint = "testnet.bitmex.com/api/v1" // rest endpoint for testnet
const wsEndpoint = "ws.bitmex.com"                   // websocket endpoint for mainnet
const wsTestEndpoint = "ws.testnet.bitmex.com"       // websocket endpoint for testnet
const rateLimitMinuteCap = 120                       // capacity for one minute request limits
const rateLimitSecondCap = 10                        // capacity for one second request limits on certain routes
const websocketConnectionTimeout = time.Second * 15  // timeout when attempting to establish websocket connection
const pingPeriod = time.Second * 5                   // time to ping from receiving the last message over websocket
const pongTimeout = time.Second * 3                  // wait time to receive pong after sending ping

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

// NewRestClient provides a variable of type *restClient which provides a method set to call
// relevant rest APIs
func (b *Bitmex) NewRestClient(config auth.Config) *restClient {

	c := restClient{
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

// NewWSConnection create a new multiplexed websocket connection and returns a variable to provide method set to
// add and authenticate several account on the connection.
func (b *Bitmex) NewWSConnection(ctx context.Context, logger *log.Logger) (*wSConnection, error) {
	var ws wSConnection

	if b.test {
		ws.host = wsTestEndpoint
	} else {
		ws.host = wsEndpoint
	}

	ws.bucketM = b.bucketM

	ws.done = make(chan struct{})

	ctx2, cancel := context.WithTimeout(ctx, websocketConnectionTimeout)
	defer cancel()

	conn, err := connect(ctx2, ws.host)
	if err != nil {
		close(ws.done)
		return nil, errors.Wrap(err, "ws connection failed to connect")
	}
	ws.conn = conn

	chWrite := make(chan []byte, 100)
	chRead := make(chan []byte, 100)
	ws.clients = make(map[string]*wsClient, 5)
	ws.removeClient = make(chan string)

	ws.chRead = chRead
	ws.chWrite = chWrite
	ws.pingPeriod = pingPeriod
	ws.pongTimeout = pongTimeout

	if ws.pingPeriod < ws.pongTimeout {
		panic("ping period cannot be smaller than pong timeout")
	}

	ws.pingTicker = time.NewTicker(ws.pingPeriod)

	ws.pongTimer = time.NewTimer(time.Hour)
	if !ws.pongTimer.Stop() {
		<-ws.pongTimer.C
	}

	ws.conn.SetPongHandler(func(string) error {
		if !ws.pongTimer.Stop() {
			_ = ws.conn.Close()
		}
		return nil
	})

	// Reader will start writer and router
	go ws.connectionReader(ctx, chRead, chWrite, logger)

	return &ws, err
}
