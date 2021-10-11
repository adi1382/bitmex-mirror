package client

import (
	"encoding/json"
	"fmt"
	"github.com/bitmex-mirror/optional"
	"github.com/google/go-querystring/query"
	"net/http"
	"strings"
	"time"
)

func (c *Client) GetOrdersRequest() *reqToGetOrders {
	return &reqToGetOrders{
		c: c,
	}
}

func (c *Client) AmendOrderRequest() *reqToAmendOrder {
	return &reqToAmendOrder{
		c: c,
	}
}

func (c *Client) PlaceOrderRequest() *reqToPlaceOrder {
	return &reqToPlaceOrder{
		c: c,
	}
}

func (c *Client) CancelOrdersRequest() *reqToCancelOrders {
	return &reqToCancelOrders{
		c: c,
	}
}

func (c *Client) CancelAllOrdersRequest() *reqToCancelAllOrders {
	return &reqToCancelAllOrders{
		c:      c,
		filter: make(map[string]interface{}),
	}
}

func (c *Client) AmendBulkOrdersRequest() *reqToAmendBulkOrders {
	return &reqToAmendBulkOrders{
		c: c,
	}
}

func (c *Client) PlaceBulkOrdersRequest() *reqToPlaceBulkOrders {
	return &reqToPlaceBulkOrders{
		c: c,
	}
}

func (c *Client) CancelAllAfterRequest() *reqToCancelAllAfter {
	return &reqToCancelAllAfter{
		c: c,
	}
}

type Order struct {
	OrderID               string    `json:"orderID"`
	ClOrdID               string    `json:"clOrdID"`
	ClOrdLinkID           string    `json:"clOrdLinkID"`
	Account               int       `json:"account"`
	Symbol                string    `json:"symbol"`
	Side                  string    `json:"side"`
	SimpleOrderQty        int       `json:"simpleOrderQty"`
	OrderQty              int       `json:"orderQty"`
	Price                 int       `json:"price"`
	DisplayQty            int       `json:"displayQty"`
	StopPx                int       `json:"stopPx"`
	PegOffsetValue        int       `json:"pegOffsetValue"`
	PegPriceType          string    `json:"pegPriceType"`
	Currency              string    `json:"currency"`
	SettlCurrency         string    `json:"settlCurrency"`
	OrdType               string    `json:"ordType"`
	TimeInForce           string    `json:"timeInForce"`
	ExecInst              string    `json:"execInst"`
	ContingencyType       string    `json:"contingencyType"`
	ExDestination         string    `json:"exDestination"`
	OrdStatus             string    `json:"ordStatus"`
	Triggered             string    `json:"triggered"`
	WorkingIndicator      bool      `json:"workingIndicator"`
	OrdRejReason          string    `json:"ordRejReason"`
	SimpleLeavesQty       int       `json:"simpleLeavesQty"`
	LeavesQty             int       `json:"leavesQty"`
	SimpleCumQty          int       `json:"simpleCumQty"`
	CumQty                int       `json:"cumQty"`
	AvgPx                 int       `json:"avgPx"`
	MultiLegReportingType string    `json:"multiLegReportingType"`
	Text                  string    `json:"text"`
	TransactTime          time.Time `json:"transactTime"`
	Timestamp             time.Time `json:"timestamp"`
}

type reqToGetOrders struct {
	c         *Client
	symbol    string
	filter    map[string]interface{}
	columns   string
	count     int
	start     int
	reverse   optional.Bool
	startTime time.Time
	endTime   time.Time
}

func (req *reqToGetOrders) Symbol(symbol string) *reqToGetOrders {
	req.symbol = symbol
	return req
}

func (req *reqToGetOrders) Filter(filter map[string]interface{}) *reqToGetOrders {
	req.filter = filter
	return req
}

func (req *reqToGetOrders) Columns(columns ...string) *reqToGetOrders {
	req.columns = strings.Join(columns, ",")
	return req
}

