package bitmex

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
//	//		_, _ = ws.AddWSClient(context.Background(), con[i], "", logger)
//	//	}()
//	//}
//	//
//	//time.Sleep(time.Second * 100)
//	//return
//	c, r := ws.AddWSClient(context.Background(), config, "", logger)
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
//	err = c.SubscribeStreams(context.Background(), "position", "order", "margin")
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
//	//err = c.SubscribeStreams("order", "position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//err = c.UnsubscribeStreams("order", "position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	////time.Sleep(time.Second)
//	//
//	//return
//	//
//	//err = c.UnsubscribeStreams("position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//time.Sleep(time.Second)
//	//
//	//err = c.SubscribeStreams("position")
//	//if err != nil {
//	//	panic(err)
//	//}
//	////err = c.Authenticate()
//	////err = c.CancelAllAfter(time.Second * 10)
//	//fmt.Println("Time taken: ", time.Now().UnixNano()-start)
//	//fmt.Println(err)
//	//return
//	//
//	//err = c.SubscribeStreams("order")
//	//fmt.Println("Done")
//	//err = c.UnsubscribeStreams("order")
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
//	//err = c.SubscribeStreams("order")
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
