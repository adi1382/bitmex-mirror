package bitmex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"4d63.com/optional"
	"github.com/valyala/fasthttp"
)

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

// GetOrders --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) GetOrders(ctx context.Context, req ReqToGetOrders) (RespToGetOrders, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetOrders{}, err
	}
	var response RespToGetOrders
	err = c.request(req, &response)
	return response, err
}

type RespToGetOrders []Order

type ReqToGetOrders struct {
	Symbol    string                 `json:"symbol,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	Columns   string                 `json:"columns,omitempty"`
	Count     int                    `json:"count,omitempty"`
	Start     int                    `json:"start,omitempty"`
	Reverse   optional.Bool          `json:"reverse,omitempty"`
	StartTime optional.Time          `json:"startTime,omitempty"`
	EndTime   optional.Time          `json:"endTime,omitempty"`
}

func (req ReqToGetOrders) path() string {
	return fmt.Sprintf("/order")
}

func (req ReqToGetOrders) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetOrders) query() (string, error) {
	return "", nil
}

func (req ReqToGetOrders) payload() (string, error) {
	b, err := json.Marshal(&req)
	fmt.Println("Payload: ", string(b))
	return string(b), err
}

// AmendOrder --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) AmendOrder(ctx context.Context, req ReqToAmendOrder) (RespToAmendOrder, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToAmendOrder{}, err
	}
	err = c.bucketS.Wait(ctx)
	if err != nil {
		return RespToAmendOrder{}, err
	}
	response := RespToAmendOrder{}
	err = c.request(req, &response)
	return response, err
}

type RespToAmendOrder Order

type ReqToAmendOrder struct {
	OrderID        string           `json:"orderID,omitempty"`
	OrigClOrdID    string           `json:"origClOrdID,omitempty"`
	ClOrdID        string           `json:"clOrdID,omitempty"`
	OrderQty       optional.Float64 `json:"orderQty,omitempty"`
	LeavesQty      optional.Float64 `json:"leavesQty,omitempty"`
	Price          optional.Float64 `json:"price,omitempty"`
	StopPx         optional.Float64 `json:"stopPx,omitempty"`
	PegOffsetValue optional.Float64 `json:"pefOffsetValue,omitempty"`
	Text           string           `json:"text,omitempty"`
}

func (req ReqToAmendOrder) path() string {
	return fmt.Sprintf("/order")
}

func (req ReqToAmendOrder) method() string {
	return fasthttp.MethodPut
}

func (req ReqToAmendOrder) query() (string, error) {
	return "", nil
}

func (req ReqToAmendOrder) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// PlaceOrder  --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) PlaceOrder(ctx context.Context, req ReqToPlaceOrder) (RespToPlaceOrder, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToPlaceOrder{}, err
	}
	err = c.bucketS.Wait(ctx)
	if err != nil {
		return RespToPlaceOrder{}, err
	}
	var response RespToPlaceOrder
	err = c.request(req, &response)
	return response, err
}

type RespToPlaceOrder Order

type ReqToPlaceOrder struct {
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
}

func (req ReqToPlaceOrder) path() string {
	return fmt.Sprintf("/order")
}

func (req ReqToPlaceOrder) method() string {
	return fasthttp.MethodPost
}

func (req ReqToPlaceOrder) query() (string, error) {
	return "", nil
}

func (req ReqToPlaceOrder) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// CancelOrders --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) CancelOrders(ctx context.Context, req ReqToCancelOrders) (RespToCancelOrders, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToCancelOrders{}, err
	}
	err = c.bucketS.Wait(ctx)
	if err != nil {
		return RespToCancelOrders{}, err
	}
	var response RespToCancelOrders
	err = c.request(req, &response)
	return response, err
}

type RespToCancelOrders Order

type ReqToCancelOrders struct {
	OrderIDs []string `json:"orderID,omitempty"`
	ClOrdIDs []string `json:"clOrdID,omitempty"`
	Text     string   `json:"text,omitempty"`
}

func (req ReqToCancelOrders) path() string {
	return fmt.Sprintf("/order")
}

func (req ReqToCancelOrders) method() string {
	return fasthttp.MethodDelete
}

func (req ReqToCancelOrders) query() (string, error) {
	return "", nil
}

func (req ReqToCancelOrders) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// CancelAllOrders --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) CancelAllOrders(ctx context.Context, req ReqToCancelAllOrders) (RespToCancelAllOrders, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToCancelAllOrders{}, err
	}
	err = c.bucketS.Wait(ctx)
	if err != nil {
		return RespToCancelAllOrders{}, err
	}
	var response RespToCancelAllOrders
	err = c.request(req, &response)
	return response, err
}

type RespToCancelAllOrders []Order

type ReqToCancelAllOrders struct {
	Symbol string                 `json:"symbol,omitempty"`
	Filter map[string]interface{} `json:"filter,omitempty"`
	Text   string                 `json:"text,omitempty"`
}

func (req ReqToCancelAllOrders) path() string {
	return fmt.Sprintf("/order/all")
}

func (req ReqToCancelAllOrders) method() string {
	return fasthttp.MethodDelete
}

func (req ReqToCancelAllOrders) query() (string, error) {
	return "", nil
}

func (req ReqToCancelAllOrders) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// CancelAllAfter implements the cancelAllAfter method
// https://www.bitmex.com/api/explorer/#!/Order/Order_cancelAllAfter
// POST /order/cancelAllAfter
func (c *restClient) CancelAllAfter(ctx context.Context, req ReqToCancelAllAfter) (RespToCancelAllAfter, error) {
	var response RespToCancelAllAfter
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToCancelAllAfter{}, err
	}
	err = c.request(req, &response)
	return response, err
}

type RespToCancelAllAfter struct{}

type ReqToCancelAllAfter struct {
	Timeout optional.Float64 `json:"timeout,omitempty"`
}

func (req ReqToCancelAllAfter) path() string {
	return fmt.Sprintf("/order/cancelAllAfter")
}

func (req ReqToCancelAllAfter) method() string {
	return fasthttp.MethodPost
}

func (req ReqToCancelAllAfter) query() (string, error) {
	return "", nil
}

func (req ReqToCancelAllAfter) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}
