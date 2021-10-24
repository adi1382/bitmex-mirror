package bitmex

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"4d63.com/optional"
	"github.com/google/go-querystring/query"
	"github.com/valyala/fasthttp"
)

//func (c *restClient) GetPositionsRequest() *reqToGetPositions {
//	r := reqToGetPositions{
//		c: c,
//	}
//}
//
//func (c *restClient) IsolateMarginRequest() (*reqToGetPositions, error) {
//	response := new(RespToIsolateMargin)
//	err := c.request(req, response)
//	return response, err
//}
//
//func (c *restClient) ChangeLeverageRequest(req *reqToChangeLeverage) (*reqToGetPositions, error) {
//	response := new(RespToChangeLeverage)
//	err := c.request(req, response)
//	return response, err
//}
//
//func (c *restClient) ChangeRiskLimitRequest(req *reqToChangeRiskLimit) (*reqToGetPositions, error) {
//	response := new(RespToChangeRiskLimit)
//	err := c.request(req, response)
//	return response, err
//}
//
//func (c *restClient) TransferMarginRequest(req *reqToTransferMargin) (*reqToGetPositions, error) {
//	response := new(RespToTransferMargin)
//	err := c.request(req, response)
//	return response, err
//}

type Position struct {
	Account              int       `json:"account"`
	Symbol               string    `json:"symbol"`
	Currency             string    `json:"currency"`
	Underlying           string    `json:"underlying"`
	QuoteCurrency        string    `json:"quoteCurrency"`
	Commission           int       `json:"commission"`
	InitMarginReq        int       `json:"initMarginReq"`
	MaintMarginReq       int       `json:"maintMarginReq"`
	RiskLimit            int       `json:"riskLimit"`
	Leverage             int       `json:"leverage"`
	CrossMargin          bool      `json:"crossMargin"`
	DeleveragePercentile int       `json:"deleveragePercentile"`
	RebalancedPnl        int       `json:"rebalancedPnl"`
	PrevRealisedPnl      int       `json:"prevRealisedPnl"`
	PrevUnrealisedPnl    int       `json:"prevUnrealisedPnl"`
	PrevClosePrice       int       `json:"prevClosePrice"`
	OpeningTimestamp     time.Time `json:"openingTimestamp"`
	OpeningQty           int       `json:"openingQty"`
	OpeningCost          int       `json:"openingCost"`
	OpeningComm          int       `json:"openingComm"`
	OpenOrderBuyQty      int       `json:"openOrderBuyQty"`
	OpenOrderBuyCost     int       `json:"openOrderBuyCost"`
	OpenOrderBuyPremium  int       `json:"openOrderBuyPremium"`
	OpenOrderSellQty     int       `json:"openOrderSellQty"`
	OpenOrderSellCost    int       `json:"openOrderSellCost"`
	OpenOrderSellPremium int       `json:"openOrderSellPremium"`
	ExecBuyQty           int       `json:"execBuyQty"`
	ExecBuyCost          int       `json:"execBuyCost"`
	ExecSellQty          int       `json:"execSellQty"`
	ExecSellCost         int       `json:"execSellCost"`
	ExecQty              int       `json:"execQty"`
	ExecCost             int       `json:"execCost"`
	ExecComm             int       `json:"execComm"`
	CurrentTimestamp     time.Time `json:"currentTimestamp"`
	CurrentQty           int       `json:"currentQty"`
	CurrentCost          int       `json:"currentCost"`
	CurrentComm          int       `json:"currentComm"`
	RealisedCost         int       `json:"realisedCost"`
	UnrealisedCost       int       `json:"unrealisedCost"`
	GrossOpenCost        int       `json:"grossOpenCost"`
	GrossOpenPremium     int       `json:"grossOpenPremium"`
	GrossExecCost        int       `json:"grossExecCost"`
	IsOpen               bool      `json:"isOpen"`
	MarkPrice            int       `json:"markPrice"`
	MarkValue            int       `json:"markValue"`
	RiskValue            int       `json:"riskValue"`
	HomeNotional         int       `json:"homeNotional"`
	ForeignNotional      int       `json:"foreignNotional"`
	PosState             string    `json:"posState"`
	PosCost              int       `json:"posCost"`
	PosCost2             int       `json:"posCost2"`
	PosCross             int       `json:"posCross"`
	PosInit              int       `json:"posInit"`
	PosComm              int       `json:"posComm"`
	PosLoss              int       `json:"posLoss"`
	PosMargin            int       `json:"posMargin"`
	PosMaint             int       `json:"posMaint"`
	PosAllowance         int       `json:"posAllowance"`
	TaxableMargin        int       `json:"taxableMargin"`
	InitMargin           int       `json:"initMargin"`
	MaintMargin          int       `json:"maintMargin"`
	SessionMargin        int       `json:"sessionMargin"`
	TargetExcessMargin   int       `json:"targetExcessMargin"`
	VarMargin            int       `json:"varMargin"`
	RealisedGrossPnl     int       `json:"realisedGrossPnl"`
	RealisedTax          int       `json:"realisedTax"`
	RealisedPnl          int       `json:"realisedPnl"`
	UnrealisedGrossPnl   int       `json:"unrealisedGrossPnl"`
	LongBankrupt         int       `json:"longBankrupt"`
	ShortBankrupt        int       `json:"shortBankrupt"`
	TaxBase              int       `json:"taxBase"`
	IndicativeTaxRate    int       `json:"indicativeTaxRate"`
	IndicativeTax        int       `json:"indicativeTax"`
	UnrealisedTax        int       `json:"unrealisedTax"`
	UnrealisedPnl        int       `json:"unrealisedPnl"`
	UnrealisedPnlPcnt    int       `json:"unrealisedPnlPcnt"`
	UnrealisedRoePcnt    int       `json:"unrealisedRoePcnt"`
	SimpleQty            int       `json:"simpleQty"`
	SimpleCost           int       `json:"simpleCost"`
	SimpleValue          int       `json:"simpleValue"`
	SimplePnl            int       `json:"simplePnl"`
	SimplePnlPcnt        int       `json:"simplePnlPcnt"`
	AvgCostPrice         int       `json:"avgCostPrice"`
	AvgEntryPrice        int       `json:"avgEntryPrice"`
	BreakEvenPrice       int       `json:"breakEvenPrice"`
	MarginCallPrice      int       `json:"marginCallPrice"`
	LiquidationPrice     int       `json:"liquidationPrice"`
	BankruptPrice        int       `json:"bankruptPrice"`
	Timestamp            time.Time `json:"timestamp"`
	LastPrice            int       `json:"lastPrice"`
	LastValue            int       `json:"lastValue"`
}

