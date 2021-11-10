package bitmex

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/bitmex-mirror/auth"
	"github.com/pkg/errors"
)

func TestErrorGrp(t *testing.T) {
	g, _ := errgroup.WithContext(context.Background())

	for i := 0; i < 10; i++ {
		i := i
		g.Go(func() error {
			fmt.Printf("%d\n", i)
			return nil
		})
	}
	err := g.Wait()
	fmt.Println(err)
}

func TestErrorGrp1(t *testing.T) {

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			fmt.Printf("%d\n", i)
			return
		}()
	}

	select {}
}

func TestBitmex_NewWSConnection2(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	b := NewBitmex(true)
	conn, err := b.NewWSConnection(context.Background(), logger)
	if err != nil {
		panic(err)
	}
	config := auth.NewConfig("YXTfYOwQn0e7gq3bziPyeNDa", "6tbsbjU1GjCGipthTHuLb2PrVswCuGBrA10JQ5IrhlGy0Ytz")
	client, err := conn.NewWSClient(context.Background(), config, "", logger)
	if err != nil {
		panic(err)
	}
	client.Connect()

	go func() {
		for i := range client.Receiver {
			logger.Println("Descpacito: ", string(i))
		}
	}()

	if err = client.Authenticate(context.Background()); err != nil {
		panic(err)
	}

	start := time.Now()
	err = client.SubscribeTables(context.Background(), "position", "margin", "order")
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	logger.Println("Tiempo 2: ", elapsed)

	//if err = client.Authenticate(context.Background()); err != nil {
	//	panic(err)
	//}
	//
	//err = client.SubscribeTables(context.Background(), "margin", "order", "position")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Got here")

	select {}
}

func TestBitmex_NewWSConnection(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)
	b := NewBitmex(true)
	conn, err := b.NewWSConnection(context.Background(), logger)
	if err != nil {
		panic(err)
	}
	config := auth.NewConfig("YXTfYOwQn0e7gq3bziPyeNDa", "6tbsbjU1GjCGipthTHuLb2PrVswCuGBrA10JQ5IrhlGy0Ytz")
	client, err := conn.NewWSClient(context.Background(), config, "", logger)
	if err != nil {
		panic(err)
	}
	client.Connect()

	go func() {
		for i := range client.Receiver {
			logger.Println("Descpacito: ", string(i))
		}
	}()

	var wg sync.WaitGroup
	wg.Add(50)

	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			if err = client.Authenticate(context.Background()); err != nil {
				panic(err)
			}
			//fmt.Println("k")
		}()
	}

	wg.Wait()
	select {}
	return

	if err = client.Authenticate(context.Background()); err != nil {
		panic(err)
	}

	start := time.Now()
	err = client.SubscribeTables(context.Background(), "margin", "order", "position")
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	logger.Println("Tiempo 1: ", elapsed)

	//if err = client.Authenticate(context.Background()); err != nil {
	//	panic(err)
	//}
	//
	//err = client.SubscribeTables(context.Background(), "margin", "order", "position")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Got here")

	select {}
}

func TestRateLimited(t *testing.T) {
	jsonStr := `{"status":429,"error":"Rate limit exceeded, retry in 1 seconds.","meta":{"retryAfter":1},"request":{"op":"subscribe","args":"orderBookL2_25"}}`
	err := wsError([]byte(jsonStr))
	fmt.Println(err)
}

func wsError(data []byte) error {
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
		return errors.Wrap(errors.Wrap(bitmexErr, "wsError: error field not found"), errStr)
	}
	bitmexErr.Message = res["error"].(string)

	if bitmexErr.StatusCode == 401 || bitmexErr.StatusCode == 403 {
		//TODO
	}

	if bitmexErr.StatusCode == 429 {
		rateLimited := false

		if m, ok := res["meta"]; ok {
			if meta, ok := m.(map[string]interface{}); ok {
				if wait, ok := meta["retryAfter"]; ok {
					if wait, ok := wait.(float64); ok {
						fmt.Println(wait, "Reached at the right place")
						rateLimited = true
					}
				}
			}
		}

		if !rateLimited {
			fmt.Println("Reached at the wrong place")
		}
		//TODO
	}

	return bitmexErr
}

