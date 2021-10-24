package bitmex

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"4d63.com/optional"
	"github.com/google/go-querystring/query"
	"github.com/valyala/fasthttp"
)

func (c *restClient) GetOrdersRequest() *reqToGetOrders {
	return &reqToGetOrders{
		c: c,
	}
}

func (c *restClient) AmendOrderRequest() *reqToAmendOrder {
	return &reqToAmendOrder{
		c: c,
	}
}

func (c *restClient) PlaceOrderRequest() *reqToPlaceOrder {
	return &reqToPlaceOrder{
		c: c,
	}
}

func (c *restClient) CancelOrdersRequest() *reqToCancelOrders {
	return &reqToCancelOrders{
		c: c,
	}
}

func (c *restClient) CancelAllOrdersRequest() *reqToCancelAllOrders {
	return &reqToCancelAllOrders{
		c: c,
	}
}

func (c *restClient) AmendBulkOrdersRequest() *reqToAmendBulkOrders {
	return &reqToAmendBulkOrders{
		c: c,
	}
}

func (c *restClient) PlaceBulkOrdersRequest() *reqToPlaceBulkOrders {
	return &reqToPlaceBulkOrders{
		c: c,
	}
}

func (c *restClient) CancelAllAfterRequest() *reqToCancelAllAfter {
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

// reqToGetOrders --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToGetOrders struct {
	c         *restClient
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
	req.reverse = optional.Bool{reverse}
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
	req.c.bucket1m.Wait(1)
	var response RespToGetOrders
	err := req.c.request(req, &response)
	return response, err
}

type RespToGetOrders []Order

func (req *reqToGetOrders) path() string {
	return fmt.Sprintf("/order")
}

func (req *reqToGetOrders) method() string {
	return fasthttp.MethodGet
}

func (req *reqToGetOrders) query() (string, error) {
	s, e := json.Marshal(req.filter)
	filterStr := string(s)
	if filterStr == "null" || e != nil {
		filterStr = ""
	}

	a := struct {
		C         *restClient   `url:"-"`
		Symbol    string        `url:"symbol,omitempty"`
		Filter    string        `url:"filter,omitempty"`
		Columns   string        `url:"columns,omitempty"`
		Count     int           `url:"count,omitempty"`
		Start     int           `url:"start,omitempty"`
		Reverse   optional.Bool `url:"reverse,omitempty"`
		StartTime time.Time     `url:"startTime,omitempty"`
		EndTime   time.Time     `url:"endTime,omitempty"`
	}{
		req.c,
		req.symbol,
		filterStr,
		req.columns,
		req.count,
		req.start,
		req.reverse,
		req.startTime,
		req.endTime,
	}
	value, err := query.Values(&a)

	if err != nil {
		return "", err
	}

	return value.Encode(), nil
}

func (req *reqToGetOrders) payload() (string, error) {
	return "", nil
}

// reqToAmendOrder --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToAmendOrder struct {
	c              *restClient
	orderID        string
	origClOrdID    string
	clOrdID        string
	orderQty       optional.Float64
	leavesQty      optional.Float64
	price          optional.Float64
	stopPx         optional.Float64
	pegOffsetValue optional.Float64
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
	req.orderQty = optional.Float64{orderQty}
	return req
}

func (req *reqToAmendOrder) LeavesQty(leavesQty float64) *reqToAmendOrder {
	req.leavesQty = optional.OfFloat64(leavesQty)
	return req
}

func (req *reqToAmendOrder) Price(price float64) *reqToAmendOrder {
	req.price = optional.OfFloat64(price)
	return req
}

func (req *reqToAmendOrder) StopPx(stopPx float64) *reqToAmendOrder {
	req.stopPx = optional.OfFloat64(stopPx)
	return req
}

func (req *reqToAmendOrder) PegOffsetValue(pegOffsetValue float64) *reqToAmendOrder {
	req.pegOffsetValue = optional.OfFloat64(pegOffsetValue)
	return req
}