// reqToGetPositions --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToGetPositions struct {
	c       *restClient
	filter  map[string]interface{}
	columns string
	count   int
}

type RespToGetPositions []Position

func (req *reqToGetPositions) Filter(filter map[string]interface{}) *reqToGetPositions {
	req.filter = filter
	return req
}

func (req *reqToGetPositions) Columns(columns ...string) *reqToGetPositions {
	req.columns = strings.Join(columns, ",")
	return req
}

func (req *reqToGetPositions) Count(count int) *reqToGetPositions {
	req.count = count
	return req
}

func (req *reqToGetPositions) Do() (RespToGetPositions, error) {
	var response RespToGetPositions
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToGetPositions) path() string {
	return fmt.Sprintf("/position")
}

func (req *reqToGetPositions) method() string {
	return fasthttp.MethodGet
}

func (req *reqToGetPositions) query() (string, error) {
	s, e := json.Marshal(req.filter)
	filterStr := string(s)
	if filterStr == "null" || e != nil {
		filterStr = ""
	}

	a := struct {
		Filter  string `url:"filter,omitempty"`
		Columns string `url:"columns,omitempty"`
		Count   int    `url:"count,omitempty"`
	}{
		filterStr,
		req.columns,
		req.count,
	}

	value, err := query.Values(&a)
	if err != nil {
		return "", err
	}
	return value.Encode(), nil
}

func (req *reqToGetPositions) payload() (string, error) {
	return "", nil
}

// reqToIsolateMargin --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToIsolateMargin struct {
	c       *restClient
	symbol  string
	enabled optional.Bool
}

type RespToIsolateMargin Position

func (req *reqToIsolateMargin) Symbol(symbol string) *reqToIsolateMargin {
	req.symbol = symbol
	return req
}

func (req *reqToIsolateMargin) Enabled(enabled bool) *reqToIsolateMargin {
	req.enabled = optional.OfBool(enabled)
	return req
}

