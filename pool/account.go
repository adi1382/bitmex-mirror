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

// Account represents a user Account.
type Account struct {
	//startMu sync.Mutex // Mutex to protect the startStreaming process

	startWSClient chan struct{}
	Done          <-chan struct{} // Channel to signal the end of the life of Account
	closing       chan string     // signals manager to close the Done channel
	//closeDone     sync.Once       // Once to close the Done channel

	config auth.Config // Config for the Account

	id    string //multi-plexed websocket id
	topic string // Topic for the Account (additional optional parameter for websocket)

	restClient *bitmex.RestClient // REST client for the Account for REST requests

	// The state of these slices are maintained by the Account with websockets and also updated at the time of
	// post, delete, put rest requests. They are protected by the respective mutexes.
	activeOrders []bitmex.Order    // current state of active orders for the Account
	positions    []bitmex.Position // current state of positions for the Account
	margins      []bitmex.Margin   // current state of margins for the Account
	ordersMu     sync.RWMutex      // Mutex to protect activeOrders
	positionsMu  sync.RWMutex      // Mutex to protect positions
	marginsMu    sync.RWMutex      // Mutex to protect margins

	// subs is used by other clients to subscribe to this Account, once subscribed all incoming messages from the
	// ws connection will be published to the subscribed channel.
	// This functionality is meant to be used by the sub-accounts to subscribe to the host-Account.
	//subs   []chan<- []byte // Channels to send data to other accounts subscribing to this Account
	//subsMu sync.RWMutex    // Mutex to protect subs slice

	// wsClient is reassigned from time to time by use of newWSClient() method of the pool.
	// the Account must ensure that wsClient variable is not used at the time of assignment.
	// It is only called by the startStreaming() method, which is protected by its mutex.
	//wsClient *bitmex.WSClient // Websocket client
	newWSClient chan *bitmex.WSClient // Channel to receive the new websocket client

	// true when respective partials are received
	// all booleans are modified under protection of partialsCond mutex
	// each time a new partial is received, the respective boolean is set to true, and the condition is broadcasted
	// subject to the conditions, the waiter can again go to the waiting state if the needed boolean is still false.
	partialsOrder    bool
	partialsPosition bool
	partialsMargin   bool
	partialsCond     *sync.Cond // Conditional variable to signal that partials were received

	// pool method newWSClient() is used directly to create new WSClient for Account, wsConnection
	// is not used in the Account variable. The reason being that wsConnection is again from time to time
	// by the pool, whenever network connection is lost. The pool ensures that the wsConn variable is always
	// protected at the time of new assignment, and hence not passed to Account.
	pool *Pool // Pool to which this Account belongs (used to create new wsClient)

	*pubSub // PubSub mechanism to publish messages to other accounts
}

func (p *Pool) newAccount(config auth.Config, topic string, id string, logger *log.Logger) (*Account, error) {
	acc := Account{
		startWSClient: make(chan struct{}, 1),
		closing:       make(chan string, 1),
		config:        config,
		id:            id,
		topic:         topic,
		restClient:    p.bitmex.NewRestClient(config),
		activeOrders:  make([]bitmex.Order, 0, 10),
		positions:     make([]bitmex.Position, 0, 10),
		margins:       make([]bitmex.Margin, 0, 10),
		//subs:             make([]chan<- []byte, 0, 10),
		newWSClient:      make(chan *bitmex.WSClient), // pool needs a guarantee of receiving, hence unbuffered
		partialsOrder:    false,
		partialsPosition: false,
		partialsMargin:   false,
		partialsCond:     sync.NewCond(&sync.Mutex{}),
		pool:             p,
		pubSub:           NewPubSub(),
	}
	done := make(chan struct{})
	acc.Done = done
	//acc := Account{
	//	activeOrders:  make([]bitmex.Order, 0, 10),
	//	positions:     make([]bitmex.Position, 0, 10),
	//	margins:       make([]bitmex.Margin, 0, 10),
	//	subs:          make([]chan<- []byte, 0, 10),
	//	Done:          make(chan struct{}),
	//	startWSClient: make(chan struct{}, 1),
	//	restClient:    p.bitmex.NewRestClient(config),
	//	pool:          p,
	//	id:            id,
	//}
	//acc.partialsCond = sync.NewCond(&sync.Mutex{})
	//acc.config = config
	//acc.topic = topic

	// first signal to start wsClient
	acc.startWSClient <- struct{}{}

	go acc.manager(done, logger)
	return &acc, nil
}

