package bitmex

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/bitmex-mirror/auth"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"golang.org/x/time/rate"
)

type wSConnection struct {
	clLock       sync.RWMutex         // RW lock for clients
	clients      map[string]*wsClient // hashmap of all clients on the multiplexed socket connection
	host         string               // testnet or mainnet host url
	conn         *websocket.Conn      // gorilla websocket connection
	bucketM      *rate.Limiter        // rate limit bucket shared with rest api calls
	chWrite      chan<- []byte        // channel to write websocket outgoing messages
	chRead       <-chan []byte        // channel to read websocket incoming messages
	pingPeriod   time.Duration        // period for ticker - input
	pongTimeout  time.Duration        // time after sending ping to get pong - input
	pingTicker   *time.Ticker         // time to send ping after receiving last message
	pongTimer    *time.Timer          // Expires when pong is not received
	done         chan struct{}        // closes when socket connection drops
	removeClient chan string          // receives a message from client to remove its api from clients
}

// AddWSClient subscribes a new stream for the client over the multiplexed websocket connection
// it returns a variable of type *wsClient which implements methods to subscribe, unsubscribe, authenticate and write
// standard message over the subscribed stream.
// These methods provide a rest-like interface which returns errors received over the subscribed stream.
func (ws *wSConnection) AddWSClient(ctx context.Context, config auth.Config, topic string, logger *log.Logger) (*wsClient, <-chan []byte) {

	c := wsClient{
		config:     config,
		topic:      topic,
		connWriter: ws.chWrite,
		bucketM:    ws.bucketM,
	}

	c.returnCh = make(chan []byte, 1000)

	c.done = make(chan struct{})
	c.closing = make(chan string)

	c.subsData = make(chan []byte, 10)
	c.unSubsData = make(chan []byte, 10)
	c.authData = make(chan []byte, 10)
	c.cancelAfterData = make(chan []byte, 10)
	c.socketMessage = make(chan []interface{}, 10)
	c.bucketM = ws.bucketM

	select {
	case <-ws.done:
		ret := make(chan []byte)
		close(ret)
		return &c, ret
	default:
	}

	go c.manager(ctx, logger, ws.done, ws.removeClient)

	ws.clLock.Lock()
	ws.clients[config.Key] = &c
	ws.clLock.Unlock()

	msg := make([]interface{}, 0, 4)
	msg = append(msg, 1, config.Key, topic)
	msgBytes, _ := json.Marshal(msg)

	select {
	case c.connWriter <- msgBytes:
	case <-c.done:
	}

	return &c, c.returnCh
}

func (ws *wSConnection) IsOpen() bool {
	select {
	case <-ws.done:
		return true
	default:
		return false
	}
}

func (ws *wSConnection) connectionReader(ctx context.Context,
	chRead chan<- []byte, chWrite <-chan []byte, logger *log.Logger) {
	defer func() {
		close(chRead)
		close(ws.done)
	}()

	ctxC, cancel := context.WithCancel(ctx)
	defer cancel()

	go ws.connectionWriter(ctxC, chWrite, logger)
	go ws.router(ctxC, logger)

	for {
		_, message, err := ws.conn.ReadMessage()

		ws.pingTicker.Reset(ws.pingPeriod)
		if err != nil {
			logger.Println("websocket: error: ", err)
			logger.Println("websocket: info: ws read channel done")
			return
		}
		chRead <- message
		//logger.Println("raw message: ", string(message))
	}
}

func (ws *wSConnection) connectionWriter(ctx context.Context, chWrite <-chan []byte, logger *log.Logger) {

	writeWait := 10 * time.Second
	fastWriteWait := time.Second

	defer func() {
		_ = ws.conn.Close()
		for len(chWrite) > 0 {
			<-chWrite
		}
	}()

	for {
		select {
		case <-ctx.Done():
			_ = ws.conn.SetWriteDeadline(time.Now().Add(fastWriteWait))
			_ = ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
			logger.Println("websocket: info: context done socket writer")
			return
		case message := <-chWrite:

			_ = ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := ws.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger.Println("websocket: error: failed to write on ws")
				return
			}

			for len(ws.chWrite) > 0 {
				message = <-chWrite
				err = ws.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					logger.Println("websocket: error: failed to write on ws")
					return
				}
			}

		case <-ws.pingTicker.C:
			if ws.pongTimer.Stop() {
				logger.Println("websocket: error: pong timer was not stopped by pong handler")
				return
			}

			_ = ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Println("websocket: error: failed to write ping on ws")
				return
			}

			ws.pongTimer.Reset(time.Second * 5)

		case <-ws.pongTimer.C:
			logger.Println("websocket: error: pong timeout expired", time.Now())
			return
		}
	}
}

func (ws *wSConnection) router(ctx context.Context, logger *log.Logger) {

	for {
		select {
		case <-ctx.Done():
			// This context gets canceled when the parent context is canceled passed into NewWSConnection func
			// or when connectionReader function returns
			// this is the only return point for the router routine
			logger.Println("websocket: info: context done router")
			return
		case apiKey := <-ws.removeClient:
			// deleting client's API key from clients map on request of wsClient variable
			ws.clLock.RLock()
			delete(ws.clients, apiKey)
			ws.clLock.RUnlock()
		case msg := <-ws.chRead:

			// msg format for multiplexed websocket message
			message := make([]interface{}, 0, 4)
			err := json.Unmarshal(msg, &message)

			if err != nil {
				logger.Println(errors.Wrap(errors.Wrap(err, string(msg)),
					"websocket: error: incoming msg unmarshal failed"))
				continue
			}

			// determining the apiKey on multiplexed client variable
			apiKey, ok := message[1].(string)
			if !ok {
				logger.Println(errors.Wrap(errors.Wrap(err, string(msg)),
					"websocket: error: incoming msg does not contain key"))
				continue
			}

			// determining wsClient variable for apiKey
			ws.clLock.RLock()
			wsCl, ok := ws.clients[apiKey]
			ws.clLock.RUnlock()

			// if apiKey is not found, unknown connection will be unsubscribed
			if !ok {
				msgType := message[0].(float64)

				// check message type, 2 is for Unsubscribing connection
				if msgType != 2 {
					topic, _ := message[2].(string)

					// preparing message for unsubscribing connection
					writeMsg := make([]interface{}, 0, 3)
					writeMsg = append(writeMsg, 2, apiKey, topic)
					writeMsgByte, _ := json.Marshal(writeMsg)

					// sending message to socket writer
					select {
					case ws.chWrite <- writeMsgByte:
					case <-ws.done:
					}
				}

				logger.Println("websocket: error: ws connection received message for invalid api: ", apiKey)
				continue
			}

			// send message to socketMessage channel operated by wsClient variable
			// if the client is closed send on socketMessage could block forever
			select {
			case wsCl.socketMessage <- message:
			case <-wsCl.done:
			}
		}
	}
}

func connect(ctx context.Context, host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/realtimemd"}
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	return conn, err
}
