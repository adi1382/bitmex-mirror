package bitmex

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bitmex-mirror/auth"
	"go.uber.org/atomic"
)

type wsMessage struct {
	Op   string        `json:"op,omitempty"`
	Args []interface{} `json:"args,omitempty"`
}

type wsClient struct {
	chClientRead    chan []byte
	connWriter      chan<- []byte
	config          auth.Config
	topic           string
	subs            []string
	subsMu          sync.Mutex
	timeout         time.Duration
	isAuthenticated bool
	isRemoved       atomic.Bool
	closeMu         sync.Mutex
	done            chan struct{}
	closeFunc       sync.Once
}

func (c *wsClient) SubscribeStreams(subs ...string) {
	c.subsMu.Lock()
	defer c.subsMu.Unlock()

	message := make([]interface{}, 0, 4)
	payload := wsMessage{
		Op: "subscribe",
	}
	for i := range subs {
		payload.Args = append(payload.Args, subs[i])
	}
	message = append(message, 0, c.config.Key, c.topic, payload)
	ml, _ := json.Marshal(message)
	c.connWriter <- ml

	c.subs = append(c.subs, subs...)
}

func (c *wsClient) UnsubscribeStreams(unSubs ...string) {
	c.subsMu.Lock()
	defer c.subsMu.Unlock()

	message := make([]interface{}, 0, 4)
	payload := wsMessage{
		Op: "unsubscribe",
	}
	for i := range unSubs {
		payload.Args = append(payload.Args, unSubs[i])
	}
	message = append(message, 0, c.config.Key, c.topic, payload)
	ml, _ := json.Marshal(message)
	c.connWriter <- ml

	for i := range unSubs {
		for j := range c.subs {
			if c.subs[j] == unSubs[i] {
				c.subs[j] = c.subs[len(c.subs)-1]
				c.subs = c.subs[:len(c.subs)-1]
				break
			}
		}
	}
}

func (c *wsClient) UnsubscribeConnection() {
	message := make([]interface{}, 0, 4)
	message = append(message, 2, c.config.Key, c.topic, "")
	ml, _ := json.Marshal(message)
	c.connWriter <- ml
	go c.remove()
}

func (c *wsClient) Authenticate() {
	msg := make([]interface{}, 0, 4)
	apiExpires := time.Now().Add(time.Second * 10).Unix()
	signature := c.config.Sign("GET/realtime" + strconv.FormatInt(apiExpires, 10))

	payload := wsMessage{
		Op: "authKeyExpires",
	}
	payload.Args = append(payload.Args, c.config.Key, apiExpires, signature)
	msg = append(msg, 0, c.config.Key, c.topic, payload)

	marshaled, _ := json.Marshal(msg)
	fmt.Println(string(marshaled))
	c.connWriter <- marshaled

	c.isAuthenticated = true
}

func (c *wsClient) AutoCancelAfter(timeout time.Duration) {
	m := make([]interface{}, 0, 4)

	message := wsMessage{
		Op: "cancelAllAfter",
	}
	message.Args = append(message.Args, timeout)
	marshal, _ := json.Marshal(message)
	m = append(m, 0, c.config.Key, c.topic, marshal)

	c.timeout = timeout
}

func (c *wsClient) remove() {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()

	c.closeFunc.Do(func() {
		c.isRemoved.Store(true)
		close(c.done)
		close(c.chClientRead)
	})
}

func (c *wsClient) onCtxDone(ctx context.Context) {
	select {
	case <-ctx.Done():
		c.remove()
		c.UnsubscribeConnection()
		return
	case <-c.done:
		return
	}
}

func (c *wsClient) sendMessage(data []byte) bool {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()

	select {
	case <-c.done:
		return false
	default:
	}

	if string(data) == "quit" {
		go c.remove()
		return false
	}

	c.chClientRead <- data
	return true
}
