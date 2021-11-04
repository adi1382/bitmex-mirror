package bitmex

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/bitmex-mirror/auth"
)

// standard syntax of payload
type wsMessage struct {
	Op   string      `json:"op,omitempty"`
	Args interface{} `json:"args,omitempty"`
}

type wsClient struct {
	// signal channel are closed by the respective functions to receive data on data channels
	// signal channel are then set to nil to stop receiving data on data channels
	// This implementation is required so the exported methods that writes over websocket are able to confirm
	// if their request was successful or not
	subsSignal        chan struct{}
	subsData          chan []byte
	unSubsSignal      chan struct{}
	unSubsData        chan []byte
	authSignal        chan struct{}
	authData          chan []byte
	cancelAfterSignal chan struct{}
	cancelAfterData   chan []byte

	subsMu   sync.Mutex // Mutex for SubscribeStreams and UnsubscribeStreams
	authMu   sync.Mutex // Mutex for Authenticate
	cancelMu sync.Mutex // Mutex for CancelAllAfter

	returnCh      chan []byte        // returned to AddWSClient caller, closes when client drops
	done          chan struct{}      // done is closed when client is closed
	closing       chan string        // signals manager to close the done channel
	socketMessage chan []interface{} // socketMessage is used by the router to send message from socket for the client

	bucketM    *rate.Limiter // rate limit bucket shared with rest API calls
	connWriter chan<- []byte // websocket writer channel
	config     auth.Config   // config contains the key and method to generate signature
	topic      string        // topic for multiplexed socket stream
}

// stop function can be called by anyone to close the client and all its processes
// it can be called concurrently any number of times
func (c *wsClient) stop(by string) {
	select {
	case c.closing <- by:
		<-c.done
	case <-c.done:
	}
}

// manager handles several operations for the wsClient variable.
// It starts at the construction of client, and it returns to close the client completely and close returnCh
// and close done channel to notify all internal operation about the closure.
// It receives on socketDone channel which closes when websocket connection is closed.
// It receives on socketMessage channel which receives data from router, manager sends this data to handleSocketMsg.
// It receives on sendCh which receives payload from handleSocketMsg, manager sends the data to receivers channel.
// It also receives on client's context and closing channel which receives when stop function is caller to close the
// client.
func (c *wsClient) manager(ctx context.Context, logger *log.Logger,
	socketDone <-chan struct{}, removeClient chan<- string) {

	defer func() {
		for len(c.socketMessage) > 0 {
			if v, ok := c.handleSocketMsg(<-c.socketMessage, logger); ok {
				select {
				case c.returnCh <- v:
					c.toWriterFunc(v)
				default:
				}
			}
		}
		close(c.done)
		close(c.returnCh)

		select {
		case removeClient <- c.config.Key:
		case <-socketDone:
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Printf("websocket: info: client: %s socket channel was done because: %s", c.config.Key,
				"client context done")
			go c.UnsubscribeConnection()
			return
		case <-socketDone:
			logger.Printf("websocket: info: client: %s socket channel was done because: %s",
				c.config.Key, "socket connection dropped")
			return
		case stoppedBy := <-c.closing:
			logger.Printf("websocket: info: client: %s socket channel was done by: %s",
				c.config.Key, stoppedBy)
			return
		case msg := <-c.socketMessage:
			if v, ok := c.handleSocketMsg(msg, logger); ok {
				select {
				case c.returnCh <- v:
					c.toWriterFunc(v)
				default:
					go c.UnsubscribeConnection()
					logger.Printf("websocket: info: client: %s socket channel was done because: %s",
						c.config.Key, "receive channel is jacked up to the tits")
					return
				}
			}
		}
	}
}

