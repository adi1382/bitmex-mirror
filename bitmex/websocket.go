package bitmex

import (
	"context"
	"encoding/json"
	"github.com/bitmex-mirror/auth"
	"github.com/gorilla/websocket"
	"github.com/juju/ratelimit"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"time"
)

func NewWSConnection(ctx context.Context, test bool, logger *log.Logger) (*wSConnection, error) {
	ws := wSConnection{
		path: "/realtimemd",
	}

	if test {
		ws.host = "ws.testnet.bitmex.com"
	} else {
		ws.host = "ws.bitmex.com"
	}

	err := ws.connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "ws connection failed to connect")
	}

	chWrite := make(chan []byte, 100)
	chRead := make(chan []byte, 100)

	ws.chRead = chRead
	ws.chWrite = chWrite
	ws.pingPeriod = time.Second * 5
	ws.pongTimeout = time.Second * 5

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

	go ws.connectionReader(chRead, logger)
	go ws.connectionWriter(ctx, chWrite, logger)
	go ws.router(ctx, logger)

	return &ws, err
}

type wSConnection struct {
	clients     map[string]*wsClient
	host        string
	path        string
	conn        *websocket.Conn
	bucket1m    *ratelimit.Bucket
	chWrite     chan<- []byte
	chRead      <-chan []byte
	pingPeriod  time.Duration
	pongTimeout time.Duration
	pingTicker  *time.Ticker
	pongTimer   *time.Timer
}

func (ws *wSConnection) AddNewClient(ctx context.Context, config auth.Config, topic string) (*wsClient, <-chan []byte) {
	c := wsClient{
		config:     config,
		topic:      topic,
		connWriter: ws.chWrite,
	}

	c.chClientRead = make(chan []byte, 100)

	msg := make([]interface{}, 0, 4)
	msg = append(msg, 1, config.Key, topic)
	marshaled, _ := json.Marshal(msg)
	c.connWriter <- marshaled
	c.done = make(chan struct{})

	go c.onCtxDone(ctx)

	ws.clients = make(map[string]*wsClient, 5)
	ws.clients[config.Key] = &c
	return &c, c.chClientRead
}

func (ws *wSConnection) connect(ctx context.Context) error {
	u := url.URL{Scheme: "wss", Host: ws.host, Path: ws.path}
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	ws.conn = conn
	return err
}

func (ws *wSConnection) connectionWriter(ctx context.Context, chWrite <-chan []byte, logger *log.Logger) {

	writeWait := 10 * time.Second
	fastWriteWait := time.Second

	defer func() {
		_ = ws.conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			_ = ws.conn.SetWriteDeadline(time.Now().Add(fastWriteWait))
			_ = ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
			logger.Println("websocket: info: context closed socket writer")
			return
		case message, ok := <-chWrite:
			// This is implemented, however program should never close "chWrite"
			// Because data is sent from many go routines
			// A quit message is received from socket reader to close the writer
			if !ok {
				_ = ws.conn.SetWriteDeadline(time.Now().Add(fastWriteWait))
				_ = ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				logger.Println("websocket: info: found ws write channel closed")
				return
			}

			if string(message) == "quit" {
				logger.Println("websocket: error: quit received on socket writer")
				return
			}

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

func (ws *wSConnection) connectionReader(chRead chan<- []byte, logger *log.Logger) {

	defer func() {
		_ = ws.conn.Close()
		close(chRead)
	}()

	for {
		_, message, err := ws.conn.ReadMessage()

		ws.pingTicker.Reset(ws.pingPeriod)
		if err != nil {
			logger.Println("websocket: error: ", err)
			logger.Println("websocket: info: ws read channel closed")
			ws.chWrite <- []byte("quit")
			return
		}
		chRead <- message
	}
}

func (ws *wSConnection) router(ctx context.Context, logger *log.Logger) {

	for {
		select {
		case <-ctx.Done():
			logger.Println("websocket: info: context closed router")
			for i := range ws.clients {
				ws.clients[i].sendMessage([]byte("quit"))
			}
			return
		case msg, ok := <-ws.chRead:

			if !ok {
				logger.Println("websocket: error: router closed for closed chRead")
				for i := range ws.clients {
					ws.clients[i].sendMessage([]byte("quit"))
				}
				return
			}

			go func(msg []byte) {
				message := make([]interface{}, 0, 4)
				err := json.Unmarshal(msg, &message)

				if err != nil {
					logger.Println(errors.Wrap(errors.Wrap(err, string(msg)),
						"websocket: error: incoming msg unmarshal failed"))
					return
				}

				apiKey, ok := message[1].(string)
				if !ok {
					logger.Println(errors.Wrap(errors.Wrap(err, string(msg)),
						"websocket: error: incoming msg does not contain key"))
					return
				}

				payload, err := json.Marshal(message[3])
				if err != nil {
					logger.Println(errors.Wrap(errors.Wrap(err, string(msg)),
						"websocket: error: no payload on incoming message"))
					return
				}

				c, ok := ws.clients[apiKey]

				if !ok {
					logger.Println("websocket: error: ws connection received message for invalid bitmex")
					return
				}

				if !c.isRemoved.Load() {
					c.chClientRead <- payload
					return
				} else {
					if m, ok := message[0].(int); ok {
						if m == 2 {
							delete(ws.clients, apiKey)
							return
						}
					}
				}

			}(msg)

		}
	}
}
