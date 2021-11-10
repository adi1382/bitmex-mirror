package pool

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/bitmex-mirror/auth"
	"github.com/bitmex-mirror/bitmex"
)

// Account represents a user account.
type account struct {
	startMu            sync.Mutex         // Mutex to protect the startStreaming process
	wsClient           *bitmex.WSClient   // Websocket client
	Done               chan struct{}      // Channel to signal the end of the life of account
	config             auth.Config        // Config for the account
	topic              string             // Topic for the account (additional optional parameter for websocket)
	ordersMu           sync.RWMutex       // Mutex to protect activeOrders
	positionsMu        sync.RWMutex       // Mutex to protect positions
	marginMu           sync.RWMutex       // Mutex to protect margins
	activeOrders       []bitmex.Order     // current state of active orders for the account
	positions          []bitmex.Position  // current state of positions for the account
	margins            []bitmex.Margin    // current state of margins for the account
	restClient         *bitmex.RestClient // REST client for the account for REST requests
	subs               []chan<- []byte    // Channels to send data to other accounts subscribing to this account
	subsMu             sync.RWMutex       // Mutex to protect subs slice
	pool               *Pool              // Pool to which this account belongs (used to create new wsClient)
	isPartialsReceived bool               // Flag to indicate if partials were received
	partialsCond       *sync.Cond         // Conditional variable to signal that partials were received
	closing            chan error
	id                 string //multi-plexed websocket id
}

func (p *Pool) newAccount(ctx context.Context, config auth.Config, topic string, logger *log.Logger, id string) (*account, error) {
	acc := account{
		activeOrders: make([]bitmex.Order, 0, 10),
		positions:    make([]bitmex.Position, 0, 10),
		margins:      make([]bitmex.Margin, 0, 10),
		subs:         make([]chan<- []byte, 0, 10),
		Done:         make(chan struct{}),
		closing:      make(chan error),
		restClient:   p.bitmex.NewRestClient(config),
		pool:         p,
		id:           id,
	}
	acc.partialsCond = sync.NewCond(&sync.Mutex{})
	acc.config = config
	acc.topic = topic

	go acc.manager(ctx, logger)
	return &acc, nil
}

func (a *account) authenticate(ctx context.Context, logger *log.Logger) {
	select {
	case <-a.wsClient.Done:
		return
	case <-ctx.Done():
		a.close(ctx.Err())
		return
	case <-a.Done:
		return
	default:
	}

	err := a.wsClient.Authenticate(ctx)
	if err == nil {
		return
	}

	logger.Println(err)

	if errors.Cause(err) == bitmex.ErrInvalidAPIKey {
		a.close(err)
		return
	}

	if errors.Cause(err) == bitmex.ErrRequestExpired {
		fmt.Println("Your system's time is not properly synchronized.")
	}

	// time to retry if fails due to Unknown/retryable error
	retry := time.After(time.Second)

	// return if context is canceled or the client is destroyed by other means
	select {
	case <-ctx.Done():
		a.close(ctx.Err())
		return
	case <-a.Done:
		return
	case <-retry:
	}

	// other possible errors
	// ErrWSVerificationTimeout, ErrServerError, ErrClientError, ErrUnexpectedError, ErrRequestExpired
	a.authenticate(ctx, logger)
}

func (a *account) WaitForPartial() {
	a.partialsCond.L.Lock()
	for !a.isPartialsReceived {
		a.partialsCond.Wait()
	}
	a.partialsCond.L.Unlock()
}

