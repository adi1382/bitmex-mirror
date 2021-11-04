package bitmex

//
//import (
//	"context"
//	"fmt"
//	"testing"
//	"time"
//
//	"github.com/bitmex-mirror/auth"
//)
//
////func BenchmarkOpt(b *testing.B) {
////	c := restClient{}
////
////	for n := 0; n < b.N; n++ {
////		_, err := c.AmendOrderRequest().OrderID("DF").OrderQty(123).payload()
////		if err != nil {
////			b.Fail()
////		}
////	}
////}
//
//func TestClient_GetOrdersRequest(t *testing.T) {
//
//	//r := ReqToGetOrders{
//	//	Symbol:    "",
//	//	Filter:    nil,
//	//	Columns:   "",
//	//	Count:     0,
//	//	Start:     0,
//	//	Reverse:   nil,
//	//	StartTime: nil,
//	//	EndTime:   nil,
//	//}
//	//k, _ := json.Marshal(r)
//	//fmt.Println(string(k))
//	//
//	//return
//
//	b := NewBitmex(true)
//	c := b.NewRestClient(auth.NewConfig("YXTfYOwQn0e7gq3bziPyeNDa", "6tbsbjU1GjCGipthTHuLb2PrVswCuGBrA10JQ5IrhlGy0Ytz"))
//	c.httpC.MaxIdleConnDuration = time.Second * 150
//
//	resp, err := c.GetOrders(context.Background(), ReqToGetOrders{})
//
//	fmt.Println(resp, err)
//
//	//for i := 0; i < 60; i++ {
//	//	_, err := c.GetOrdersRequest().Symbol("XBTUSD").Filter(map[string]interface{}{"open": true}).Do()
//	//	if err != nil {
//	//		t.Error(err)
//	//	}
//	//}
//
//	//startTime := time.Now().UnixNano()
//	//_, _ = c.GetOrdersRequest().Symbol("XBTUSD").Filter(map[string]interface{}{"open": true}).Do()
//	//fmt.Println("Total Time: ", time.Now().UnixNano()-startTime)
//
//	fmt.Println("**********************************************")
//
//	//time.Sleep(time.Second*10)
//	//
//	//startTime = time.Now().UnixNano()
//	//_, _ = c.GetOrdersRequest().Symbol("XBTUSD").Filter(map[string]interface{}{"open": true}).Do()
//	//fmt.Println("Total Time: ", time.Now().UnixNano()-startTime)
//	//
//	//fmt.Println("**********************************************")
//	//
//	//
//	//
//	//time.Sleep(time.Second*180)
//	//
//	//startTime = time.Now().UnixNano()
//	//_, _ = c.GetOrdersRequest().Symbol("XBTUSD").Filter(map[string]interface{}{"open": true}).Do()
//	//fmt.Println("Total Time: ", time.Now().UnixNano()-startTime)
//
//	//fmt.Println(do)
//}
//
////
////func TestClient_PlaceBulkOrdersRequest(t *testing.T) {
////	b := NewBitmex(true)
////	c := b.NewRestClient(auth.NewConfig("pkOAYJNujj_-cW3gN0FjPizp", "LgiVEp09S4TFw9FfoMuurP-7lVJ6DPfbCRHWhomkqn7qx-F4"))
////
////	resp, err := c.GetOrders(context.Background(), ReqToGetOrders{})
////	fmt.Println(resp, err)
////
////	//for i := 0; i < 60; i++ {
////	//
////	//}
////	//fmt.Println(do)
////}
////
////func TestClient_GetInstrumentsRequest(t *testing.T) {
////	c := NewRestClient(auth.NewConfig("pkOAYJNujj_-cW3gN0FjPizp", "LgiVEp09S4TFw9FfoMuurP-7lVJ6DPfbCRHWhomkqn7qx-F4"), true)
////
////	_, err := c.GetInstruments().Symbol("XBTUSD").Do()
////	if err != nil {
////		t.Error(err)
////	}
////
////	//for i := 0; i < 60; i++ {
////	//
////	//}
////	//fmt.Println(do)
////}
////
////func TestClient_AmendOrderRequest(t *testing.T) {
////
////}
////
////func TestClient_PlaceOrderRequest(t *testing.T) {
////
////}
////
////func TestClient_CancelOrdersRequest(t *testing.T) {
////
////}
////
////func TestClient_CancelAllOrdersRequest(t *testing.T) {
////
////}
////
////func TestClient_AmendBulkOrdersRequest(t *testing.T) {
////
////}
////
////func TestClient_CancelAllAfterRequest(t *testing.T) {
////
////}
////
////func TestClient_GetOrdersRequestQuery(t *testing.T) {
////
////}
////
////func TestClient_AmendOrderRequestPayload(t *testing.T) {
////	c := restClient{}
////
////	payload, err := c.AmendOrderRequest().OrderID("FF").OrderQty(100).Price(0).payload()
////	if err != nil {
////		t.Errorf("payload generate error")
////	}
////	expected := `{"orderID":"FF","orderQty":100,"price":0}`
////
////	if payload != expected {
////		t.Fatalf("Incorrect Payload")
////	}
////}
////
////func TestClient_PlaceOrderRequestPayload(t *testing.T) {
////	c := restClient{}
////	payload, err := c.PlaceOrderRequest().OrderQty(100).DisplayQty(0).Symbol("XBTUSD").payload()
////	if err != nil {
////		t.Errorf("payload generate error")
////	}
////	expected := `{"symbol":"XBTUSD","orderQty":100,"displayQty":0}`
////
////	if payload != expected {
////		t.Fatalf("incorret payload")
////	}
////
////}
////
////func TestClient_CancelOrdersRequestPayload(t *testing.T) {
////	c := restClient{}
////	payload, err := c.CancelOrdersRequest().AddOrderIDs("abc", "def").payload()
////	if err != nil {
////		t.Errorf("payload generate error")
////	}
////	expected := `{"orderID":["abc","def"]}`
////
////	if payload != expected {
////		t.Fatalf("incorrect payload")
////	}
////}
////
////func TestClient_CancelAllOrdersRequestPayload(t *testing.T) {
////	c := restClient{}
////	payload, err := c.CancelAllOrdersRequest().Symbol("XBTUSD").Filter(map[string]interface{}{"side": "Buy"}).payload()
////	if err != nil {
////		t.Errorf("payload generate error")
////	}
////	expected := `{"symbol":"XBTUSD","filter":{"side":"Buy"}}`
////
////	if payload != expected {
////		t.Fatalf("incorrect payload")
////	}
////}
////
////func TestClient_AmendBulkOrdersRequestPayload(t *testing.T) {
////	c := restClient{}
////
////	payload, err := c.AmendBulkOrdersRequest().AddAmendedOrder(
////		c.AmendOrderRequest().OrderID("abc").Price(0),
////		c.AmendOrderRequest().ClOrdID("def").OrderQty(0).Price(100)).payload()
////	if err != nil {
////		t.Errorf("payload generate error")
////	}
////
////	expected := `{"orders":["{\"orderID\":\"abc\",\"price\":0}","{\"clOrdID\":\"def\",\"orderQty\":0,\"price\":100}"]}`
////
////	if payload != expected {
////		t.Fatalf("incorrect payload")
////	}
////}
////
////func TestClient_PlaceBulkOrdersRequestPayload(t *testing.T) {
////
////}
////
////func TestClient_CancelAllAfterRequestPayload(t *testing.T) {
////	c := restClient{}
////
////	payload, err := c.CancelAllAfterRequest().Timeout(0).payload()
////	if err != nil {
////		t.Errorf("payload generate error")
////	}
////	expected := `{"timeout":0}`
////
////	if payload != expected {
////		t.Fatalf("incorrect payload")
////	}
////}