func (a *Account) authenticate(wsClient *bitmex.WSClient, logger *log.Logger) {
	select {
	case <-wsClient.Done:
		return
	//case <-ctx.Done():
	//	a.close(ctx.Err().Error(), logger)
	//	return
	case <-a.Done:
		return
	default:
	}

	err := wsClient.Authenticate(context.TODO())
	if err == nil {
		fmt.Println("Authenticated: ", a.id)
		return
	}

	logger.Println(errors.Wrap(err, "authentication error"))

	if errors.Cause(err) == bitmex.ErrInvalidAPIKey {
		//fmt.Println("attempting close: ", a.id)
		a.close("authentication failed")
		//fmt.Println("attempted close: ", a.id)
		return
	}

	if errors.Cause(err) == bitmex.ErrRequestExpired {
		fmt.Println("Your system's time is not properly synchronized.")
	}

	// time to retry if fails due to Unknown/retryable error
	retry := time.After(time.Second)

	// return if context is canceled or the client is destroyed by other means
	select {
	//case <-ctx.Done():
	//	a.close(ctx.Err().Error(), logger)
	//	return
	case <-a.Done:
		return
	case <-retry:
	}

	// other possible errors
	logger.Println("retrying auth for", a.id)
	// ErrWSVerificationTimeout, ErrServerError, ErrClientError, ErrUnexpectedError, ErrRequestExpired
	a.authenticate(wsClient, logger)
}

func (a *Account) WaitForPartial() {
	a.partialsCond.L.Lock()
	for !a.partialsOrder || !a.partialsPosition || !a.partialsMargin {
		a.partialsCond.Wait()
	}
	a.partialsCond.L.Unlock()
}

func (a *Account) subscribeAllTables(wsClient *bitmex.WSClient, logger *log.Logger) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		a.subscribeTable(wsClient, bitmex.WSTableOrder, logger)
	}()
	go func() {
		defer wg.Done()
		a.subscribeTable(wsClient, bitmex.WSTablePosition, logger)
	}()
	go func() {
		defer wg.Done()
		a.subscribeTable(wsClient, bitmex.WSTableMargin, logger)
	}()
	wg.Wait()
}

func (a *Account) subscribeTable(wsClient *bitmex.WSClient, table string, logger *log.Logger) {

	select {
	case <-wsClient.Done:
		//fmt.Println("ws client done", a.id)
		return
	//case <-ctx.Done():
	//	a.close(ctx.Err().Error(), logger)
	//	return
	case <-a.Done:
		//fmt.Println(fmt.Sprintf("Account done top: %s", a.id))
		return
	default:
	}

	err := wsClient.SubscribeTableWithPartials(context.TODO(), table)
	if err == nil {
		fmt.Printf("subscribe success for %s : %s\n", table, a.id)
		return
	}

	logger.Println(err, a.id)

	if errors.Cause(err) == bitmex.ErrInvalidAPIKey {
		a.close("subscribeTable")
		return
	}

	if errors.Cause(err) == bitmex.ErrBadRequest || errors.Cause(err) == bitmex.ErrClientError {
		fmt.Println("Unsubscribing")
		_ = wsClient.Unsubscribe(context.TODO())
	}

	if errors.Cause(err) == bitmex.ErrRequestExpired {
		fmt.Println("Your system's time is not properly synchronized.")
	}

	// time to retry if fails due to Unknown/retryable error
	retry := time.After(time.Second)

	// return if context is canceled or the client is destroyed by other means
	select {
	//case <-ctx.Done():
	//	fmt.Println(fmt.Sprintf("context canceled for the Account: %s", a.id))
	//	a.close(errors.Wrap(ctx.Err(), fmt.Sprintf("context canceled for the Account: %s", a.id)).Error(), logger)
	//	return
	case <-a.Done:
		fmt.Println(fmt.Sprintf("account Done: %s", a.id))
		return
	case <-retry:
	}

	// other possible errors
	// ErrWSVerificationTimeout, ErrServerError, ErrClientError, ErrUnexpectedError, ErrRequestExpired
	fmt.Println("retrying for: ", a.id)
	a.subscribeTable(wsClient, table, logger)
}

// OrdersCopy returns a deep copy of the active orders.
func (a *Account) OrdersCopy() []bitmex.Order {
	a.ordersMu.RLock()
	defer a.ordersMu.RUnlock()
	orders := make([]bitmex.Order, len(a.activeOrders))
	copy(orders, a.activeOrders)
	return orders
}

// PositionsCopy returns a deep copy of the open positions.
func (a *Account) PositionsCopy() []bitmex.Position {
	a.positionsMu.RLock()
	defer a.positionsMu.RUnlock()
	positions := make([]bitmex.Position, len(a.positions))
	copy(positions, a.positions)
	return positions
}

func (a *Account) MarginsCopy() []bitmex.Margin {
	a.marginsMu.RLock()
	defer a.marginsMu.RUnlock()
	margins := make([]bitmex.Margin, len(a.margins))
	copy(margins, a.margins)
	return margins
}