//
//import (
//	"context"
//	"fmt"
//	"log"
//	"os"
//	"testing"
//	"time"
//
//	"github.com/bitmex-mirror/auth"
//)
//
//func TestConnectCon(t *testing.T) {
//	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.LUTC|log.Lshortfile)
//	//logger.Println("sadf")
//	//os.Exit(0)
//	ws, err := NewWSConnection(context.Background(), true, logger)
//	if err != nil {
//		panic(err)
//	}
//
//	config := auth.NewConfig("YXTfYOwQn0e7gq3bziPyeNDa", "6tbsbjU1GjCGipthTHuLb2PrVswCuGBrA10JQ5IrhlGy0Ytz")
//
//	//con := make([]auth.Config, 0, 1000)
//	//
//	//for i := 0; i < 1000; i++ {
//	//	con = append(con, auth.NewConfig(fmt.Sprintf("asfiuberiviu%d", i), ""))
//	//}
//	//
//	//for i := range con {
//	//	go func() {
//	//		_, _ = ws.NewWSClient(context.Background(), con[i], "", logger)
//	//	}()
//	//}
//	//
//	//time.Sleep(time.Second * 100)
//	//return
//	c, r := ws.NewWSClient(context.Background(), config, "", logger)
//
//	go func() {
//		for {
//			v, ok := <-r
//			if !ok {
//				return
//			}
//			fmt.Println(string(v))
//		}
//	}()
//
//	fmt.Println("got")
//	err = c.Authenticate(context.Background())
//	fmt.Println(err)
//	err = c.SubscribeTables(context.Background(), "position", "order", "margin")
//	fmt.Println(err)
//
//	select {}
//
//	panic("sad")
//	fmt.Println("got")
//	fmt.Println(err)
//	if err != nil {
//		panic(err)
//	}
//	//err = c.Authenticate()
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	for {
//		//fmt.Println("ff")
//		err = c.CancelAllAfter(context.Background(), time.Second)
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//	}
//
//	return
//
//	//err = c.SubscribeTables("order", "position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//err = c.UnsubscribeTables("order", "position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	////time.Sleep(time.Second)
//	//
//	//return
//	//
//	//err = c.UnsubscribeTables("position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//time.Sleep(time.Second)
//	//
//	//err = c.SubscribeTables("position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	////err = c.Authenticate()
//	////err = c.CancelAllAfter(time.Second * 10)
//	//fmt.Println("Time taken: ", time.Now().UnixNano()-start)
//	//fmt.Println(err)
//	//return
//	//
//	//err = c.SubscribeTables("order")
//	//fmt.Println("Done")
//	//err = c.UnsubscribeTables("order")
//	//
//	//fmt.Println(err)
//	//return
//	//
//	//go func() {
//	//	for {
//	//		msg, ok := <-r
//	//		if !ok {
//	//			logger.Println("reader done")
//	//			return
//	//		}
//	//		logger.Println("Message received")
//	//		logger.Println(string(msg))
//	//	}
//	//}()
//	//
//	//err = c.SubscribeTables("order")
//	//fmt.Println("Subs Error: ", err)
//	////go func() {
//	////	time.Sleep(time.Second)
//	////	_, err = rc.GetOrdersRequest().Symbol("XBTUSD").Do()
//	////	if err != nil {
//	////		panic(err)
//	////	}
//	////	for i := 0; i < 20; i++ {
//	////		time.Sleep(time.Millisecond * 100)
//	////		c.CancelAllAfter(time.Second / time.Millisecond)
//	////	}
//	////	_, err = rc.GetOrdersRequest().Symbol("XBTUSD").Do()
//	////	if err != nil {
//	////		panic(err)
//	////	}
//	////	//time.Sleep(time.Second * 30)
//	////	//logger.Println("Stopping now")
//	////	//c.UnsubscribeConnection()
//	////}()
//	//
//	////logger.Println("sending auth")
//	////c.Authenticate()
//	select {}
//	//cancel()
//
//}
