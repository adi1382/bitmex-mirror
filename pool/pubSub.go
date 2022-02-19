package pool

import "sync"

func NewPubSub() *pubSub {
	return &pubSub{
		unsubscribing: make(chan struct{}, 1),
	}
}

type pubSub struct {
	subs []struct {
		ch   chan []byte
		done chan struct{}
	} // slice of subscribers
	subsMu        sync.Mutex    // protects the subs slice
	unsubscribing chan struct{} // buffered channel of 1 to avoid multiple unsubscribe run
}

// Subscribe subscribes the ch channel to the main receiver channel and all data for the client is copied and forwarded
// to the ch channel.
// it returns an unsubscribe function which must be called to unsubscribe the client from the main receiver channel.
// Failing to do so, will result in deadlocks.
func (c *pubSub) Subscribe(ch chan []byte) func() {
	done := make(chan struct{})
	var once sync.Once
	unsubscribe := func() { once.Do(func() { close(done) }) }
	c.subsMu.Lock()
	defer c.subsMu.Unlock()
	c.subs = append(c.subs, struct {
		ch   chan []byte
		done chan struct{}
	}{ch, done})
	return unsubscribe
}

// publish should be used by the manager to send client's websocket data.
// The msg will be sent to all active subscriptions.
func (c *pubSub) publish(msg []byte) {
	c.subsMu.Lock()
	defer c.subsMu.Unlock()

	if len(c.subs) == 0 {
		return
	}

	for i := range c.subs {
		select {
		case <-c.subs[i].done:
			select {
			case c.unsubscribing <- struct{}{}:
				go c.removeClosed()
			default:
			}
		case c.subs[i].ch <- msg:
		}
	}
}

func (c *pubSub) removeClosed() {
	c.subsMu.Lock()
	defer func() {
		<-c.unsubscribing
		c.subsMu.Unlock()
	}()

	k := 0
	for i := range c.subs {
		select {
		case <-c.subs[i].done:
		default:
			c.subs[k] = c.subs[i]
			k++
		}
	}
	c.subs = c.subs[:k]
}