// SubscribeStreams is used to subscribe different tables like orders, position, etc.
// When subscribing a private table, you must ensure that the client stream is successfully authenticated
// using the Authenticate function.
// This function does not return any error other than 429 rate-limiting error.
// You must ensure that you subscribe to tables that exist, and subscribe private tables only after
// successful authentication. If the subscription fails due to any error other than 429 error, the function will just
// return a nil error value, however the error message can be caught through receivers channel.
func (c *wsClient) SubscribeStreams(ctx context.Context, subs ...string) error {
	// return if client is already closed
	select {
	case <-c.done:
		return errors.New("client closed")
	default:
	}

	// Locking subsMu to prevent concurrent calls to SubscribeStreams or UnsubscribeStreams
	c.subsMu.Lock()

	defer func() {
		// Setting subsSignal to nil so no data is received on subsData from toWriterFunc()
		c.subsSignal = nil
		// flushing subsData
		for len(c.subsData) > 0 {
			<-c.subsData
		}
		c.subsMu.Unlock()
	}()

	// preparing message to send multiplexed socket connection
	message := make([]interface{}, 0, 4)
	payload := wsMessage{
		Op: "subscribe",
	}
	payload.Args = subs
	message = append(message, 0, c.config.Key, c.topic, payload)
	msgByte, _ := json.Marshal(message)

	// this context is used to taking tokens from bucket
	// cancel func cancels when client is closed
	ctxWait, cancelWait := context.WithCancel(ctx)
	defer cancelWait()

	go func() {
		select {
		case <-ctxWait.Done():
		case <-c.done:
			cancelWait()
		}
	}()

	// waiting for tokens
	// number of tokens = number of tables subscribing
	if err := c.bucketM.WaitN(ctxWait, len(subs)); err != nil {
		return err
	}

	// setting a timeout for confirmation from socket connection, this should not take long
	ctxConf, cancelConf := context.WithTimeout(ctxWait, wsConfirmTimeout)
	defer cancelConf()

	// sending message to socket writer or returning if client is already closed
	select {
	case <-c.done:
		return errors.New("client closed")
	default:
		select {
		case <-c.done:
			return errors.New("client closed")
		case c.connWriter <- msgByte:
		}
	}

	// making subsSignal alive for receiving data on subsData
	c.subsSignal = make(chan struct{})
	close(c.subsSignal)

	// retChan will be sent to different goroutines to receive error if any from socket
	// For the case of SubscribeStreams and UnsubscribeStreams functions only 429 error code will be returned
	// all other cases should be handled by the caller of the API using main receiver channel
	retChan := make(chan error)

	// socket will send a success or fail message for each table subscribed
	num := len(subs)

	// receiving data from subsData and checking for any error
	for {
		select {
		case res := <-c.subsData:
			// starting a go routine for each incoming msg
			// If any routine is blocked due to extra messages than anticipated with same request then
			// the routine will return automatically on return of SubscribeStreams because of CtxConf
			// payload is reused and implements slice,
			// so it should be ensured the values are not deleted down the call stack
			// res is explicitly passed to prevent sharing the value
			go func(ctxConf context.Context, retChan chan<- error, res []byte, payload wsMessage) {
				// validating the message if it is for the same request this function sent
				if request, ok := requestField(res); ok && validateSubsReq(request, payload) {
					// checking the request was successful
					if isSuccessful(res) {
						select {
						case retChan <- nil:
							return
						case <-ctxConf.Done():
							return
						}
					}

					// request failed, sending error on retChan
					select {
					case retChan <- wsError(res):
						return
					case <-ctxConf.Done():
						return
					}
				}
			}(ctxConf, retChan, res, payload)
		case err := <-retChan:
			// decrementing counter
			// expected number of messages wih same request = number of tables subscribed
			num--
			if err != nil {
				// returning error only in the case of 429 status code
				// this is done to promote responsible use of the API, and the error message could get complicated
				// as a different error could be received for each table subscribed
				if err.(APIError).StatusCode == 429 {
					return err
				}
			}
			// return if all messages are successful
			if num == 0 {
				return nil
			}
		// if all messages with validated request are not received before canceling of context
		case <-ctxConf.Done():
			return errors.New("context canceled/timeout: message sent but could not verify")
		// Returning if client is closed
		case <-c.done:
			return errors.New("client closed")
		}
	}
}

// UnsubscribeStreams it is to unsubscribe counterpart of SubscribeStreams
// They are implemented in the same way, use this function to unsubscribe already subscribed tables.
func (c *wsClient) UnsubscribeStreams(ctx context.Context, unSubs ...string) error {
	// This API is implemented the same way as of SubscribeStreams

	select {
	case <-c.done:
		return errors.New("client closed")
	default:
	}

	c.subsMu.Lock()
	defer func() {
		c.unSubsSignal = nil
		for len(c.unSubsData) > 0 {
			<-c.unSubsData
		}
		c.subsMu.Unlock()
	}()

	// preparing message for multiplexed socket connection
	message := make([]interface{}, 0, 4)
	payload := wsMessage{
		Op: "unsubscribe",
	}
	payload.Args = unSubs
	message = append(message, 0, c.config.Key, c.topic, payload)
	msgByte, _ := json.Marshal(message)

	ctxWait, cancelWait := context.WithCancel(ctx)
	defer cancelWait()

	go func() {
		select {
		case <-ctxWait.Done():
		case <-c.done:
			cancelWait()
		}
	}()

	if err := c.bucketM.WaitN(ctxWait, len(unSubs)); err != nil {
		return err
	}

	ctxConf, cancelConf := context.WithTimeout(ctxWait, wsConfirmTimeout)
	defer cancelConf()

	select {
	case <-c.done:
		return errors.New("client closed")
	default:
		select {
		case <-c.done:
			return errors.New("client closed")
		case c.connWriter <- msgByte:
		}
	}

	c.unSubsSignal = make(chan struct{})
	close(c.unSubsSignal)

	retChan := make(chan error)

	num := len(unSubs)

	for {
		select {
		case res := <-c.unSubsData:
			go func(ctxConf context.Context, retChan chan<- error, res []byte, payload wsMessage) {
				if request, ok := requestField(res); ok && validateSubsReq(request, payload) {
					if isSuccessful(res) {
						select {
						case retChan <- nil:
							return
						case <-ctxConf.Done():
							return
						}
					}
					select {
					case retChan <- wsError(res):
						return
					case <-ctxConf.Done():
						return
					}
				}
			}(ctxConf, retChan, res, payload)
		case err := <-retChan:
			num--
			if err != nil {
				if err.(APIError).StatusCode == 429 {
					return err
				}
			}
			if num == 0 {
				return nil
			}
		case <-ctxConf.Done():
			return errors.New("context canceled/timeout: message sent but could not verify")
		case <-c.done:
			return errors.New("client closed")
		}
	}
}

