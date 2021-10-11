package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"strings"
	"time"
)

func (c *Client) GetPositions(req *reqToGetPositions) (*RespToGetPositions, error) {
	response := new(RespToGetPositions)
	err := c.request(req, response)
	return response, err
}

func (c *Client) ChangeLeverage(req *ReqToChangeLeverage) (*RespToChangeLeverage, error) {
	response := new(RespToChangeLeverage)
	err := c.request(req, response)
	return response, err
}

func (c *Client) ChangeRiskLimit(req *ReqToChangeRiskLimit) (*RespToChangeRiskLimit, error) {
	response := new(RespToChangeRiskLimit)
	err := c.request(req, response)
	return response, err
}

func (c *Client) IsolateMargin(req *ReqToIsolateMargin) (*RespToIsolateMargin, error) {
	response := new(RespToIsolateMargin)
	err := c.request(req, response)
	return response, err
}

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

type reqToGetPositions struct {
	filter  map[string]interface{} `url:"filter,omitempty"`
	columns string                 `url:"columns,omitempty"`
	count   int                    `url:"count,omitempty"`
}

func (req *reqToGetPositions) AddFilter(key string, value interface{}) *reqToGetPositions {
	req.filter[key] = value
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

type RespToGetPositions []Position

func (req *reqToGetPositions) path() string {
	return fmt.Sprintf("/position")
}

func (req *reqToGetPositions) method() string {
	return http.MethodPost
}

func (req *reqToGetPositions) query() string {
	value, _ := query.Values(req)
	return value.Encode()
}

func (req *reqToGetPositions) payload() string {
	return ""
}

// ReqToIsolateMargin --
type ReqToIsolateMargin struct {
	Symbol  string `url:"symbol,omitempty" json:"symbol,omitempty"`
	Enabled *bool  `url:"enabled,omitempty" json:"enabled,omitempty"`
}

type RespToIsolateMargin Position

func (req *ReqToIsolateMargin) path() string {
	return fmt.Sprintf("/position/isolate")
}

func (req *ReqToIsolateMargin) method() string {
	return http.MethodPost
}

func (req *ReqToIsolateMargin) query() string {
	return ""
}

func (req *ReqToIsolateMargin) payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}

// ReqToChangeLeverage --
type ReqToChangeLeverage struct {
	Symbol   string   `url:"symbol" json:"symbol"`
	Leverage *float64 `url:"leverage" json:"leverage"`
}

type RespToChangeLeverage Position

func (req *ReqToChangeLeverage) path() string {
	return fmt.Sprintf("/position/leverage")
}

func (req *ReqToChangeLeverage) method() string {
	return http.MethodPost
}

func (req *ReqToChangeLeverage) query() string {
	return ""
}

func (req *ReqToChangeLeverage) payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}

// ReqToChangeRiskLimit --
type ReqToChangeRiskLimit struct {
	Symbol    string   `url:"symbol" json:"symbol"`
	RiskLimit *float64 `url:"riskLimit" json:"riskLimit"`
}

type RespToChangeRiskLimit Position

func (req *ReqToChangeRiskLimit) path() string {
	return fmt.Sprintf("/position/riskLimit")
}

func (req *ReqToChangeRiskLimit) method() string {
	return http.MethodPost
}

func (req *ReqToChangeRiskLimit) query() string {
	return ""
}

func (req *ReqToChangeRiskLimit) payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}

// ReqToTransferMargin --
type ReqToTransferMargin struct {
	Symbol string   `url:"symbol" json:"symbol"`
	Amount *float64 `url:"amount" json:"amount"`
}

type RespToTransferMargin Position

func (req *ReqToTransferMargin) path() string {
	return fmt.Sprintf("/position/transferMargin")
}

func (req *ReqToTransferMargin) method() string {
	return http.MethodPost
}

func (req *ReqToTransferMargin) query() string {
	return ""
}

func (req *ReqToTransferMargin) payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