func (a *account) subscribeStreams(ctx context.Context, logger *log.Logger) {

	//fmt.Println("start for, ", a.id)
	//defer fmt.Println("done for ", a.id)

	select {
	case <-a.wsClient.Done:
		fmt.Println("ws client done", a.id)
		return
	case <-ctx.Done():
		a.close(ctx.Err())
		return
	case <-a.Done:
		fmt.Println(fmt.Sprintf("account done top: %s", a.id))
		return
	default:
	}

	err := a.wsClient.SubscribeTables(ctx, bitmex.WSTableOrder, bitmex.WSTablePosition, bitmex.WSTableMargin)
	if err == nil {
		//fmt.Println("success for: ", a.id)
		return
	}

	logger.Println(err, a.id)

	if errors.Cause(err) == bitmex.ErrInvalidAPIKey {
		a.close(err)
		return
	}

	if errors.Cause(err) == bitmex.ErrRequestExpired {
		fmt.Println("Your system's time is not properly synchronized.")
	}

	// time to retry if fails due to Unknown/retryable error
	retry := time.After(time.Second)

	// return if context is canceled or the client is destroyed by other means
	select {
	case <-ctx.Done():
		fmt.Println(fmt.Sprintf("context canceled for the account: %s", a.id))
		a.close(errors.Wrap(ctx.Err(), fmt.Sprintf("context canceled for the account: %s", a.id)))
		return
	case <-a.Done:
		fmt.Println(fmt.Sprintf("account Done: %s", a.id))
		return
	case <-retry:
	}

	// other possible errors
	// ErrWSVerificationTimeout, ErrServerError, ErrClientError, ErrUnexpectedError, ErrRequestExpired
	fmt.Println("retry for: ", a.id)
	a.subscribeStreams(ctx, logger)
}

// OrdersCopy returns a deep copy of the active orders.
func (a *account) OrdersCopy() []bitmex.Order {
	a.ordersMu.RLock()
	defer a.ordersMu.RUnlock()
	orders := make([]bitmex.Order, len(a.activeOrders))
	copy(orders, a.activeOrders)
	return orders
}

// PositionsCopy returns a deep copy of the open positions.
func (a *account) PositionsCopy() []bitmex.Position {
	a.positionsMu.RLock()
	defer a.positionsMu.RUnlock()
	positions := make([]bitmex.Position, len(a.positions))
	copy(positions, a.positions)
	return positions
}

// startStreaming creates a new WSClient on the common websocket connection,
// it also authenticates the sub connection and subscribes to the required tables.
// this is also responsible for unlocking partials mutex
func (a *account) startStreaming(ctx context.Context, logger *log.Logger) {
	a.startMu.Lock()
	select {
	case <-a.Done:
		a.startMu.Unlock()
		return
	default:
	}
	// The manager can start working after this assignment
	a.wsClient = a.pool.newWSClient(ctx, a.config, a.topic, logger, a.id)

	// the function returns here so the caller (account's manager) can continue receiving messages
	// this is important because the websocket receiver will start receiving data after the connect function is called
	// Although, any concurrent calls to start function will block until the below routine returns.
	go func() {
		defer a.startMu.Unlock()
		a.wsClient.Connect()
		a.authenticate(ctx, logger)
		a.subscribeStreams(ctx, logger)

		select {
		case <-a.Done:
			fmt.Println("partials not received: ", a.id)
			return
		case <-a.wsClient.Done:
			fmt.Println("partials not received: ", a.id)
			return
		default:
		}
		a.partialsCond.L.Lock()
		a.isPartialsReceived = true
		a.partialsCond.Broadcast()
		logger.Println("Partials received: ", a.id)
		a.partialsCond.L.Unlock()
	}()
}

func (a *account) manager(ctx context.Context, logger *log.Logger) {
	a.startStreaming(ctx, logger)

	exit := func(Reason error) {
		a.pool.internalRemovals <- Removal{
			APIKey: a.config.Key,
			Reason: Reason,
		}
		close(a.Done)
	}

	for {
		select {
		case msg := <-a.wsClient.Receiver:
			for i := range a.subs {
				a.subs[i] <- msg
			}
			a.handleMessage(msg, logger)
		case <-a.wsClient.Done:
			a.startStreaming(ctx, logger)
			continue
		case <-ctx.Done():
			//a.close(errors.Wrap(ctx.Err(), fmt.Sprintf("context canceled for the account: %s", a.config.Key)))
			exit(errors.Wrap(ctx.Err(), fmt.Sprintf("context canceled for the account: %s", a.config.Key)))
			return
		case reason := <-a.closing:
			exit(reason)
			return
			//case <-a.Done:
			//	return
		}
	}
}

func (a *account) subscribe(ch chan<- []byte) {
	a.subsMu.Lock()
	defer a.subsMu.Unlock()
	a.subs = append(a.subs, ch)
}

func (a *account) close(Reason error) {
	fmt.Println(fmt.Sprintf("%s: account closed: ", a.id), Reason)
	select {
	case a.closing <- Reason:
		<-a.Done
	case <-a.Done:
	}
}
