package bitmex

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bitmex-mirror/auth"
	"github.com/pkg/errors"

	"golang.org/x/time/rate"
)

// NewWSClient subscribes a new stream for the client over the multiplexed websocket connection
// it returns a variable of type *WSClient which implements methods to subscribe, unsubscribe, authenticate and write
// standard message over the subscribed stream.
// These methods provide a rest-like interface which returns errors received over the subscribed stream.
func (ws *WSConnection) NewWSClient(ctx context.Context, config auth.Config,
	topic string, logger *log.Logger, id string) (*WSClient, error) {

	c := WSClient{
		config:     config,
		topic:      topic,
		connWriter: ws.chWrite,
		bucketM:    ws.bucketM,
		id:         id,
	}

	receiver := make(chan []byte, 10000)

	c.Done = make(chan struct{})
	c.closing = make(chan string)

	//c.subsData = make(chan []byte, 10)
	//c.unSubsData = make(chan []byte, 10)
	//c.authData = make(chan []byte, 10)
	//c.cancelAfterData = make(chan []byte, 10)
	c.socketMessage = make(chan []interface{}, 10)
	c.bucketM = ws.bucketM

	select {
	case <-ws.Done:
		return nil, ErrWSConnClosed
	default:
	}
	c.SocketDone = ws.Done

	go c.manager(ctx, logger, receiver, ws.removeClient)
	c.Receiver = receiver

	ws.clLock.Lock()
	ws.clients[id] = &c
	ws.clLock.Unlock()

	return &c, nil
}

type WSClient struct {
	// signal channel are closed by the respective functions to receive data on data channels
	// signal channel are then set to nil to stop receiving data on data channels
	// This implementation is required so the exported methods that writes over websocket are able to confirm
	// if their request was successful or not
	//subsSignal        chan struct{}
	//subsData          chan []byte
	//unSubsSignal      chan struct{}
	//unSubsData        chan []byte
	//authSignal        chan struct{}
	//authData          chan []byte
	//cancelAfterSignal chan struct{}
	//cancelAfterData   chan []byte

	id              string
	subscriptions   []pubSub
	subscriptionsMu sync.RWMutex

	//subsMu   sync.Mutex // Mutex for SubscribeTables and UnsubscribeTables
	//authMu   sync.Mutex // Mutex for Authenticate
	//cancelMu sync.Mutex // Mutex for CancelAllAfter

	Receiver      <-chan []byte      // returned to NewWSClient caller, closes when client drops
	Done          chan struct{}      // Done is closed when client is closed
	SocketDone    chan struct{}      // SocketDone is closed when ws connection is closed
	closing       chan string        // signals manager to close the done channel
	socketMessage chan []interface{} // socketMessage is used by the router to send message from socket for the client

	bucketM    *rate.Limiter // rate limit bucket shared with rest API calls
	connWriter chan<- []byte // websocket writer channel
	config     auth.Config   // config contains the key and method to generate signature
	topic      string        // topic for multiplexed socket stream
}

// Connect subscribes to the multiplexed socket stream, this function must be called before calling any other method
// Expect to receive data Receiver channel after calling this method, so have the read on the Receiver already ready
// before calling this function.
// This function should only be called once.
func (c *WSClient) Connect() {
	msg := make([]interface{}, 0, 4)
	msg = append(msg, 1, c.id, c.topic)
	msgBytes, _ := json.Marshal(msg)

	select {
	case c.connWriter <- msgBytes:
	case <-c.Done:
	case <-c.SocketDone:
	}
}

