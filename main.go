package main

import (
	"fmt"
	"github.com/bitmex-mirror/position"
	"github.com/bitmex-mirror/rest"
	"os"
)

func main() {
	r := rest.NewRestObject()

	c := r.AddNewClient("pkOAYJNujj_-cW3gN0FjPizp", "LgiVEp09S4TFw9FfoMuurP-7lVJ6DPfbCRHWhomkqn7qx-F4", true)

	//c.CancelAllAfter(&order.ReqToCancelAllAfter{Timeout: 60000})

	//allOrders, err := c.CancelAllOrders(&order.ReqToCancelAllOrders{
	//	Symbol: "XBTUSD",
	//	Filter: map[string]interface{}{"side": "Buy"},
	//})
	//if err != nil {
	//	panic(err)
	//}

	leverage, err := c.ChangeLeverage(&position.ReqToChangeLeverage{
		Symbol:   "XBTUSD",
		Leverage: 0,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(leverage)
	os.Exit(0)

	//c.GetOrders(&order.ReqToGetOrders{
	//	Symbol:  "XBTUSD",
	//	Reverse: true,
	//	Count:   2,
	//})

	//c.AmendOrder(&order.ReqToAmendOrder{
	//	OrderID:        "710ddeca-c147-45ac-bc5b-fa00dbb3467c",
	//	Price:          36000,
	//})
	//
	//o1 := order.ReqToAmendOrder{
	//	OrderID: "710ddeca-c147-45ac-bc5b-fa00dbb3467c",
	//	Price:   32000,
	//}
	//
	//o2 := order.ReqToAmendOrder{
	//	OrderID: "16146363-6597-4ea3-b699-0db3f7250d49",
	//	Price:   37000,
	//}
	//
	//var ab order.ReqToAmendBulkOrders
	//ab.Orders = append(ab.Orders, o1, o2)
	//
	//c.AmendBulkOrders(&ab)

	//c.CancelOrders(&order.ReqToCancelOrders{
	//	OrderIDs: []string{"4ba1ee5f-a70e-44b8-a51a-772993160b6a", "b8ec7535-3f04-40ab-a301-253a89259a42"},
	//})

	//var bulkOrder order.ReqToPlaceBulkOrder
	//
	//ord1 := order.ReqToPlaceOrder{
	//	Symbol:   "XBTUSD",
	//	Side:     "Buy",
	//	OrderQty: 100,
	//	Price:    30000,
	//}
	//
	//ord2 := order.ReqToPlaceOrder{
	//	Symbol:   "XBTUSD",
	//	Side:     "Buy",
	//	OrderQty: 100,
	//	Price:    35000,
	//}
	//
	//bulkOrder.Orders = append(bulkOrder.Orders, ord1, ord2)
	//
	//orders, err := c.PlaceBuckOrders(&bulkOrder)
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("##########################################")
	fmt.Println("##########################################")
	fmt.Println("##########################################")
	fmt.Println("##########################################")
	fmt.Println("##########################################")

	//fmt.Println(allOrders)

}

//func main() {
//	apiKey := "pkOAYJNujj_-cW3gN0FjPizp"
//	apiSecret := "LgiVEp09S4TFw9FfoMuurP-7lVJ6DPfbCRHWhomkqn7qx-F4"
//
//	auth := bitmexgo.NewAPIKeyContext(apiKey, apiSecret)
//
//	// Create a testnet API client instance
//	testnetClient := bitmexgo.NewAPIClient(bitmexgo.NewTestnetConfiguration())
//
//	var params bitmexgo.OrderGetOrdersOpts
//
//	params.Symbol.Set("XBTUSD")
//	params.Reverse.Set(true)
//	params.Count.Set(1)
//
//	orders, _, err := testnetClient.OrderApi.OrderGetOrders(auth, &params)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(orders)
//
//}
