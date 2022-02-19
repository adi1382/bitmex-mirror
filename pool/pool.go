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
	p := Pool{
		closing:          make(chan string, 1),
		startConn:        make(chan struct{}, 1),
		wsClientReq:      make(chan wsClientParams, 1),
		Removals:         make(chan Removal),
		internalRemovals: make(chan Removal),
		bitmex:           bitmex.NewBitmex(test),
		accMu:            sync.RWMutex{},
		accounts:         make(map[string]*Account, 1),
	}

	done := make(chan struct{})
	p.Done = done

	p.startConn <- struct{}{}

	//p.bitmex = bitmex.NewBitmex(test)
	//p.accounts = make(map[string]*Account, 1)
	//p.Removals = make(chan Removal)
	//p.internalRemovals = make(chan Removal)
	//p.startConn = make(chan struct{}, 1)

	go p.manager(done, logger) // the manager will call connect

	return &p
}

type Pool struct {
	Done             <-chan struct{}
	closing          chan string
	startConn        chan struct{}
	wsClientReq      chan wsClientParams
	Removals         chan Removal
	internalRemovals chan Removal
	bitmex           *bitmex.Bitmex
	accMu            sync.RWMutex
	accounts         map[string]*Account
}

type wsClientParams struct {
	config   auth.Config
	topic    string
	id       string
	receiver chan []byte
	logger   *log.Logger
}

type Removal struct {
	id     string
	Reason error
}

func (p *Pool) Put(config auth.Config, topic string, id string, logger *log.Logger) {
	acc, err := p.newAccount(config, topic, id, logger)
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

func (p *Pool) Get(id string) (*Account, error) {
	p.accMu.RLock()
	defer p.accMu.RUnlock()

	acc, ok := p.accounts[id]
	if !ok {
		return nil, errors.New("account not found")
	}

	return acc, nil
}

func (p *Pool) manager(done chan struct{}, logger *log.Logger) {

	defer func() {
		close(done)
		return
	}()

	for {
		select {
		case by := <-p.closing:
			logger.Printf("websocket connection closed by: %s", by)
			return
		case <-p.startConn:
			wsConn := p.connect(logger)
			p.wsConnWorker(wsConn, logger)
		case r := <-p.internalRemovals:
			p.accMu.Lock()
			delete(p.accounts, r.id)
			p.accMu.Unlock()
			logger.Printf("Removed account: %s Reason: %s", r.id, r.Reason)
			p.Removals <- r
			//case <-p.wsConn.Done:
			//	logger.Println("Connection to Bitmex Websocket closed")
			//	logger.Println("Restarting websocket connection")
			//	p.connect(ctx, logger)
			//	logger.Println("websocket connection restored")
			//	logger.Println("Broadcasted")
		}
	}
}

func (p *Pool) wsConnWorker(wsConn *bitmex.WSConnection, logger *log.Logger) {
	defer wsConn.Close()

	for {
		select {
		case <-p.Done:
			// no need to restart in this case
			return
		case <-wsConn.Done:
			logger.Println("Connection to Bitmex Websocket closed")
			logger.Println("Restarting websocket connection")
			p.startConn <- struct{}{}
			return
		case req := <-p.wsClientReq:
			p.accMu.RLock()
			// check if the Account is still alive
			select {
			case <-p.accounts[req.id].Done:
				continue
			default:
			}

			// create a new websocket client for the Account
			wsClient, err := wsConn.SubscribeNewClient(req.config, req.topic, req.id, req.receiver, req.logger)
			if err != nil {
				// err means connection was closed
				// we need to restart the connection
				logger.Printf("Error subscribing to Bitmex Websocket: %s", req.id)
				// since the connection is closed, client cannot be constructed, so sending req back to channel
				// to be processed again by new connection
				p.startConn <- struct{}{}

				// sending the request back to the channel
				go func() {
					p.wsClientReq <- req
				}()

				return
			}
			// send the WSClient to the Account
			select {
			case p.accounts[req.id].newWSClient <- wsClient:
			case <-p.accounts[req.id].Done:
				//fmt.Println("this ran")
				// If Account is closed then unsubscribe the client from connection
				_ = wsClient.Unsubscribe(context.TODO())
			}
			p.accMu.RUnlock()
		}
	}
}

func (p *Pool) connect(logger *log.Logger) *bitmex.WSConnection {
	connection, err := p.bitmex.NewWSConnection(context.TODO(), logger)
	if err == nil {
		return connection
	}
	time.Sleep(time.Second)
	fmt.Println("Error connecting to Bitmex Websocket:", err)
	logger.Println("Failed to connect to Bitmex Websocket:", err)
	return p.connect(logger)
}

// newWSClient reassigns the wsClient variable to a new WSClient for the Account object.
// This reassignment is necessary because the WSClient is a singleton and once closed due to network issues, or
// closed by host, it cannot be reopened and must be recreated.
// It requires an active websocket connection on which this new client can be added, so if the websocket connection
// is not available , it will wait until it is.
// This function returns a channel that will be closed when wsClient is assigned.
// the Account caller must not use the wsClient until the channel is closed.
//func (p *Pool) newWSClient(config auth.Config,
//	topic string, id string, receiver chan []byte, logger *log.Logger) (*bitmex.WSClient, error) {
//
//	p.connCond.L.Lock()
//	for p.wsConn == nil || !p.wsConn.IsOpen() {
//		p.connCond.Wait() // wait for the websocket connection to be established by pool's manager
//	}
//	p.connCond.L.Unlock()
//
//	select {
//	case <-ctx.Done():
//		return nil, bitmex.ErrContextCanceled
//	default:
//		if client, err := p.wsConn.SubscribeNewClient(config, topic, id, receiver, logger); err == nil {
//			return client, nil
//		} else {
//			return p.newWSClient(config, topic, id, receiver, logger)
//		}
//	}
//}