func (c *WSClient) VerifyPartials(ctx context.Context, tables ...string) <-chan bool {
	ClientSubs := make(chan []byte, 1)
	ctx2, cancel := context.WithCancel(ctx)
	c.subscribe(ctx2, ClientSubs)

	type partials struct {
		Table  string            `json:"table"`
		Action string            `json:"action"`
		Data   []json.RawMessage `json:"data"`
	}

	returnCh := make(chan bool, 1)

	go func() {
		defer cancel()
		for {
			select {
			case msg := <-ClientSubs:
				//fmt.Println("partial message: ", c.id, string(msg))
				var partialMessage partials
				if err := json.Unmarshal(msg, &partialMessage); err != nil {
					break
				}
				if partialMessage.Action == WSDataActionPartial {
					for i := range tables {
						if partialMessage.Table == tables[i] {
							tables[i] = tables[len(tables)-1]
							tables = tables[:len(tables)-1]
							break
						}
					}
					if len(tables) == 0 {
						returnCh <- true
						return
					}
				}
			case <-ctx2.Done():
				returnCh <- false
				return
			}
		}
	}()

	return returnCh
}

// SubscribeTables is used to subscribe different tables like orders, position, etc.
// When subscribing a private table, you must ensure that the client stream is successfully authenticated
// using the Authenticate function.
// This function does not return any error other than 429 rate-limiting error.
// You must ensure that you subscribe to tables that exist, and subscribe private tables only after
// successful authentication. If the subscription fails due to any error other than 429 error, the function will just
// return a nil error value, however the error message can be caught through receivers channel.
// update 429 is also handled by retry mechanism, so receiving error is very rare.
func (c *WSClient) SubscribeTables(ctx context.Context, tables ...string) error {
	// return if client is already closed
	select {
	case <-c.Done:
		select {
		case <-c.SocketDone:
			return ErrWSConnClosed
		default:
			return ErrWSClientClosed
		}
	default:
	}

	// preparing message to send multiplexed socket connection
	message := make([]interface{}, 0, 4)
	payload := wsMessage{
		Op: "subscribe",
	}
	payload.Args = tables
	message = append(message, 0, c.id, c.topic, payload)
	msgByte, _ := json.Marshal(message)

	// this context is used to taking tokens from bucket
	// cancel func cancels when client is closed
	//fmt.Println(ctx.Deadline())
	//fmt.Println("Deadline: ",)
	ctxWait, cancelWait := context.WithCancel(ctx)
	defer cancelWait()

	go func() {
		select {
		case <-ctxWait.Done():
		case <-c.Done:
			cancelWait()
		}
	}()

	// waiting for tokens
	// number of tokens = number of tables subscribing
	//fmt.Println("Allow: ", c.bucketM.AllowN(time.Now(), len(tables)))
	if err := c.bucketM.WaitN(ctxWait, len(tables)); err != nil {
		return err
	}
	//fmt.Println("released", c.id, time.Now())

	// setting a timeout for confirmation from socket connection, this should not take long
	ctxConf, cancelConf := context.WithTimeout(ctxWait, wsConfirmTimeout)
	defer cancelConf()

	table2 := make([]string, len(tables))
	copy(table2, tables)
	partials := c.VerifyPartials(ctxConf, table2...)

	// sending message to socket writer or returning if client is already closed

	ctxSubs, subscriptionCancel := context.WithCancel(ctxConf)
	defer subscriptionCancel()

	ClientSubs := make(chan []byte, 1)
	c.subscribe(ctxSubs, ClientSubs)

	select {
	case <-c.Done:
		select {
		case <-c.SocketDone:
			return ErrWSConnClosed
		default:
			return ErrWSClientClosed
		}
	default:
		select {
		case <-c.Done:
			select {
			case <-c.SocketDone:
				return ErrWSConnClosed
			default:
				return ErrWSClientClosed
			}
		case c.connWriter <- msgByte:
		}
	}
	//fmt.Println("sent", c.id, time.Now())

	n := len(tables)
	success := make(chan struct{}, 3) // this channel is closed when the message is confirmed
	g, ctxErr := errgroup.WithContext(ctxConf)

	var retErr error

L:
	for {
		select {
		case msg := <-ClientSubs:
			fmt.Println("Verification Message: ", c.id, string(msg))
			g.Go(func() error {
				if request, ok := requestField(msg); ok && validateSubsReq(request, payload) {
					if isSuccessful(msg) {
						select {
						case success <- struct{}{}:
						}
						return nil
					}
					return c.wsError(msg)
				}
				return nil
			})
		case <-ctxErr.Done():
			retErr = ErrWSVerificationTimeout
			break L
		case <-success:
			n--
			if n == 0 {
				break L
			}
		}
	}

	// unsubscribe from stream
	subscriptionCancel()

	//fmt.Println("verified", c.id, time.Now())

	if err := g.Wait(); err != nil {
		cause := errors.Cause(err)
		//if cause == ErrTooManyRequests || cause == ErrServerOverloaded {
		//	time.Sleep(time.Millisecond * 500)
		//	return c.SubscribeTables(ctx, tables...)
		//}
		retErr = err
		if cause == ErrBadRequest {
			// don't return a bad request error
			// no need to wait for partials
			return nil
		}
	}

	if retErr != nil {
		return retErr
	}

	select {
	case t := <-partials:
		//fmt.Println("partials", t, c.id, time.Now())
		if t {
			return nil
		} else {
			c.UnsubscribeConnection()
			return ErrWSConnClosed
		}
	}
}

