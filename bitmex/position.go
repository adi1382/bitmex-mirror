package bitmex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"4d63.com/optional"
	"github.com/valyala/fasthttp"
)

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

// GetPositions --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) GetPositions(ctx context.Context, req ReqToGetPositions) (RespToGetPositions, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetPositions{}, err
	}
	var response RespToGetPositions
	err = c.request(req, &response)
	return response, err
}

type ReqToGetPositions struct {
	Filter  string `json:"filter,omitempty"`
	Columns string `json:"columns,omitempty"`
	Count   int    `json:"count,omitempty"`
}

type RespToGetPositions []Position

func (req ReqToGetPositions) path() string {
	return fmt.Sprintf("/position")
}

func (req ReqToGetPositions) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetPositions) query() (string, error) {
	return "", nil
}

func (req ReqToGetPositions) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// IsolateMargin --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) IsolateMargin(ctx context.Context, req ReqToIsolateMargin) (RespToIsolateMargin, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToIsolateMargin{}, err
	}
	var response RespToIsolateMargin
	err = c.request(req, &response)
	return response, err
}

type ReqToIsolateMargin struct {
	Symbol  string        `json:"symbol,omitempty"`
	Enabled optional.Bool `json:"enabled,omitempty"`
}

type RespToIsolateMargin Position

func (req ReqToIsolateMargin) path() string {
	return fmt.Sprintf("/position/isolate")
}

func (req ReqToIsolateMargin) method() string {
	return fasthttp.MethodPost
}

func (req ReqToIsolateMargin) query() (string, error) {
	return "", nil
}

func (req ReqToIsolateMargin) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// ChangeLeverage --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) ChangeLeverage(ctx context.Context, req ReqToChangeLeverage) (RespToChangeLeverage, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToChangeLeverage{}, err
	}
	var response RespToChangeLeverage
	err = c.request(req, &response)
	return response, err
}

type ReqToChangeLeverage struct {
	Symbol   string           `json:"symbol,omitempty"`
	Leverage optional.Float64 `json:"leverage,omitempty"`
}

type RespToChangeLeverage Position

func (req ReqToChangeLeverage) path() string {
	return fmt.Sprintf("/position/leverage")
}

func (req ReqToChangeLeverage) method() string {
	return fasthttp.MethodPost
}

func (req ReqToChangeLeverage) query() (string, error) {
	return "", nil
}

func (req ReqToChangeLeverage) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// ChangeRiskLimit --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) ChangeRiskLimit(ctx context.Context, req ReqToChangeRiskLimit) (RespToChangeRiskLimit, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToChangeRiskLimit{}, err
	}
	var response RespToChangeRiskLimit
	err = c.request(req, &response)
	return response, err
}

type ReqToChangeRiskLimit struct {
	Symbol    string           `json:"symbol,omitempty"`
	RiskLimit optional.Float64 `json:"riskLimit,omitempty"`
}

type RespToChangeRiskLimit Position

func (req ReqToChangeRiskLimit) path() string {
	return fmt.Sprintf("/position/riskLimit")
}

func (req ReqToChangeRiskLimit) method() string {
	return fasthttp.MethodPost
}

func (req ReqToChangeRiskLimit) query() (string, error) {
	return "", nil
}

func (req ReqToChangeRiskLimit) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// TransferMargin --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) TransferMargin(ctx context.Context, req ReqToTransferMargin) (RespToTransferMargin, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToTransferMargin{}, err
	}
	var response RespToTransferMargin
	err = c.request(req, &response)
	return response, err
}

type ReqToTransferMargin struct {
	Symbol string           `json:"symbol,omitempty"`
	Amount optional.Float64 `json:"amount,omitempty"`
}

type RespToTransferMargin Position

func (req ReqToTransferMargin) path() string {
	return fmt.Sprintf("/position/transferMargin")
}

func (req ReqToTransferMargin) method() string {
	return fasthttp.MethodPost
}

func (req ReqToTransferMargin) query() (string, error) {
	return "", nil
}

func (req ReqToTransferMargin) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}