func (req *reqToAmendOrder) Text(text string) *reqToAmendOrder {
	req.text = text
	return req
}

func (req *reqToAmendOrder) Do() (RespToAmendOrder, error) {
	response := RespToAmendOrder{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToAmendOrder Order

func (req *reqToAmendOrder) path() string {
	return fmt.Sprintf("/order")
}

func (req *reqToAmendOrder) method() string {
	return fasthttp.MethodPut
}

func (req *reqToAmendOrder) query() (string, error) {
	return "", nil
}

func (req *reqToAmendOrder) payload() (string, error) {
	a := struct {
		C              *restClient      `json:"-"`
		OrderID        string           `json:"orderID,omitempty"`
		OrigClOrdID    string           `json:"origClOrdID,omitempty"`
		ClOrdID        string           `json:"clOrdID,omitempty"`
		OrderQty       optional.Float64 `json:"orderQty,omitempty"`
		LeavesQty      optional.Float64 `json:"leavesQty,omitempty"`
		Price          optional.Float64 `json:"price,omitempty"`
		StopPx         optional.Float64 `json:"stopPx,omitempty"`
		PegOffsetValue optional.Float64 `json:"pefOffsetValue,omitempty"`
		Text           string           `json:"text,omitempty"`
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
	return string(b), err
}

// reqToPlaceOrder  --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToPlaceOrder struct {
	c              *restClient
	symbol         string
	side           string
	orderQty       optional.Float64
	price          optional.Float64
	displayQty     optional.Float64
	stopPx         optional.Float64
	clOrdID        string
	pegOffsetValue optional.Float64
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
	req.orderQty = optional.OfFloat64(orderQty)
	return req
}

func (req *reqToPlaceOrder) Price(price float64) *reqToPlaceOrder {
	req.price = optional.OfFloat64(price)
	return req
}

func (req *reqToPlaceOrder) DisplayQty(displayQty float64) *reqToPlaceOrder {
	req.displayQty = optional.OfFloat64(displayQty)
	return req
}

func (req *reqToPlaceOrder) StopPx(stopPx float64) *reqToPlaceOrder {
	req.stopPx = optional.OfFloat64(stopPx)
	return req
}

func (req *reqToPlaceOrder) ClOrdID(clOrdID string) *reqToPlaceOrder {
	req.clOrdID = clOrdID
	return req
}

func (req *reqToPlaceOrder) PegOffsetValue(pegOffsetValue float64) *reqToPlaceOrder {
	req.pegOffsetValue = optional.OfFloat64(pegOffsetValue)
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
	return fasthttp.MethodPost
}

func (req *reqToPlaceOrder) query() (string, error) {
	return "", nil
}

func (req *reqToPlaceOrder) payload() (string, error) {

	a := struct {
		C              *restClient      `json:"-"`
		Symbol         string           `json:"symbol,omitempty"`
		Side           string           `json:"side,omitempty"`
		OrderQty       optional.Float64 `json:"orderQty,omitempty"`
		Price          optional.Float64 `json:"price,omitempty"`
		DisplayQty     optional.Float64 `json:"displayQty,omitempty"`
		StopPx         optional.Float64 `json:"stopPx,omitempty"`
		ClOrdID        string           `json:"clOrdID,omitempty"`
		PegOffsetValue optional.Float64 `json:"pefOffsetValue,omitempty"`
		PegPriceType   string           `json:"pegPriceType,omitempty"`
		OrdTyp         string           `json:"ordType,omitempty"`
		TimeInForce    string           `json:"timeInForce,omitempty"`
		ExecInst       string           `json:"execInst,omitempty"`
		Text           string           `json:"text,omitempty"`
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
	return string(b), err
}

// reqToCancelOrders --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToCancelOrders struct {
	c        *restClient
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
	return fasthttp.MethodDelete
}

func (req *reqToCancelOrders) query() (string, error) {
	return "", nil
}

func (req *reqToCancelOrders) payload() (string, error) {

	a := struct {
		C        *restClient `json:"-"`
		OrderIDs []string    `json:"orderID,omitempty"`
		ClOrdIDs []string    `json:"clOrdID,omitempty"`
		Text     string      `json:"text,omitempty"`
	}{
		req.c,
		req.orderIDs,
		req.clOrdIDs,
		req.text,
	}

	b, err := json.Marshal(&a)
	return string(b), err
}

// reqToCancelAllOrders --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToCancelAllOrders struct {
	c      *restClient
	symbol string
	filter map[string]interface{}
	text   string
}

func (req *reqToCancelAllOrders) Symbol(symbol string) *reqToCancelAllOrders {
	req.symbol = symbol
	return req
}

func (req *reqToCancelAllOrders) Filter(filter map[string]interface{}) *reqToCancelAllOrders {
	req.filter = filter
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
	return fasthttp.MethodDelete
}

func (req *reqToCancelAllOrders) query() (string, error) {
	return "", nil
}

func (req *reqToCancelAllOrders) payload() (string, error) {

	a := struct {
		C      *restClient            `json:"-"`
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
	return string(b), err
}

// reqToAmendBulkOrders --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToAmendBulkOrders struct {
	c      *restClient
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
	return fasthttp.MethodPut
}

func (req *reqToAmendBulkOrders) query() (string, error) {
	return "", nil
}

func (req *reqToAmendBulkOrders) payload() (string, error) {

	a := make([]string, 0, 2)

	for i := range req.orders {
		if payload, err := req.orders[i].payload(); err != nil {
			return "", err
		} else {
			a = append(a, payload)
		}
	}

	c := struct {
		Orders []string `json:"orders,omitempty"`
	}{
		a,
	}

	b, err := json.Marshal(&c)
	return string(b), err
}

// reqToPlaceBulkOrders --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToPlaceBulkOrders struct {
	c      *restClient
	orders []*reqToPlaceOrder
}