// Authenticate sends an authentication message over the websocket connection for the client.
// If it returns a non-nil error then most likely connection is now unsubscribed for the client, and you will
// need to use the NewWSClient function again
func (c *WSClient) Authenticate(ctx context.Context) error {
	// This API is implemented in the very similar way as that of SubscribeTables
	// Authentication through websocket does consume a token from request limiter
	// It does not need num variable as there is only one expected message from websocket

	select {
	case <-c.Done:
		select {
		case <-c.SocketDone:
			return ErrWSConnClosed
		default:
			return ErrWSClientClosed
		}
	default:
	}

	// preparing message to send multiplexed socket connection
	///////////////////////////////////////////////////////////////////////////////////
	msg := make([]interface{}, 0, 4)
	apiExpires := time.Now().Add(requestTimeout).Unix()
	signature := c.config.Sign("GET/realtime" + strconv.FormatInt(apiExpires, 10))
	payload := wsMessage{
		Op: "authKeyExpires",
	}
	payload.Args = []interface{}{c.config.Key, apiExpires, signature}
	msg = append(msg, 0, c.id, c.topic, payload)
	msgByte, _ := json.Marshal(msg)
	///////////////////////////////////////////////////////////////////////////////////

	ctxConf, cancelConf := context.WithTimeout(ctx, wsConfirmTimeout)
	defer cancelConf()

	ClientSubs := make(chan []byte, 1)
	c.subscribe(ctxConf, ClientSubs)

	select {
	case <-c.Done:
		select {
		case <-c.SocketDone:
			return ErrWSConnClosed
		default:
			return ErrWSClientClosed
		}
	default:
		select {
		case <-c.Done:
			select {
			case <-c.SocketDone:
				return ErrWSConnClosed
			default:
				return ErrWSClientClosed
			}
		case c.connWriter <- msgByte:
		}
	}

	success := make(chan struct{}, 1)
	g, ctxErr := errgroup.WithContext(ctx)

	var errOld error

L:
	for {
		select {
		case msg := <-ClientSubs:
			g.Go(func() error {
				if request, ok := requestField(msg); ok && validateAuthReq(request, payload) {
					if isSuccessful(msg) {
						select {
						case success <- struct{}{}:
							return nil
						case <-ctx.Done():
							return nil
						}
					}
					return c.wsError(msg)
				}
				return nil
			})
		case <-ctxErr.Done():
			errOld = ErrWSVerificationTimeout
			break L
		case <-success:
			break L
		}
	}

	if err := g.Wait(); err != nil {
		cause := errors.Cause(err)
		if cause == ErrTooManyRequests || cause == ErrServerOverloaded {
			time.Sleep(time.Millisecond * 500)
			return c.Authenticate(ctx)
		}
		if cause == ErrBadRequest {
			// don't return a bad request error
			return nil
		}
		return err
	}

	return errOld
}