// Authenticate sends an authentication message over the websocket connection for the client.
// If it returns a non-nil error then most likely connection is now unsubscribed for the client, and you will
// need to use the AddWSClient function again
func (c *wsClient) Authenticate(ctx context.Context) error {
	// This API is implemented in the very similar way as that of SubscribeStreams
	// Authentication through websocket does consume a token from request limiter
	// It does not need num variable as there is only one expected message from websocket

	select {
	case <-c.done:
		return errors.New("client closed")
	default:
	}

	c.authMu.Lock()
	defer func() {
		c.authSignal = nil
		for len(c.authData) > 0 {
			<-c.authData
		}
		c.authMu.Unlock()
	}()

	// preparing message
	msg := make([]interface{}, 0, 4)
	apiExpires := time.Now().Add(RequestTimeout).Unix()
	signature := c.config.Sign("GET/realtime" + strconv.FormatInt(apiExpires, 10))
	payload := wsMessage{
		Op: "authKeyExpires",
	}
	payload.Args = []interface{}{c.config.Key, apiExpires, signature}
	msg = append(msg, 0, c.config.Key, c.topic, payload)
	msgByte, _ := json.Marshal(msg)

	ctxConf, cancelConf := context.WithTimeout(ctx, wsConfirmTimeout)
	defer cancelConf()

	select {
	case <-c.done:
		return errors.New("client closed")
	default:
		select {
		case <-c.done:
			return errors.New("client closed")
		case c.connWriter <- msgByte:
		}
	}

	c.authSignal = make(chan struct{})
	close(c.authSignal)

	retChan := make(chan error)

	for {
		select {
		case res := <-c.authData:
			go func(ctxConf context.Context, retChan chan<- error, res []byte, payload wsMessage) {
				if request, ok := requestField(res); ok && validateAuthReq(request, payload) {
					if isSuccessful(res) {
						select {
						case retChan <- nil:
							return
						case <-ctxConf.Done():
							return
						}
					}
					select {
					case retChan <- wsError(res):
						return
					case <-ctxConf.Done():
						return
					}
				}
			}(ctxConf, retChan, res, payload)
		case err := <-retChan:
			return err
		case <-ctxConf.Done():
			return errors.New("context canceled/timeout: message sent but could not verify")
			//case <-c.done:
			//	return errors.New("client done")
			// when authentication failed the socket unsubscribes the client which closes the c.done channel
			// this code is commented to prevent returning 'client done' error instead of more important
			// websocket error
		}
	}
}

