package bitmex

import (
	"context"
	"fmt"
	"github.com/bitmex-mirror/auth"
	"log"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestConnectCon(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	logger.Println("sadf")
	//os.Exit(0)
	ctx, cancelFunc := context.WithCancel(context.Background())
	ws, err := NewWSConnection(ctx, true, logger)
	ctx2, cancel := context.WithCancel(ctx)
	if err != nil {
		panic(err)
	}

	config := auth.NewConfig("pkOAYJNujj_-cW3gN0FjPizp", "LgiVEp09S4TFw9FfoMuurP-7lVJ6DPfbCRHWhomkqn7qx-F4")
	//time.Sleep(time.Second * 30)
	c, r := ws.AddNewClient(ctx2, config, "")

	go func() {
		for {
			msg, ok := <-r
			if !ok {
				return
			}
			fmt.Println("Message received")
			fmt.Println(string(msg))
		}
	}()
	c.Authenticate()
	time.Sleep(time.Second * 10)
	fmt.Println(runtime.NumGoroutine())
	fmt.Println("context canceled")
	cancelFunc()
	time.Sleep(time.Second)
	fmt.Println(runtime.NumGoroutine())
	//select {}
	cancel()

}