// UnsubscribeConnection gracefully closes the client without closing the websocket connection, it notifies the host
// that no data is required on this client and close the receiver channel.
// this function should not be called on another routine, and it should be allowed to block until the client is closed
// don't attempt to restart the client until this function returns.
func (c *WSClient) UnsubscribeConnection() {
	select {
	case <-c.Done:
		return
	default:
	}

	msg := make([]interface{}, 0, 3)
	msg = append(msg, 2, c.id, c.topic)
	msgByte, _ := json.Marshal(msg)

	select {
	case <-c.Done:
		return
	default:
		select {
		case <-c.Done:
			return
		case c.connWriter <- msgByte:
		}
	}

	// Once called, the UnsubscribeConnection function will take over the socketMessage channel
	// this might be shared with extractPayload function, in either case
	// if Unsubscribe message is detected, stop function will be triggered.

	// If UnsubscribeConnection is called from outside caller, extractPayload function will continue to work
	// but if UnsubscribeConnection is called from manager, extractPayload will not be reachable, as
	// UnsubscribeConnection will block the manager.

	// This function shall only return after ensuring that the host has sent Unsubscribe message
	t := time.After(wsConfirmTimeout)

	for {
		select {
		case <-c.SocketDone:
			return
		case <-c.Done:
			return
		case <-t:
			// retry Unsubscribe as client is still alive and no unsubscribe message has been received
			c.UnsubscribeConnection()
		case socketMessage := <-c.socketMessage:
			if len(socketMessage) <= 3 {
				if v, ok := msg[0].(float64); ok {
					if v == 2 {
						go c.stop("Unsubscribe function")
						return
					}
				}
			}
		}
	}
}

// stop function can be called by anyone to close the client and all its processes
// it can be called concurrently any number of times
func (c *WSClient) stop(by string) {
	select {
	case c.closing <- by:
		<-c.Done
	case <-c.Done:
	}
}

// manager handles several operations for the WSClient variable.
// It starts at the construction of client, and it returns to close the client completely and close Receiver
// and close done channel to notify all internal operation about the closure.
// It receives on socketDone channel which closes when websocket connection is closed.
// It receives on socketMessage channel which receives data from router, manager sends this data to extractPayload.
// It receives on sendCh which receives payload from extractPayload, manager sends the data to receivers channel.
// It also receives on client's context and closing channel which receives when stop function is caller to close the
// client.
func (c *WSClient) manager(ctx context.Context, logger *log.Logger, receiver chan<- []byte, removeClient chan<- string) {

	defer func() {
		for len(c.socketMessage) > 0 {
			if v, ok := c.extractPayload(<-c.socketMessage, logger); ok {
				select {
				case receiver <- v:
					go c.publish(v)
					//c.toWriterFunc(v)
				default:
				}
			}
		}
		close(c.Done)
		//close(c.Receiver)

		select {
		case removeClient <- c.id:
			fmt.Println("client removed")
		case <-c.SocketDone:
		}
	}()

	for {
		select {
		// canceled by context from the one who created this wsClient
		case <-ctx.Done():
			logger.Printf("websocket: info: client: %s socket channel was done because: %s", c.config.Key,
				"client context done")
			c.UnsubscribeConnection()
			return
		// ws connection is closed
		case <-c.SocketDone:
			logger.Printf("websocket: info: client: %s socket channel was done because: %s",
				c.config.Key, "socket connection dropped")
			return
		// stop function is called
		case stoppedBy := <-c.closing:
			logger.Printf("websocket: info: client: %s socket channel was done by: %s",
				c.config.Key, stoppedBy)
			return
		// data from router
		case msg := <-c.socketMessage:
			if v, ok := c.extractPayload(msg, logger); ok {
				select {
				case receiver <- v:
					c.publish(v)
					//c.toWriterFunc(v)
				default:
					c.UnsubscribeConnection()
					logger.Printf("websocket: info: client: %s socket channel was done because: %s",
						c.config.Key, "receive channel is jacked up to the tits")
					return
				}
			}
		}
	}
}

// isSuccessful checks for the 'success' key and return its value
// if the field is not found, the map lookup automatically return the false default value.
func isSuccessful(data []byte) bool {
	var res map[string]bool
	_ = json.Unmarshal(data, &res)
	return res["success"]
}