// CancelAllAfter implements the cancelAllAfter subscription message of the websocket
func (c *wsClient) CancelAllAfter(ctx context.Context, timeout time.Duration) error {
	// This API is implemented in the very similar way as that of SubscribeStreams
	// This message consumes one token from the request limiter
	// It does not need num variable as there is only one expected message from websocket

	select {
	case <-c.done:
		return errors.New("client closed")
	default:
	}

	c.cancelMu.Lock()
	defer func() {
		c.cancelAfterSignal = nil
		for len(c.cancelAfterData) > 0 {
			<-c.cancelAfterData
		}
		c.cancelMu.Unlock()
	}()

	msg := make([]interface{}, 0, 4)
	payload := wsMessage{
		Op: "cancelAllAfter",
	}
	payload.Args = timeout.Milliseconds()
	msg = append(msg, 0, c.config.Key, c.topic, payload)
	msgByte, _ := json.Marshal(msg)

	ctxWait, cancelWait := context.WithCancel(ctx)
	defer cancelWait()

	go func() {
		select {
		case <-ctxWait.Done():
		case <-c.done:
			cancelWait()
		}
	}()

	if err := c.bucketM.Wait(ctxWait); err != nil {
		return err
	}

	ctxConf, cancelConf := context.WithTimeout(ctxWait, wsConfirmTimeout)
	defer cancelConf()

	select {
	case <-c.done:
		return errors.New("client closed")
	default:
		select {
		case <-c.done:
			return errors.New("client closed")
		case c.connWriter <- msgByte:
		}
	}

	c.cancelAfterSignal = make(chan struct{})
	close(c.cancelAfterSignal)

	retChan := make(chan error)

	// custom successful function for cancelAllAfter
	isSuccess := func(data []byte) bool {
		var res map[string]string
		_ = json.Unmarshal(data, &res)
		_, ok := res["cancelTime"]
		return ok
	}

	for {
		select {
		case res := <-c.cancelAfterData:
			go func(ctxConf context.Context, retChan chan<- error, res []byte, payload wsMessage) {
				if request, ok := requestField(res); ok {
					if validateCancelAllAfterReq(request, payload) {
						if isSuccess(res) {
							select {
							case retChan <- nil:
								return
							case <-ctxConf.Done():
								return
							}
						}
						select {
						case retChan <- wsError(res):
							return
						case <-ctxConf.Done():
							return
						}
					}
				}
			}(ctxConf, retChan, res, payload)
		case err := <-retChan:
			return err
		case <-ctxConf.Done():
			return errors.New("context canceled/timeout: message sent but could not verify")
		case <-c.done:
			return errors.New("client done")
		}
	}

}

//func (c *wsClient) AutoCancelAllAfter(ctx context.Context, timeout time.Duration, every time.Duration) {
//	burstOffset := int(math.Ceil(1 / every.Minutes()))
//	m := make(chan int)
//	go func() {
//		defer close(m)
//		for {
//			select {
//			case m <- burstOffset:
//			case <-c.done:
//			}
//		}
//	}()
//}

// UnsubscribeConnection gracefully closes the client without closing the websocket connection, it notifies the host
// that no data is required on this client and close the receiver channel.
func (c *wsClient) UnsubscribeConnection() {
	select {
	case <-c.done:
		return
	default:
	}

	msg := make([]interface{}, 0, 3)
	msg = append(msg, 2, c.config.Key, c.topic)
	msgByte, _ := json.Marshal(msg)

	select {
	case <-c.done:
		return
	default:
		select {
		case <-c.done:
			return
		case c.connWriter <- msgByte:
		}
	}

	go c.stop("outside caller")
}

// toWriterFunc sends the processed payload data to the one of the socket writer functions SubscribeStreams,
// UnsubscribeStreams, Authenticate, CancelAllAfter, these messages are sent so the writer function can verify if
// their request was successful or not signal channels are closed by the writer functions when the data
// sending starts
func (c *wsClient) toWriterFunc(data []byte) {
	select {
	case <-c.subsSignal:
		c.subsData <- data
	default:
	}

	select {
	case <-c.unSubsSignal:
		c.unSubsData <- data
	default:
	}

	select {
	case <-c.authSignal:
		c.authData <- data
	default:
	}

	select {
	case <-c.cancelAfterSignal:
		c.cancelAfterData <- data
	default:
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
func wsError(data []byte) APIError {
	var resStatus map[string]int
	var resError map[string]string

	_ = json.Unmarshal(data, &resStatus)
	_ = json.Unmarshal(data, &resError)

	bitmexErr := APIError{Name: "WebsocketError"}
	bitmexErr.StatusCode = resStatus["status"]
	bitmexErr.Message = resError["error"]

	return bitmexErr
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
// the request message sent by the functions SubscribeStreams and UnsubscribeStreams,
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

// validateCancelAllAfterReq verifies if the request field on the incoming socket message is the same as that of
// the request message sent by cancelAllAfter function, it returns true if the validation is successful.
func validateCancelAllAfterReq(req []byte, payload wsMessage) bool {
	request := struct {
		Op   string `json:"op,omitempty"`
		Args int    `json:"args,omitempty"`
	}{}

	if err := json.Unmarshal(req, &request); err == nil {
		if request.Op == payload.Op {
			if request.Args == payload.Args.(int) {
				return true
			}
		}
	}

	return false
}

// handleSocketMsg processes the message from websocket, extract the payload byte and send it over to sendCh.
// This function should not be called concurrently as an unsubscribeConnection message from websocket can
// provoke this function to call stop() which would prevent the previous messages from ever reaching to sendCh
// It is absolutely critical that the caller receives all messages for the client up until the last message
// of unsubscribe connection from the socket
func (c *wsClient) handleSocketMsg(msg []interface{}, logger *log.Logger) ([]byte, bool) {
	if len(msg) <= 3 {
		if v, ok := msg[0].(float64); ok {
			if v == 2 {
				go c.stop("websocket")
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
