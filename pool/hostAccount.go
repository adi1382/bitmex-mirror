package pool

//
//import (
//	"context"
//	"fmt"
//	"github.com/bitmex-mirror/auth"
//	"github.com/bitmex-mirror/bitmex"
//	"github.com/pkg/errors"
//	"log"
//	"sync"
//)
//
//func NewHostAccount(ctx context.Context, config auth.Config,
//	topic string, id bool, logger *log.Logger) (*hostAccount, error) {
//
//	host := hostAccount{
//		done: make(chan struct{}),
//	}
//
//	host.bitmex = bitmex.NewBitmex(id)
//
//	host.Account = &Account{
//		activeOrders:  make([]bitmex.Order, 0, 10),
//		positions: make([]bitmex.Position, 0, 10),
//		margins:       make([]bitmex.Margin, 0, 2),
//		restClient:    host.bitmex.NewRestClient(config),
//	}
//
//	var err error
//
//	host.connection, err = host.bitmex.NewWSConnection(ctx, logger)
//	if err != nil {
//		return nil, err
//	}
//
//	host.config = config
//	host.topic = topic
//	wsClient, socketReceiver := host.connection.SubscribeNewClient(ctx, config, topic, logger)
//
//	wsReceiver := make(chan []byte, 100)
//	go host.router(ctx, socketReceiver, wsReceiver)
//	host.wsUpdater = wsReceiver
//
//	err = wsClient.Authenticate(ctx)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to authenticate")
//	}
//
//	err = wsClient.SubscribeTables(ctx, bitmex.WSTableOrder, bitmex.WSTablePosition, bitmex.WSTableMargin)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to subscribe tables")
//	}
//
//	go host.updateFromSocket(ctx, logger)
//
//	return &host, nil
//}
//
//type hostAccount struct {
//	config     auth.Config
//	topic      string
//	bitmex     *bitmex.Bitmex
//	connection *bitmex.WSConnection
//	subs       []chan<- []byte
//	subsMu     sync.Mutex
//	done       chan struct{}
//	*Account
//}
//
//func (h *hostAccount) AddSubAccount(ctx context.Context, config auth.Config, logger *log.Logger) error {
//	sub := subAccount{}
//
//	sub.Account = &Account{
//		activeOrders:  make([]bitmex.Order, 0, 10),
//		positions: make([]bitmex.Position, 0, 10),
//		margins:       make([]bitmex.Margin, 0, 2),
//		restClient:    h.bitmex.NewRestClient(config),
//	}
//
//	var err error
//
//	wsClient, receiver := h.connection.SubscribeNewClient(ctx, config, "", logger)
//	sub.wsUpdater = receiver
//
//	go sub.updateFromSocket(ctx, logger)
//
//	err = wsClient.Authenticate(ctx)
//	if err != nil {
//		return errors.Wrap(err, "failed to authenticate")
//	}
//
//	err = wsClient.SubscribeTables(ctx, bitmex.WSTableOrder, bitmex.WSTablePosition, bitmex.WSTableMargin)
//	if err != nil {
//		return errors.Wrap(err, "failed to subscribe tables")
//	}
//
//	hostReceiver := make(chan []byte, 100)
//	h.subscribe(hostReceiver)
//	sub.hostReceiver = hostReceiver
//
//	//TODO: REMOVE THIS AFTER TESTS
//	go func() {
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			case msg := <-hostReceiver:
//				fmt.Println("message from host: ", string(msg))
//			}
//		}
//	}()
//
//	return nil
//}
//
//func (h *hostAccount) restart(ctx context.Context, logger *log.Logger) (chan []byte, bool) {
//	wsClient, socketReceiver := h.connection.SubscribeNewClient(ctx, h.config, h.topic, logger)
//}
//
//func (h *hostAccount) router(ctx context.Context, socketReceiver <-chan []byte,
//	wsReceiver chan<- []byte, logger *log.Logger) {
//	for {
//		select {
//		case msg, ok := <-socketReceiver:
//			if !ok {
//				if h.connection.IsOpen() {
//					socketReceiver, success := h.restart(ctx, logger)
//					continue
//				}
//				h.close(wsReceiver)
//				return
//			}
//			for i := range h.subs {
//				h.subs[i] <- msg
//			}
//			wsReceiver <- msg
//
//		case <-ctx.Done():
//			//TODO: Unsubscribe
//			h.close(wsReceiver)
//			return
//		}
//	}
//}
//
//// once close is called it's the end of the hostAccount variable forever
//// new hostAccount should be created from operator if needed
//func (h *hostAccount) close(wsReceiver chan<- []byte) {
//	close(h.done)
//	h.subsMu.Lock()
//	defer h.subsMu.Unlock()
//	for i := range h.subs {
//		close(h.subs[i])
//	}
//	close(wsReceiver)
//}
//
//func (h *hostAccount) subscribe(ch chan<- []byte) {
//	h.subsMu.Lock()
//	defer h.subsMu.Unlock()
//	h.subs = append(h.subs, ch)
//}