// requestField scans for the field 'request' if found it returns the request data on the incoming socket message.
// if not found it returns nil with a false for not found.
func requestField(data []byte) ([]byte, bool) {
	socketRes := make(map[string]interface{})
	if err := json.Unmarshal(data, &socketRes); err != nil {
		return nil, false
	}

	if request, ok := socketRes["request"]; ok {
		marshaledReq, _ := json.Marshal(request)
		return marshaledReq, true
	}
	return nil, false
}

// wsError scans the socket message for error and status field and return an error value of type APIError
// which implements the error interface.
// TODO: request timeout error handling
func (c *WSClient) wsError(data []byte) error {
	errStr := string(data)
	bitmexErr := APIError{Name: "WebsocketError"}
	var res map[string]interface{}

	var err error
	var ok bool

	if err = json.Unmarshal(data, &res); err != nil {
		return errors.Wrap(ErrUnexpectedError, fmt.Sprintf("wsError: unmarshal error Data: %s", errStr))
	}

	if _, ok = res["status"].(float64); !ok {
		return errors.Wrap(ErrUnexpectedError, fmt.Sprintf("wsError: status field not found: %s", errStr))
	}
	bitmexErr.StatusCode = int64(res["status"].(float64))

	if _, ok = res["error"].(string); !ok {
		bitmexErr.Message = fmt.Sprintf("wsError: error field not found: %s", errStr)
	} else {
		bitmexErr.Message = res["error"].(string)
	}

	if bitmexErr.StatusCode == 400 {
		// catch request expired error
		if strings.Contains(bitmexErr.Message, "request") && strings.Contains(bitmexErr.Message, "expired") {
			return errors.Wrap(ErrRequestExpired, bitmexErr.Message)
		}

		return errors.Wrap(ErrBadRequest, bitmexErr.Error())
	}

	if bitmexErr.StatusCode == 401 || bitmexErr.StatusCode == 403 {
		return errors.Wrap(ErrInvalidAPIKey, bitmexErr.Error())
	}

	if bitmexErr.StatusCode == 503 {
		return ErrServerOverloaded
	}

	if bitmexErr.StatusCode == 429 {
		rateLimited := false

		if m, ok := res["meta"]; ok {
			if meta, ok := m.(map[string]interface{}); ok {
				if wait, ok := meta["retryAfter"]; ok {
					if wait, ok := wait.(float64); ok {
						rateLimited = true
						c.RateLimited(time.Second*time.Duration(wait), rateLimitMinuteCap)
					}
				}
			}
		}

		if !rateLimited {
			c.RateLimited(time.Minute, rateLimitMinuteCap)
		}
		return ErrTooManyRequests
	}

	if bitmexErr.StatusCode >= 500 {
		return ErrServerError
	} else if bitmexErr.StatusCode >= 400 {
		return ErrClientError
	}
	// If http code is less than 400 than it's not really an error
	return nil
}