func (req *reqToGetOrders) Count(count int) *reqToGetOrders {
	req.count = count
	return req
}

func (req *reqToGetOrders) Start(start int) *reqToGetOrders {
	req.start = start
	return req
}

func (req *reqToGetOrders) Reverse(reverse bool) *reqToGetOrders {
	req.reverse.Set(reverse)
	return req
}

func (req *reqToGetOrders) StartTime(startTime time.Time) *reqToGetOrders {
	req.startTime = startTime
	return req
}

func (req *reqToGetOrders) EndTime(endTime time.Time) *reqToGetOrders {
	req.endTime = endTime
	return req
}

func (req *reqToGetOrders) Do() (RespToGetOrders, error) {
	var response RespToGetOrders
	err := req.c.request(req, &response)
	return response, err
}

type RespToGetOrders []Order

func (req *reqToGetOrders) path() string {
	return fmt.Sprintf("/order")
}

func (req *reqToGetOrders) method() string {
	return http.MethodGet
}

func (req *reqToGetOrders) query() string {
	s, e := json.Marshal(req.filter)
	filterStr := string(s)
	if filterStr == "null" || e != nil {
		filterStr = ""
	}

	s, e = json.Marshal(req.reverse)
	reverseStr := string(s)
	if reverseStr == "null" || e != nil {
		reverseStr = ""
	}

	a := struct {
		C         *Client   `url:"-"`
		Symbol    string    `url:"symbol,omitempty"`
		Filter    string    `url:"filter,omitempty"`
		Columns   string    `url:"columns,omitempty"`
		Count     int       `url:"count,omitempty"`
		Start     int       `url:"start,omitempty"`
		Reverse   string    `url:"reverse,omitempty"`
		StartTime time.Time `url:"startTime,omitempty"`
		EndTime   time.Time `url:"endTime,omitempty"`
	}{
		req.c,
		req.symbol,
		filterStr,
		req.columns,
		req.count,
		req.start,
		reverseStr,
		req.startTime,
		req.endTime,
	}
	value, _ := query.Values(&a)

	if filterStr != "" {
		return value.Encode()
	}

	return value.Encode()
}

func (req *reqToGetOrders) payload() string {
	return ""
}

// reqToAmendOrder --

type reqToAmendOrder struct {
	c              *Client
	orderID        string
	origClOrdID    string
	clOrdID        string
	orderQty       *float64
	leavesQty      *float64
	price          *float64
	stopPx         *float64
	pegOffsetValue *float64
	text           string
}

func (req *reqToAmendOrder) OrderID(orderID string) *reqToAmendOrder {
	req.orderID = orderID
	return req
}

func (req *reqToAmendOrder) OrigClOrdID(origClOrdID string) *reqToAmendOrder {
	req.origClOrdID = origClOrdID
	return req
}

func (req *reqToAmendOrder) ClOrdID(clOrdID string) *reqToAmendOrder {
	req.clOrdID = clOrdID
	return req
}

func (req *reqToAmendOrder) OrderQty(orderQty float64) *reqToAmendOrder {
	req.orderQty = &orderQty
	return req
}

func (req *reqToAmendOrder) LeavesQty(leavesQty float64) *reqToAmendOrder {
	req.leavesQty = &leavesQty
	return req
}

func (req *reqToAmendOrder) Price(price float64) *reqToAmendOrder {
	req.price = &price
	return req
}

func (req *reqToAmendOrder) StopPx(stopPx float64) *reqToAmendOrder {
	req.stopPx = &stopPx
	return req
}

func (req *reqToAmendOrder) PegOffsetValue(pegOffsetValue float64) *reqToAmendOrder {
	req.pegOffsetValue = &pegOffsetValue
	return req
}

func (req *reqToAmendOrder) Text(text string) *reqToAmendOrder {
	req.text = text
	return req
}

func (req *reqToAmendOrder) Do() (RespToAmendOrder, error) {
	response := RespToAmendOrder{}
	fmt.Println("rr")
	fmt.Println(response)
	err := req.c.request(req, &response)
	return response, err
}