// startStreaming creates a new WSClient on the common websocket connection,
// it also authenticates the sub connection and subscribes to the required tables.
// this is also responsible for unlocking partials mutex
//func (a *Account) startStreaming(ctx context.Context, logger *log.Logger) {
//	a.resetPartials()
//
//	a.startMu.Lock()
//	select {
//	case <-a.Done:
//		a.startMu.Unlock()
//		return
//	default:
//	}
//	// The manager can start working after this assignment
//	success := a.pool.newWSClient(ctx, a.config, a.topic, logger, a.id)
//
//	select {
//	case <-ctx.Done():
//		a.startMu.Unlock()
//		return
//	case <-success:
//	}
//
//	// the function returns here so the caller (Account's manager) can continue receiving messages
//	// this is important because the websocket receiver will start receiving data after the connect function is called
//	// Although, any concurrent calls to start function will block until the below routine returns.
//	go func() {
//		defer a.startMu.Unlock()
//		a.wsClient.Connect()
//		a.authenticate(ctx, logger)
//		a.subscribeAllTables(ctx, logger)
//	}()
//}

func (a *Account) manager(done chan struct{}, logger *log.Logger) {
	defer func() {
		close(done)

		// Release Wait Partials
		a.partialsCond.L.Lock()
		a.partialsOrder = true
		a.partialsPosition = true
		a.partialsMargin = true
		a.partialsCond.Broadcast()
		a.partialsCond.L.Unlock()

		//fmt.Println("Account closed : ", a.id)
	}()
	//a.startStreaming(ctx, logger)

	// return if context is canceled, or Account already closed by startStreaming()
	for {
		select {
		case by := <-a.closing:
			logger.Println("account closed by: ", by)
			return

		case <-a.startWSClient:
			a.resetPartials()
			var wsClient *bitmex.WSClient
			receiver := make(chan []byte, 1)

			// send wsClient creation request to Pool
			select {
			case by := <-a.closing:
				logger.Println("account closed by: ", by)
				return
			case a.pool.wsClientReq <- wsClientParams{
				config:   a.config,
				topic:    a.topic,
				id:       a.id,
				receiver: receiver,
				logger:   logger,
			}:
			}

			// wait from client from Pool

			select {
			case by := <-a.closing:
				logger.Println("account closed by: ", by)
				return
			case wsClient = <-a.newWSClient:
				//fmt.Println("new client received")
			}

			//wsClient, err := a.pool.newWSClient(a.config, a.topic, a.id, receiver, logger)
			//if err != nil {
			//	a.close(bitmex.ErrContextCanceled.Error(), logger)
			//	return
			//}
			go a.wsClientWorker(wsClient, receiver, logger)
		}
	}
}

// wsClient managers handles the receiver(wsClient) and returns at the end of the wsClient life cycle.
// it is again restarted by the manager, with a new wsClient object.
// wsClient can be closed and destroyed any number of time due to websocket connection closure or unsubscription of
// wsClient by the host.
// wsClient value on the Account type is reassigned with new wsClient object, after which this routine should be called
// to manage the all operations on the new wsClient object.
func (a *Account) wsClientWorker(wsClient *bitmex.WSClient,
	receiver chan []byte, logger *log.Logger) {

	defer func() {
		//logger.Println("unsubscribingx")
		err := wsClient.Unsubscribe(context.Background())
		if err != nil {
			logger.Println(errors.Wrap(err, fmt.Sprintf("failed to unsubscribe account: %s", a.id)))
			logger.Println("Failed to Unsubscribe account: ", a.id)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		a.authenticate(wsClient, logger)
		a.subscribeAllTables(wsClient, logger)
		wg.Done()
	}()

	for {
		select {
		case <-a.Done:
			wg.Wait()
			return
		case <-wsClient.Done:
			wg.Wait()
			select {
			case <-a.Done:
				return
			default:
				// signaling manager to start a new wsClient
				a.startWSClient <- struct{}{}
				return
			}
		case msg := <-receiver:
			//for i := range a.subs {
			//	a.subs[i] <- msg
			//}
			a.publish(msg)
			a.handleMessage(msg, logger)
		}
	}
}

//func (a *Account) subscribe(ch chan<- []byte) {
//	a.subsMu.Lock()
//	defer a.subsMu.Unlock()
//	a.subs = append(a.subs, ch)
//}

func (a *Account) Close() {
	a.close("outside caller")
}

func (a *Account) close(by string) {
	select {
	case a.closing <- by:
		<-a.Done
	case <-a.Done:
	}

	//a.closeDone.Do(func() {
	//	close(a.Done)
	//	a.pool.internalRemovals <- Removal{
	//		id:     a.id,
	//		Reason: errors.New(reason),
	//	}
	//	a.resetPartials()
	//})
	//logger.Printf("closing Account: %s | reason: %s", a.id, reason)
}

func (a *Account) resetPartials() {
	a.partialsCond.L.Lock()
	a.partialsOrder, a.partialsPosition, a.partialsMargin = false, false, false
	a.partialsCond.L.Unlock()
}