func (req *reqToPlaceBulkOrders) AddOrders(orders ...*reqToPlaceOrder) *reqToPlaceBulkOrders {
	req.orders = append(req.orders, orders...)
	return req
}

func (req *reqToPlaceBulkOrders) Do() (RespToPlaceBulkOrders, error) {
	req.c.bucket1m.Wait(int64(math.Max(math.Ceil(float64(len(req.orders))/10), 1)))
	req.c.bucket10s.Wait(1)
	response := RespToPlaceBulkOrders{}
	err := req.c.request(req, &response)
	return response, err
}

type RespToPlaceBulkOrders []Order

func (req *reqToPlaceBulkOrders) path() string {
	return fmt.Sprintf("/order/bulk")
}

func (req *reqToPlaceBulkOrders) method() string {
	return fasthttp.MethodPost
}

func (req *reqToPlaceBulkOrders) query() (string, error) {
	return "", nil
}

func (req *reqToPlaceBulkOrders) payload() (string, error) {

	a := make([]string, 0, 2)

	for i := range req.orders {
		if payload, err := req.orders[i].payload(); err != nil {
			return "", err
		} else {
			a = append(a, payload)
		}
	}

	c := struct {
		Orders []string `json:"orders,omitempty"`
	}{
		a,
	}

	b, err := json.Marshal(&c)
	return string(b), err
}

// reqToCancelAllAfter --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToCancelAllAfter struct {
	c       *restClient
	timeout optional.Float64
}

func (req *reqToCancelAllAfter) Timeout(timeout float64) *reqToCancelAllAfter {
	req.timeout = optional.OfFloat64(timeout)
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
	return fasthttp.MethodPost
}

func (req *reqToCancelAllAfter) query() (string, error) {
	return "", nil
}

func (req *reqToCancelAllAfter) payload() (string, error) {

	a := struct {
		C       *restClient      `json:"-"`
		Timeout optional.Float64 `json:"timeout,omitempty"`
	}{
		req.c,
		req.timeout,
	}

	b, err := json.Marshal(&a)
	return string(b), err
}