type RespToAmendOrder Order

func (req *reqToAmendOrder) path() string {
	return fmt.Sprintf("/order")
}

func (req *reqToAmendOrder) method() string {
	return http.MethodPut
}

func (req *reqToAmendOrder) query() string {
	return ""
}

func (req *reqToAmendOrder) payload() string {

	a := struct {
		C              *Client  `json:"-"`
		OrderID        string   `json:"orderID,omitempty"`
		OrigClOrdID    string   `json:"origClOrdID,omitempty"`
		ClOrdID        string   `json:"clOrdID,omitempty"`
		OrderQty       *float64 `json:"orderQty,omitempty"`
		LeavesQty      *float64 `json:"leavesQty,omitempty"`
		Price          *float64 `json:"price,omitempty"`
		StopPx         *float64 `json:"stopPx,omitempty"`
		PegOffsetValue *float64 `json:"pefOffsetValue,omitempty"`
		Text           string   `json:"text,omitempty"`
	}{
		req.c,
		req.orderID,
		req.origClOrdID,
		req.clOrdID,
		req.orderQty,
		req.leavesQty,
		req.price,
		req.stopPx,
		req.pegOffsetValue,
		req.text,
	}

	b, err := json.Marshal(&a)
	if err != nil {
		return ""
	}
	return string(b)
}

// GetOrdersRequest reqToPlaceOrder  --

type reqToPlaceOrder struct {
	c              *Client
	symbol         string
	side           string
	orderQty       *float64
	price          *float64
	displayQty     *float64
	stopPx         *float64
	clOrdID        string
	pegOffsetValue *float64
	pegPriceType   string
	ordTyp         string
	timeInForce    string
	execInst       string
	text           string
}

func (req *reqToPlaceOrder) Symbol(symbol string) *reqToPlaceOrder {
	req.symbol = symbol
	return req
}

func (req *reqToPlaceOrder) Side(side string) *reqToPlaceOrder {
	req.side = side
	return req
}

func (req *reqToPlaceOrder) OrderQty(orderQty float64) *reqToPlaceOrder {
	req.orderQty = &orderQty
	return req
}

func (req *reqToPlaceOrder) Price(price float64) *reqToPlaceOrder {
	req.price = &price
	return req
}

func (req *reqToPlaceOrder) DisplayQty(displayQty float64) *reqToPlaceOrder {
	req.displayQty = &displayQty
	return req
}

func (req *reqToPlaceOrder) StopPx(stopPx float64) *reqToPlaceOrder {
	req.stopPx = &stopPx
	return req
}

func (req *reqToPlaceOrder) ClOrdID(clOrdID string) *reqToPlaceOrder {
	req.clOrdID = clOrdID
	return req
}

func (req *reqToPlaceOrder) PegOffsetValue(pegOffsetValue float64) *reqToPlaceOrder {
	req.pegOffsetValue = &pegOffsetValue
	return req
}

func (req *reqToPlaceOrder) PegPriceType(pegPriceType string) *reqToPlaceOrder {
	req.pegPriceType = pegPriceType
	return req
}

func (req *reqToPlaceOrder) OrdTyp(ordTyp string) *reqToPlaceOrder {
	req.ordTyp = ordTyp
	return req
}

func (req *reqToPlaceOrder) TimeInForce(timeInForce string) *reqToPlaceOrder {
	req.timeInForce = timeInForce
	return req
}

func (req *reqToPlaceOrder) ExecInst(execInst string) *reqToPlaceOrder {
	req.execInst = execInst
	return req
}

func (req *reqToPlaceOrder) Text(text string) *reqToPlaceOrder {
	req.text = text
	return req
}