func (req *reqToIsolateMargin) Do() (RespToIsolateMargin, error) {
	var response RespToIsolateMargin
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToIsolateMargin) path() string {
	return fmt.Sprintf("/position/isolate")
}

func (req *reqToIsolateMargin) method() string {
	return fasthttp.MethodPost
}

func (req *reqToIsolateMargin) query() (string, error) {
	return "", nil
}

func (req *reqToIsolateMargin) payload() (string, error) {

	a := struct {
		Symbol  string        `json:"symbol,omitempty"`
		Enabled optional.Bool `json:"enabled,omitempty"`
	}{
		req.symbol,
		req.enabled,
	}

	b, err := json.Marshal(&a)
	return string(b), err
}

// reqToChangeLeverage --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToChangeLeverage struct {
	c        *restClient
	symbol   string
	leverage optional.Float64
}

type RespToChangeLeverage Position

func (req *reqToChangeLeverage) Symbol(symbol string) *reqToChangeLeverage {
	req.symbol = symbol
	return req
}

func (req *reqToChangeLeverage) Leverage(leverage float64) *reqToChangeLeverage {
	req.leverage = optional.OfFloat64(leverage)
	return req
}

func (req *reqToChangeLeverage) Do() (RespToChangeLeverage, error) {
	var response RespToChangeLeverage
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToChangeLeverage) path() string {
	return fmt.Sprintf("/position/leverage")
}

func (req *reqToChangeLeverage) method() string {
	return fasthttp.MethodPost
}

func (req *reqToChangeLeverage) query() (string, error) {
	return "", nil
}

func (req *reqToChangeLeverage) payload() (string, error) {

	a := struct {
		Symbol   string           `json:"symbol,omitempty"`
		Leverage optional.Float64 `json:"leverage,omitempty"`
	}{
		req.symbol,
		req.leverage,
	}

	b, err := json.Marshal(&a)
	return string(b), err
}

// reqToChangeRiskLimit --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToChangeRiskLimit struct {
	c         *restClient
	symbol    string
	riskLimit optional.Float64
}

type RespToChangeRiskLimit Position

func (req *reqToChangeRiskLimit) Symbol(symbol string) *reqToChangeRiskLimit {
	req.symbol = symbol
	return req
}

func (req *reqToChangeRiskLimit) RiskLimit(riskLimit float64) *reqToChangeRiskLimit {
	req.riskLimit = optional.OfFloat64(riskLimit)
	return req
}

func (req *reqToChangeRiskLimit) Do() (RespToChangeRiskLimit, error) {
	var response RespToChangeRiskLimit
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToChangeRiskLimit) path() string {
	return fmt.Sprintf("/position/riskLimit")
}

func (req *reqToChangeRiskLimit) method() string {
	return fasthttp.MethodPost
}

func (req *reqToChangeRiskLimit) query() (string, error) {
	return "", nil
}

func (req *reqToChangeRiskLimit) payload() (string, error) {

	a := struct {
		Symbol    string           `json:"symbol,omitempty"`
		RiskLimit optional.Float64 `json:"riskLimit,omitempty"`
	}{
		req.symbol,
		req.riskLimit,
	}

	b, err := json.Marshal(&a)
	return string(b), err
}

// reqToTransferMargin --
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////
type reqToTransferMargin struct {
	c      *restClient
	symbol string
	amount optional.Float64
}

type RespToTransferMargin Position

func (req *reqToTransferMargin) Symbol(symbol string) *reqToTransferMargin {
	req.symbol = symbol
	return req
}

func (req *reqToTransferMargin) Amount(amount float64) *reqToTransferMargin {
	req.amount = optional.OfFloat64(amount)
	return req
}

func (req *reqToTransferMargin) Do() (RespToTransferMargin, error) {
	var response RespToTransferMargin
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToTransferMargin) path() string {
	return fmt.Sprintf("/position/transferMargin")
}

func (req *reqToTransferMargin) method() string {
	return fasthttp.MethodPost
}

func (req *reqToTransferMargin) query() (string, error) {
	return "", nil
}

func (req *reqToTransferMargin) payload() (string, error) {

	a := struct {
		Symbol string           `json:"symbol,omitempty"`
		Amount optional.Float64 `json:"amount,omitempty"`
	}{
		req.symbol,
		req.amount,
	}

	b, err := json.Marshal(&a)
	return string(b), err
}