// equalStrSlice checks if all the elements of slice 'a' are present in slice 'b' and vice-versa
// it is allowed sort the slice, but is not allowed to remove/replace/add any elements from either of the slice
func equalStrSlice(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// validateSubsReq verifies if the request field on the incoming socket message is the same as that of
// the request message sent by the functions SubscribeTables and UnsubscribeTables,
// it returns true if the validation is successful.
func validateSubsReq(req []byte, payload wsMessage) bool {
	request := struct {
		Op   string   `json:"op,omitempty"`
		Args []string `json:"args,omitempty"`
	}{}

	if err := json.Unmarshal(req, &request); err == nil {
		if request.Op == payload.Op {
			if strSlice, ok := payload.Args.([]string); ok {
				return equalStrSlice(request.Args, strSlice)
			}
		}
	}
	return false
}

// validateAuthReq verifies if the request field on the incoming socket message is the same as that of
// the request message sent by Authenticate function, it returns true if the validation is successful.
func validateAuthReq(req []byte, payload wsMessage) bool {
	request := struct {
		Op   string        `json:"op,omitempty"`
		Args []interface{} `json:"args,omitempty"`
	}{}

	payloadArgs := payload.Args.([]interface{})

	if err := json.Unmarshal(req, &request); err == nil {
		if request.Op == payload.Op && len(request.Args) == len(payloadArgs) {
			for i := range request.Args {
				switch value := request.Args[i].(type) {
				case float64:
				L1:
					for j := range payloadArgs {
						switch payloadValue := payloadArgs[j].(type) {
						case int64:
							if payloadValue == int64(value) {
								payloadArgs[j] = payloadArgs[len(payloadArgs)-1]
								payloadArgs = payloadArgs[:len(payloadArgs)-1]
								break L1
							}
						}
					}
				case string:
				L2:
					for j := range payloadArgs {
						switch payloadValue := payloadArgs[j].(type) {
						case string:
							if payloadValue == value {
								payloadArgs[j] = payloadArgs[len(payloadArgs)-1]
								payloadArgs = payloadArgs[:len(payloadArgs)-1]
								break L2
							}
						}
					}
				}
			}

			return len(payloadArgs) == 0
		}
	}
	return false
}

// extractPayload processes the message from websocket, extract the payload byte and send it over to sendCh.
// This function should not be called concurrently as an unsubscribeConnection message from websocket can
// provoke this function to call stop() which would prevent the previous messages from ever reaching to sendCh
// It is absolutely critical that the caller receives all messages for the client up until the last message
// of unsubscribe connection from the socket
func (c *WSClient) extractPayload(msg []interface{}, logger *log.Logger) ([]byte, bool) {
	if len(msg) <= 3 {
		if v, ok := msg[0].(float64); ok {
			if v == 2 {
				go c.stop("Websocket Sent a Unsubscribe message on multiplexed connection")
				return nil, false
			}
		}
		return nil, false
	}

	payload, err := json.Marshal(msg[3])
	if err != nil {
		logger.Println(errors.Wrap(errors.Wrap(err, string(payload)),
			"websocket: error: no payload on incoming message"))
		return nil, false
	}

	return payload, true
}

// standard syntax of payload
type wsMessage struct {
	Op   string      `json:"op,omitempty"`
	Args interface{} `json:"args,omitempty"`
}

func (c *WSClient) RateLimited(timeToSleep time.Duration, tokens int) {
	fmt.Println("^^^^^^^^^^^^^^^^^^^^Rate Limit Exceeded^^^^^^^^^^^^^^^^^^^^")
	defer c.bucketM.SetBurst(tokens)

	c.bucketM.SetBurst(0)
	// wait for at least 3 seconds in case of being rate limited
	t := time.NewTimer(timeToSleep + time.Second*3)

	select {
	case <-t.C:
	case <-c.Done:
	}
}

type pubSub struct {
	ch   chan []byte
	done <-chan struct{}
}

func (c *WSClient) subscribe(ctx context.Context, ch chan []byte) {
	c.subscriptionsMu.Lock()
	defer c.subscriptionsMu.Unlock()
	p := pubSub{
		ch:   ch,
		done: ctx.Done(),
	}
	c.subscriptions = append(c.subscriptions, p)
}

func (c *WSClient) publish(msg []byte) {
	c.subscriptionsMu.Lock()
	defer c.subscriptionsMu.Unlock()

	if len(c.subscriptions) == 0 {
		return
	}

	//fmt.Println("Publish: ", c.id, len(c.subscriptions), string(msg))

	//https://stackoverflow.com/questions/20545743/how-to-remove-items-from-a-slice-while-ranging-over-it
	k := 0 // output index
	for i := range c.subscriptions {
		select {
		case <-c.subscriptions[i].done:
		case c.subscriptions[i].ch <- msg:
			c.subscriptions[k] = c.subscriptions[i]
			k++
		}
	}
	c.subscriptions = c.subscriptions[:k]
}

// UnsubscribeTables it is to unsubscribe counterpart of SubscribeTables
// They are implemented in the same way, use this function to unsubscribe already subscribed tables.
//func (c *WSClient) UnsubscribeTables(ctx context.Context, unSubs ...string) error {
//	// This API is implemented the same way as of SubscribeTables
//
//	select {
//	case <-c.Done:
//		select {
//		case <-c.SocketDone:
//			return ErrWSConnClosed
//		default:
//			return ErrWSClientClosed
//		}
//	default:
//	}
//
//	c.subsMu.Lock()
//	defer func() {
//		c.unSubsSignal = nil
//		for len(c.unSubsData) > 0 {
//			<-c.unSubsData
//		}
//		c.subsMu.Unlock()
//	}()
//
//	// preparing message for multiplexed socket connection
//	message := make([]interface{}, 0, 4)
//	payload := wsMessage{
//		Op: "unsubscribe",
//	}
//	payload.Args = unSubs
//	message = append(message, 0, c.config.Key, c.topic, payload)
//	msgByte, _ := json.Marshal(message)
//
//	ctxWait, cancelWait := context.WithCancel(ctx)
//	defer cancelWait()
//
//	go func() {
//		select {
//		case <-ctxWait.Done():
//		case <-c.Done:
//			cancelWait()
//		}
//	}()
//
//	if err := c.bucketM.WaitN(ctxWait, len(unSubs)); err != nil {
//		return err
//	}
//
//	ctxConf, cancelConf := context.WithTimeout(ctxWait, wsConfirmTimeout)
//	defer cancelConf()
//
//	select {
//	case <-c.Done:
//		select {
//		case <-c.SocketDone:
//			return ErrWSConnClosed
//		default:
//			return ErrWSClientClosed
//		}
//	default:
//		select {
//		case <-c.Done:
//			select {
//			case <-c.SocketDone:
//				return ErrWSConnClosed
//			default:
//				return ErrWSClientClosed
//			}
//		case c.connWriter <- msgByte:
//		}
//	}
//
//	c.unSubsSignal = make(chan struct{})
//	close(c.unSubsSignal)
//
//	retChan := make(chan error)
//
//	num := len(unSubs)
//
//	for {
//		select {
//		case res := <-c.unSubsData:
//			go func(ctxConf context.Context, retChan chan<- error, res []byte, payload wsMessage) {
//				if request, ok := requestField(res); ok && validateSubsReq(request, payload) {
//					if isSuccessful(res) {
//						select {
//						case retChan <- nil:
//							return
//						case <-ctxConf.Done():
//							return
//						}
//					}
//					select {
//					case retChan <- c.wsError(res):
//						return
//					case <-ctxConf.Done():
//						return
//					}
//				}
//			}(ctxConf, retChan, res, payload)
//		case err := <-retChan:
//			num--
//			if err != nil {
//				// retry if the error is 429 or 503
//				if err2 := errors.Cause(err); err2 == ErrTooManyRequests || err2 == ErrServerOverloaded {
//					time.Sleep(time.Millisecond * 500)
//					return c.UnsubscribeTables(ctx, unSubs...)
//				}
//
//				if errors.Cause(err) == ErrBadRequest {
//					return nil
//				}
//
//				return err
//			}
//			if num == 0 {
//				return nil
//			}
//		case <-ctxConf.Done():
//			return ErrWSVerificationTimeout
//		case <-c.Done:
//			select {
//			case <-c.SocketDone:
//				return ErrWSConnClosed
//			default:
//				return ErrWSClientClosed
//			}
//		}
//	}
//}
//// CancelAllAfter implements the cancelAllAfter subscription message of the websocket
//func (c *WSClient) CancelAllAfter(ctx context.Context, timeout time.Duration) error {
//	// This API is implemented in the very similar way as that of SubscribeTables
//	// This message consumes one token from the request limiter
//	// It does not need num variable as there is only one expected message from websocket
//
//	select {
//	case <-c.Done:
//		return ErrWSClientClosed
//	default:
//	}
//
//	c.cancelMu.Lock()
//	defer func() {
//		c.cancelAfterSignal = nil
//		for len(c.cancelAfterData) > 0 {
//			<-c.cancelAfterData
//		}
//		c.cancelMu.Unlock()
//	}()
//
//	msg := make([]interface{}, 0, 4)
//	payload := wsMessage{
//		Op: "cancelAllAfter",
//	}
//	payload.Args = timeout.Milliseconds()
//	msg = append(msg, 0, c.config.Key, c.topic, payload)
//	msgByte, _ := json.Marshal(msg)
//
//	ctxWait, cancelWait := context.WithCancel(ctx)
//	defer cancelWait()
//
//	go func() {
//		select {
//		case <-ctxWait.Done():
//		case <-c.Done:
//			cancelWait()
//		}
//	}()
//
//	if err := c.bucketM.Wait(ctxWait); err != nil {
//		return err
//	}
//
//	ctxConf, cancelConf := context.WithTimeout(ctxWait, wsConfirmTimeout)
//	defer cancelConf()
//
//	select {
//	case <-c.Done:
//		return ErrWSClientClosed
//	default:
//		select {
//		case <-c.Done:
//			return ErrWSClientClosed
//		case c.connWriter <- msgByte:
//		}
//	}
//
//	c.cancelAfterSignal = make(chan struct{})
//	close(c.cancelAfterSignal)
//
//	retChan := make(chan error)
//
//	// custom successful function for cancelAllAfter
//	isSuccess := func(data []byte) bool {
//		var res map[string]string
//		_ = json.Unmarshal(data, &res)
//		_, ok := res["cancelTime"]
//		return ok
//	}
//
//	for {
//		select {
//		case res := <-c.cancelAfterData:
//			go func(ctxConf context.Context, retChan chan<- error, res []byte, payload wsMessage) {
//				if request, ok := requestField(res); ok {
//					if validateCancelAllAfterReq(request, payload) {
//						if isSuccess(res) {
//							select {
//							case retChan <- nil:
//								return
//							case <-ctxConf.Done():
//								return
//							}
//						}
//						select {
//						case retChan <- c.wsError(res):
//							return
//						case <-ctxConf.Done():
//							return
//						}
//					}
//				}
//			}(ctxConf, retChan, res, payload)
//		case err := <-retChan:
//			return err
//		case <-ctxConf.Done():
//			return ErrWSVerificationTimeout
//		case <-c.Done:
//			return ErrWSClientClosed
//		}
//	}
//
//}
//
////func (c *WSClient) AutoCancelAllAfter(ctx context.Context, timeout time.Duration, every time.Duration) {
////	burstOffset := int64(math.Ceil(1 / every.Minutes()))
////	m := make(chan int64)
////	go func() {
////		defer close(m)
////		for {
////			select {
////			case m <- burstOffset:
////			case <-c.done:
////			}
////		}
////	}()
////}
//// validateCancelAllAfterReq verifies if the request field on the incoming socket message is the same as that of
//// the request message sent by cancelAllAfter function, it returns true if the validation is successful.
//func validateCancelAllAfterReq(req []byte, payload wsMessage) bool {
//	request := struct {
//		Op   string `json:"op,omitempty"`
//		Args int64  `json:"args,omitempty"`
//	}{}
//
//	if err := json.Unmarshal(req, &request); err == nil {
//		if request.Op == payload.Op {
//			if request.Args == payload.Args.(int64) {
//				return true
//			}
//		}
//	}
//
//	return false
//}
// toWriterFunc sends the processed payload data to the one of the socket writer functions SubscribeTables,
// UnsubscribeTables, Authenticate, CancelAllAfter, these messages are sent so the writer function can verify if
// their request was successful or not signal channels are closed by the writer functions when the data
// sending starts
//func (c *WSClient) toWriterFunc(data []byte) {
//	select {
//	case <-c.subsSignal:
//		c.subsData <- data
//	default:
//	}
//
//	select {
//	case <-c.unSubsSignal:
//		c.unSubsData <- data
//	default:
//	}
//
//	select {
//	case <-c.authSignal:
//		c.authData <- data
//	default:
//	}
//
//	select {
//	case <-c.cancelAfterSignal:
//		c.cancelAfterData <- data
//	default:
//	}
//}