func (req *reqToPlaceOrder) Do() (RespToPlaceOrders, error) {
	response := RespToPlaceOrders{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToPlaceOrders Order

func (req *reqToPlaceOrder) path() string {
	return fmt.Sprintf("/order")
}

func (req *reqToPlaceOrder) method() string {
	return http.MethodPost
}

func (req *reqToPlaceOrder) query() string {
	return ""
}

func (req *reqToPlaceOrder) payload() string {

	a := struct {
		C              *Client  `json:"-"`
		Symbol         string   `json:"symbol,omitempty"`
		Side           string   `json:"side,omitempty"`
		OrderQty       *float64 `json:"orderQty,omitempty"`
		Price          *float64 `json:"price,omitempty"`
		DisplayQty     *float64 `json:"displayQty,omitempty"`
		StopPx         *float64 `json:"stopPx,omitempty"`
		ClOrdID        string   `json:"clOrdID,omitempty"`
		PegOffsetValue *float64 `json:"pefOffsetValue,omitempty"`
		PegPriceType   string   `json:"pegPriceType,omitempty"`
		OrdTyp         string   `json:"ordType,omitempty"`
		TimeInForce    string   `json:"timeInForce,omitempty"`
		ExecInst       string   `json:"execInst,omitempty"`
		Text           string   `json:"text,omitempty"`
	}{
		req.c,
		req.symbol,
		req.side,
		req.orderQty,
		req.price,
		req.displayQty,
		req.stopPx,
		req.clOrdID,
		req.pegOffsetValue,
		req.pegPriceType,
		req.ordTyp,
		req.timeInForce,
		req.execInst,
		req.text,
	}

	b, err := json.Marshal(&a)
	if err != nil {
		return ""
	}
	return string(b)
}

// reqToCancelOrders --

type reqToCancelOrders struct {
	c        *Client
	orderIDs []string
	clOrdIDs []string
	text     string
}

func (req *reqToCancelOrders) AddOrderIDs(orderIDs ...string) *reqToCancelOrders {
	req.orderIDs = append(req.orderIDs, orderIDs...)
	return req
}

func (req *reqToCancelOrders) AddClOrdIDs(clOrdIDs ...string) *reqToCancelOrders {
	req.clOrdIDs = append(req.clOrdIDs, clOrdIDs...)
	return req
}

func (req *reqToCancelOrders) Text(text string) *reqToCancelOrders {
	req.text = text
	return req
}

func (req *reqToCancelOrders) Do() (RespToCancelOrders, error) {
	response := RespToCancelOrders{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToCancelOrders Order

func (req *reqToCancelOrders) path() string {
	return fmt.Sprintf("/order")
}

func (req *reqToCancelOrders) method() string {
	return http.MethodDelete
}

func (req *reqToCancelOrders) query() string {
	return ""
}

func (req *reqToCancelOrders) payload() string {

	a := struct {
		C        *Client  `json:"-"`
		OrderIDs []string `json:"orderID,omitempty"`
		ClOrdIDs []string `json:"clOrdID,omitempty"`
		Text     string   `json:"text,omitempty"`
	}{
		req.c,
		req.orderIDs,
		req.clOrdIDs,
		req.text,
	}

	b, err := json.Marshal(&a)
	if err != nil {
		return ""
	}
	return string(b)
}

// reqToCancelAllOrders --

type reqToCancelAllOrders struct {
	c      *Client
	symbol string
	filter map[string]interface{}
	text   string
}

func (req *reqToCancelAllOrders) Symbol(symbol string) *reqToCancelAllOrders {
	req.symbol = symbol
	return req
}

func (req *reqToCancelAllOrders) AddFilter(key string, value interface{}) *reqToCancelAllOrders {
	req.filter[key] = value
	return req
}

func (req *reqToCancelAllOrders) Text(text string) *reqToCancelAllOrders {
	req.text = text
	return req
}

func (req *reqToCancelAllOrders) Do() (RespToCancelAllOrders, error) {
	response := RespToCancelAllOrders{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToCancelAllOrders []Order

func (req *reqToCancelAllOrders) path() string {
	return fmt.Sprintf("/order/all")
}

func (req *reqToCancelAllOrders) method() string {
	return http.MethodDelete
}

func (req *reqToCancelAllOrders) query() string {
	return ""
}

func (req *reqToCancelAllOrders) payload() string {

	a := struct {
		C      *Client                `json:"-"`
		Symbol string                 `json:"symbol,omitempty"`
		Filter map[string]interface{} `json:"filter,omitempty"`
		Text   string                 `json:"text,omitempty"`
	}{
		req.c,
		req.symbol,
		req.filter,
		req.text,
	}

	b, err := json.Marshal(&a)
	if err != nil {
		return ""
	}
	return string(b)
}

// reqToAmendBulkOrders --
type reqToAmendBulkOrders struct {
	c      *Client
	orders []*reqToAmendOrder
}

func (req *reqToAmendBulkOrders) AddAmendedOrder(order ...*reqToAmendOrder) *reqToAmendBulkOrders {
	req.orders = append(req.orders, order...)
	return req
}

func (req *reqToAmendBulkOrders) Do() (RespToAmendBulkOrders, error) {
	response := RespToAmendBulkOrders{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToAmendBulkOrders []Order

func (req *reqToAmendBulkOrders) path() string {
	return fmt.Sprintf("/order/bulk")
}

func (req *reqToAmendBulkOrders) method() string {
	return http.MethodPut
}

func (req *reqToAmendBulkOrders) query() string {
	return ""
}

func (req *reqToAmendBulkOrders) payload() string {

	a := make([]string, 0, 2)

	for i := range req.orders {
		a = append(a, req.orders[i].payload())
	}

	c := struct {
		Orders []string `json:"orders,omitempty"`
	}{
		a,
	}

	b, err := json.Marshal(&c)
	if err != nil {
		return ""
	}
	return string(b)
}

// reqToPlaceBulkOrders --
type reqToPlaceBulkOrders struct {
	c      *Client
	orders []*reqToPlaceOrder
}

func (req *reqToPlaceBulkOrders) AddOrders(orders ...*reqToPlaceOrder) *reqToPlaceBulkOrders {
	req.orders = append(req.orders, orders...)
	return req
}

func (req *reqToPlaceBulkOrders) Do() (RespToPlaceBulkOrders, error) {
	response := RespToPlaceBulkOrders{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToPlaceBulkOrders []Order

func (req *reqToPlaceBulkOrders) path() string {
	return fmt.Sprintf("/order/bulk")
}

func (req *reqToPlaceBulkOrders) method() string {
	return http.MethodPost
}

func (req *reqToPlaceBulkOrders) query() string {
	return ""
}

func (req *reqToPlaceBulkOrders) payload() string {

	a := make([]string, 0, 2)

	for i := range req.orders {
		a = append(a, req.orders[i].payload())
	}

	c := struct {
		Orders []string `json:"orders,omitempty"`
	}{
		a,
	}

	b, err := json.Marshal(&c)
	if err != nil {
		return ""
	}
	return string(b)
}

// reqToCancelAllAfter --
type reqToCancelAllAfter struct {
	c       *Client
	timeout *float64
}

func (req *reqToCancelAllAfter) Timeout(timeout float64) *reqToCancelAllAfter {
	req.timeout = &timeout
	return req
}

func (req *reqToCancelAllAfter) Do() (RespToCancelAllAfter, error) {
	response := RespToCancelAllAfter{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToCancelAllAfter struct{}

func (req *reqToCancelAllAfter) path() string {
	return fmt.Sprintf("/order/cancelAllAfter")
}

func (req *reqToCancelAllAfter) method() string {
	return http.MethodPost
}

func (req *reqToCancelAllAfter) query() string {
	return ""
}

func (req *reqToCancelAllAfter) payload() string {

	a := struct {
		C       *Client  `json:"-"`
		Timeout *float64 `json:"timeout,omitempty"`
	}{
		req.c,
		req.timeout,
	}

	b, err := json.Marshal(&a)
	if err != nil {
		return ""
	}
	return string(b)
}
