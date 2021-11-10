package pool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bitmex-mirror/auth"
	"github.com/bitmex-mirror/bitmex"
	"github.com/pkg/errors"
)

func NewPool(test bool, logger *log.Logger) *Pool {
	p := Pool{}
	p.connectionCond = sync.NewCond(&sync.Mutex{})

	bit := bitmex.NewBitmex(test)
	p.bitmex = bit
	p.accounts = make(map[string]*account, 1)
	p.Removals = make(chan Removal)
	p.internalRemovals = make(chan Removal)

	ctx := context.Background()

	go p.manager(ctx, logger) // the manager will call connect

	return &p
}

func (p *Pool) connect(ctx context.Context, logger *log.Logger) {
	connection, err := p.bitmex.NewWSConnection(ctx, logger)
	if err == nil {
		p.wsConn = connection
		return
	}
	time.Sleep(time.Second)
	fmt.Println("Error connecting to Bitmex Websocket:", err)
	logger.Println("Failed to connect to Bitmex Websocket:", err)
	p.connect(ctx, logger)
}

type Pool struct {
	Removals         chan Removal
	internalRemovals chan Removal
	bitmex           *bitmex.Bitmex
	wsConn           *bitmex.WSConnection
	accMu            sync.RWMutex
	accounts         map[string]*account
	connectionCond   *sync.Cond
}

type Removal struct {
	APIKey string
	Reason error
}

func (p *Pool) Put(ctx context.Context, config auth.Config, topic string, logger *log.Logger, id string) {
	acc, err := p.newAccount(ctx, config, topic, logger, id)
	if err != nil {
		p.Removals <- Removal{config.Key, err}
		logger.Printf("Error creating account: %s", err)
		return
	}
	p.accMu.Lock()
	//p.accounts[config.Key] = acc
	p.accounts[id] = acc
	p.accMu.Unlock()
}

func (p *Pool) Get(api string) (*account, error) {
	p.accMu.RLock()
	defer p.accMu.RUnlock()

	var acc *account
	var ok bool
	if acc, ok = p.accounts[api]; !ok {
		return nil, errors.New("account not found")
	}
	return acc, nil
}

func (p *Pool) manager(ctx context.Context, logger *log.Logger) {
	p.connectionCond.L.Lock()
	p.connect(ctx, logger)
	p.connectionCond.Broadcast()
	p.connectionCond.L.Unlock()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-p.internalRemovals:
			p.accMu.Lock()
			delete(p.accounts, r.APIKey)
			p.accMu.Unlock()
			logger.Printf("Removed account: %s Reason: %s", r.APIKey, r.Reason)
			p.Removals <- r
		case <-p.wsConn.Done:
			logger.Println("Connection to Bitmex Websocket closed")
			logger.Println("Restarting websocket connection")
			p.connectionCond.L.Lock()
			p.connect(ctx, logger)
			logger.Println("websocket connection restored")
			p.connectionCond.Broadcast()
			logger.Println("Broadcasted")
			p.connectionCond.L.Unlock()
		}
	}
}

// newWSClient is used by underlying account variables to create their own wsClient on the common
// websocket connection.
// all parameters value are the same as that of Put function which is used by the user at the time of creation of
// new account
func (p *Pool) newWSClient(ctx context.Context, config auth.Config, topic string, logger *log.Logger, id string) *bitmex.WSClient {
	p.connectionCond.L.Lock()
	for p.wsConn == nil || !p.wsConn.IsOpen() {
		p.connectionCond.Wait() // wait for the websocket connection to be established by pool's manager
	}
	client, err := p.wsConn.NewWSClient(ctx, config, topic, logger, id)
	p.connectionCond.L.Unlock()

	if err != nil {
		return p.newWSClient(ctx, config, topic, logger, id)
	}
	return client
}
