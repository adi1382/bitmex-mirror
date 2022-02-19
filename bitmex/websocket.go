package bitmex

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"golang.org/x/time/rate"
)

// NewWSConnection create a new multiplexed websocket connection and returns a variable to provide method set to
// add and authenticate several account on the connection.
// context can be used to set a timeout to establish the connection, default is 15 seconds.
// canceling the context will return this function immediately, but does not guarantee that the connection is closed.
// if this function returns a nil error, then you should call the Close() method on the returned *WSConnection
// to close the connection properly freeing up resources.
func (b *Bitmex) NewWSConnection(ctx context.Context, logger *log.Logger) (*WSConnection, error) {
	var ws WSConnection

	if b.test {
		ws.host = wsTestEndpoint
	} else {
		ws.host = wsEndpoint
	}

	ws.bucketM = b.bucketM

	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, websocketConnectionTimeout)
		defer cancel()
	}

	conn, err := connect(ctx, ws.host)
	if err != nil {
		return nil, errors.Wrap(err, "ws connection failed to connect")
	}

	done := make(chan struct{})
	ws.Done = done
	ws.conn = conn
	ws.closing = make(chan string)

	chWrite := make(chan []byte, 1)
	chRead := make(chan []byte, 1)
	ws.clients = make(map[string]*WSClient, 5)
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

	go ws.connectionReader(chRead, done, logger)
	go ws.connectionWriter(chWrite, logger)
	go ws.connectionRouter(logger)

	return &ws, err
}

type WSConnection struct {
	clientsMu    sync.RWMutex         // RW lock for clients
	clients      map[string]*WSClient // hashmap of all clients on the multiplexed socket connection
	host         string               // testnet or mainnet host url
	conn         *websocket.Conn      // gorilla websocket connection
	bucketM      *rate.Limiter        // rate limit bucket shared with rest api calls
	chWrite      chan<- []byte        // channel to write websocket outgoing messages
	chRead       <-chan []byte        // channel to read websocket incoming messages
	pingPeriod   time.Duration        // period for ticker - input
	pongTimeout  time.Duration        // time after sending ping to get pong - input
	pingTicker   *time.Ticker         // time to send ping after receiving last message
	pongTimer    *time.Timer          // Expires when pong is not received
	Done         <-chan struct{}      // closes when socket connection drops
	closing      chan string          // intermediary channel to close the connection
	removeClient chan string          // receives a message from client to remove its api from clients
}

// IsOpen returns true if the websocket connection is open.
func (ws *WSConnection) IsOpen() bool {
	select {
	case <-ws.Done:
		return false
	default:
		return true
	}
}

func (ws *WSConnection) Close() {
	ws.close("outside caller")
}

func (ws *WSConnection) close(by string) {
	select {
	case ws.closing <- by:
		<-ws.Done
	case <-ws.Done:
	}
}

// connectionReader reads incoming messages from the websocket connection and sends them to the router.
// It must be run as a goroutine.
func (ws *WSConnection) connectionReader(chRead chan<- []byte, Done chan struct{}, logger *log.Logger) {
	defer func() {
		close(Done)
	}()

	for {
		_, message, err := ws.conn.ReadMessage()

		ws.pingTicker.Reset(ws.pingPeriod)
		if err != nil {
			logger.Println("websocket: error: ", err)
			logger.Println("websocket: info: ws read channel done")
			return
		}
		chRead <- message
		//logger.Println(">>>", string(message))
	}
}

// connectionWriter writes outgoing messages to the websocket connection.
// It also handles ping/pong messages, it must be run on its own goroutine.
func (ws *WSConnection) connectionWriter(chWrite <-chan []byte, logger *log.Logger) {

	defer func() {
		// closing the connection will make the connectionReader return as ws.conn.ReadMessage() will return an error
		_ = ws.conn.Close()
	}()

	for {
		select {
		case by := <-ws.closing:
			logger.Println("websocket: info: ws connection is being closed by: ", by)
			return
		case <-ws.Done:
			logger.Println("websocket: info: ws writer returns because connection closed")
			return
		case message := <-chWrite:
			_ = ws.conn.SetWriteDeadline(time.Now().Add(wsWriteTimeout))
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
			_ = ws.conn.SetWriteDeadline(time.Now().Add(wsWriteTimeout))
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

func (ws *WSConnection) connectionRouter(logger *log.Logger) {

	for {
		select {
		case <-ws.Done:
			// This channel is closed when the connectionReader returns
			logger.Println("router returns because ws connection closed")
			return
		case apiKey := <-ws.removeClient:
			// deleting client's API key from clients map on request of WSClient variable
			ws.clientsMu.Lock()
			delete(ws.clients, apiKey)
			ws.clientsMu.Unlock()
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

			// determining WSClient variable for apiKey
			ws.clientsMu.RLock()
			wsCl, ok := ws.clients[apiKey]
			ws.clientsMu.RUnlock()

			// if apiKey is not found, unknown connection will be unsubscribed
			if !ok {
				msgType := message[0].(float64)
				//
				//// check message type, 2 is for Unsubscribing connection
				if msgType != 2 {
					logger.Println("websocket: error: ws connection received message for invalid api: ", apiKey)
				}
				//	topic, _ := message[2].(string)
				//
				//	// preparing message for unsubscribing connection
				//	writeMsg := make([]interface{}, 0, 3)
				//	writeMsg = append(writeMsg, 2, apiKey, topic)
				//	writeMsgByte, _ := json.Marshal(writeMsg)
				//
				//	// sending message to socket writer
				//	select {
				//	case ws.chWrite <- writeMsgByte:
				//	case <-ws.Done:
				//	}
				//}
				continue
			}

			// send message to socketMessage channel operated by WSClient variable
			// if the client is closed send on socketMessage could block forever
			select {
			case wsCl.socketMessage <- message:
			case <-wsCl.Done:
			case <-ws.Done:
			}
		}
	}
}

// Dials a websocket connection to the specified URL.
func connect(ctx context.Context, host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/realtimemd"}
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	return conn, err
}
