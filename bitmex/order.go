package bitmex

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitmex-mirror/optional"
	"github.com/valyala/fasthttp"
)

//type Order struct {
//	Account               int64     `json:"account"`
//	AvgPx                 float64   `json:"avgPx"`
//	ClOrdID               string    `json:"clOrdID"`
//	ClOrdLinkID           string    `json:"clOrdLinkID"`
//	ContingencyType       string    `json:"contingencyType"`
//	CumQty                float64   `json:"cumQty"`
//	Currency              string    `json:"currency"`
//	DisplayQty            float64   `json:"displayQty"`
//	ExDestination         string    `json:"exDestination"`
//	ExecInst              string    `json:"execInst"`
//	LeavesQty             float64   `json:"leavesQty"`
//	MultiLegReportingType string    `json:"multiLegReportingType"`
//	OrdRejReason          string    `json:"ordRejReason"`
//	OrdStatus             string    `json:"ordStatus"`
//	OrdType               string    `json:"ordType"`
//	OrderID               string    `json:"orderID"`
//	OrderQty              float64   `json:"orderQty"`
//	PegOffsetValue        float64   `json:"pegOffsetValue"`
//	PegPriceType          string    `json:"pegPriceType"`
//	Price                 float64   `json:"price"`
//	SettlCurrency         string    `json:"settlCurrency"`
//	Side                  string    `json:"side"`
//	SimpleCumQty          float64   `json:"simpleCumQty"`
//	SimpleLeavesQty       float64   `json:"simpleLeavesQty"`
//	SimpleOrderQty        float64   `json:"simpleOrderQty"`
//	StopPx                float64   `json:"stopPx"`
//	Symbol                string    `json:"symbol"`
//	Text                  string    `json:"text"`
//	TimeInForce           string    `json:"timeInForce"`
//	Timestamp             time.Time `json:"timestamp"`
//	TransactTime          time.Time `json:"transactTime"`
//	Triggered             string    `json:"triggered"`
//	WorkingIndicator      bool      `json:"workingIndicator"`
//}

type Order struct {
	Account               int64            `json:"account"`
	AvgPx                 optional.Float64 `json:"avgPx"`
	ClOrdID               optional.String  `json:"clOrdID"`
	ClOrdLinkID           optional.String  `json:"clOrdLinkID"`
	ContingencyType       optional.String  `json:"contingencyType"`
	CumQty                optional.Float64 `json:"cumQty"`
	Currency              optional.String  `json:"currency"`
	DisplayQty            optional.Float64 `json:"displayQty"`
	ExDestination         optional.String  `json:"exDestination"`
	ExecInst              optional.String  `json:"execInst"`
	LeavesQty             optional.Float64 `json:"leavesQty"`
	MultiLegReportingType optional.String  `json:"multiLegReportingType"`
	OrdRejReason          optional.String  `json:"ordRejReason"`
	OrdStatus             optional.String  `json:"ordStatus"`
	OrdType               optional.String  `json:"ordType"`
	OrderID               string           `json:"orderID"`
	OrderQty              optional.Float64 `json:"orderQty"`
	PegOffsetValue        optional.Float64 `json:"pegOffsetValue"`
	PegPriceType          optional.String  `json:"pegPriceType"`
	Price                 optional.Float64 `json:"price"`
	SettlCurrency         optional.String  `json:"settlCurrency"`
	Side                  optional.String  `json:"side"`
	SimpleCumQty          optional.Float64 `json:"simpleCumQty"`
	SimpleLeavesQty       optional.Float64 `json:"simpleLeavesQty"`
	SimpleOrderQty        optional.Float64 `json:"simpleOrderQty"`
	StopPx                optional.Float64 `json:"stopPx"`
	Symbol                optional.String  `json:"symbol"`
	Text                  optional.String  `json:"text"`
	TimeInForce           optional.String  `json:"timeInForce"`
	Timestamp             optional.Time    `json:"timestamp"`
	TransactTime          optional.Time    `json:"transactTime"`
	Triggered             optional.String  `json:"triggered"`
	WorkingIndicator      optional.Bool    `json:"workingIndicator"`
}

// GetOrders --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetOrders(ctx context.Context, req ReqToGetOrders) (RespToGetOrders, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetOrders{}, err
	}
	var response []Order
	err = c.request(req, &response)
	return response, err
}

type RespToGetOrders []Order

type ReqToGetOrders struct {
	Symbol    string                 `json:"symbol,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	Columns   string                 `json:"columns,omitempty"`
	Count     int64                  `json:"count,omitempty"`
	Start     int64                  `json:"start,omitempty"`
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
func (c *RestClient) AmendOrder(ctx context.Context, req ReqToAmendOrder) (RespToAmendOrder, error) {
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
func (c *RestClient) PlaceOrder(ctx context.Context, req ReqToPlaceOrder) (RespToPlaceOrder, error) {
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
func (c *RestClient) CancelOrders(ctx context.Context, req ReqToCancelOrders) (RespToCancelOrders, error) {
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

type RespToCancelOrders []Order

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
func (c *RestClient) CancelAllOrders(ctx context.Context, req ReqToCancelAllOrders) (RespToCancelAllOrders, error) {
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
func (c *RestClient) CancelAllAfter(ctx context.Context, req ReqToCancelAllAfter) (RespToCancelAllAfter, error) {
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
