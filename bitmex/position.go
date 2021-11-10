package bitmex

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitmex-mirror/optional"
	"github.com/valyala/fasthttp"
)

//type Position struct {
//	Account              int64     `json:"account"`
//	AvgCostPrice         float64   `json:"avgCostPrice"`
//	AvgEntryPrice        float64   `json:"avgEntryPrice"`
//	BankruptPrice        float64   `json:"bankruptPrice"`
//	BreakEvenPrice       float64   `json:"breakEvenPrice"`
//	Commission           float64   `json:"commission"`
//	CrossMargin          bool      `json:"crossMargin"`
//	Currency             string    `json:"currency"`
//	CurrentComm          float64   `json:"currentComm"`
//	CurrentCost          float64   `json:"currentCost"`
//	CurrentQty           float64   `json:"currentQty"`
//	CurrentTimestamp     time.Time `json:"currentTimestamp"`
//	DeleveragePercentile float64   `json:"deleveragePercentile"`
//	ExecBuyCost          float64   `json:"execBuyCost"`
//	ExecBuyQty           float64   `json:"execBuyQty"`
//	ExecComm             float64   `json:"execComm"`
//	ExecCost             float64   `json:"execCost"`
//	ExecQty              float64   `json:"execQty"`
//	ExecSellCost         float64   `json:"execSellCost"`
//	ExecSellQty          float64   `json:"execSellQty"`
//	ForeignNotional      float64   `json:"foreignNotional"`
//	GrossExecCost        float64   `json:"grossExecCost"`
//	GrossOpenCost        float64   `json:"grossOpenCost"`
//	GrossOpenPremium     float64   `json:"grossOpenPremium"`
//	HomeNotional         float64   `json:"homeNotional"`
//	IndicativeTax        float64   `json:"indicativeTax"`
//	IndicativeTaxRate    float64   `json:"indicativeTaxRate"`
//	InitMargin           float64   `json:"initMargin"`
//	InitMarginReq        float64   `json:"initMarginReq"`
//	IsOpen               bool      `json:"isOpen"`
//	LastPrice            float64   `json:"lastPrice"`
//	LastValue            float64   `json:"lastValue"`
//	Leverage             float64   `json:"leverage"`
//	LiquidationPrice     float64   `json:"liquidationPrice"`
//	LongBankrupt         float64   `json:"longBankrupt"`
//	MaintMargin          float64   `json:"maintMargin"`
//	MaintMarginReq       float64   `json:"maintMarginReq"`
//	MarginCallPrice      float64   `json:"marginCallPrice"`
//	MarkPrice            float64   `json:"markPrice"`
//	MarkValue            float64   `json:"markValue"`
//	OpenOrderBuyCost     float64   `json:"openOrderBuyCost"`
//	OpenOrderBuyPremium  float64   `json:"openOrderBuyPremium"`
//	OpenOrderBuyQty      float64   `json:"openOrderBuyQty"`
//	OpenOrderSellCost    float64   `json:"openOrderSellCost"`
//	OpenOrderSellPremium float64   `json:"openOrderSellPremium"`
//	OpenOrderSellQty     float64   `json:"openOrderSellQty"`
//	OpeningComm          float64   `json:"openingComm"`
//	OpeningCost          float64   `json:"openingCost"`
//	OpeningQty           float64   `json:"openingQty"`
//	OpeningTimestamp     time.Time `json:"openingTimestamp"`
//	PosAllowance         float64   `json:"posAllowance"`
//	PosComm              float64   `json:"posComm"`
//	PosCost              float64   `json:"posCost"`
//	PosCost2             float64   `json:"posCost2"`
//	PosCross             float64   `json:"posCross"`
//	PosInit              float64   `json:"posInit"`
//	PosLoss              float64   `json:"posLoss"`
//	PosMaint             float64   `json:"posMaint"`
//	PosMargin            float64   `json:"posMargin"`
//	PosState             string    `json:"posState"`
//	PrevClosePrice       float64   `json:"prevClosePrice"`
//	PrevRealisedPnl      float64   `json:"prevRealisedPnl"`
//	PrevUnrealisedPnl    float64   `json:"prevUnrealisedPnl"`
//	QuoteCurrency        string    `json:"quoteCurrency"`
//	RealisedCost         float64   `json:"realisedCost"`
//	RealisedGrossPnl     float64   `json:"realisedGrossPnl"`
//	RealisedPnl          float64   `json:"realisedPnl"`
//	RealisedTax          float64   `json:"realisedTax"`
//	RebalancedPnl        float64   `json:"rebalancedPnl"`
//	RiskLimit            float64   `json:"riskLimit"`
//	RiskValue            float64   `json:"riskValue"`
//	SessionMargin        float64   `json:"sessionMargin"`
//	ShortBankrupt        float64   `json:"shortBankrupt"`
//	SimpleCost           float64   `json:"simpleCost"`
//	SimplePnl            float64   `json:"simplePnl"`
//	SimplePnlPcnt        float64   `json:"simplePnlPcnt"`
//	SimpleQty            float64   `json:"simpleQty"`
//	SimpleValue          float64   `json:"simpleValue"`
//	Symbol               string    `json:"symbol"`
//	TargetExcessMargin   float64   `json:"targetExcessMargin"`
//	TaxBase              float64   `json:"taxBase"`
//	TaxableMargin        float64   `json:"taxableMargin"`
//	Timestamp            time.Time `json:"timestamp"`
//	Underlying           string    `json:"underlying"`
//	UnrealisedCost       float64   `json:"unrealisedCost"`
//	UnrealisedGrossPnl   float64   `json:"unrealisedGrossPnl"`
//	UnrealisedPnl        float64   `json:"unrealisedPnl"`
//	UnrealisedPnlPcnt    float64   `json:"unrealisedPnlPcnt"`
//	UnrealisedRoePcnt    float64   `json:"unrealisedRoePcnt"`
//	UnrealisedTax        float64   `json:"unrealisedTax"`
//	VarMargin            float64   `json:"varMargin"`
//}

type Position struct {
	Account              int64            `json:"account"`
	AvgCostPrice         optional.Float64 `json:"avgCostPrice"`
	AvgEntryPrice        optional.Float64 `json:"avgEntryPrice"`
	BankruptPrice        optional.Float64 `json:"bankruptPrice"`
	BreakEvenPrice       optional.Float64 `json:"breakEvenPrice"`
	Commission           optional.Float64 `json:"commission"`
	CrossMargin          optional.Bool    `json:"crossMargin"`
	Currency             optional.String  `json:"currency"`
	CurrentComm          optional.Float64 `json:"currentComm"`
	CurrentCost          optional.Float64 `json:"currentCost"`
	CurrentQty           optional.Float64 `json:"currentQty"`
	CurrentTimestamp     optional.Time    `json:"currentTimestamp"`
	DeleveragePercentile optional.Float64 `json:"deleveragePercentile"`
	ExecBuyCost          optional.Float64 `json:"execBuyCost"`
	ExecBuyQty           optional.Float64 `json:"execBuyQty"`
	ExecComm             optional.Float64 `json:"execComm"`
	ExecCost             optional.Float64 `json:"execCost"`
	ExecQty              optional.Float64 `json:"execQty"`
	ExecSellCost         optional.Float64 `json:"execSellCost"`
	ExecSellQty          optional.Float64 `json:"execSellQty"`
	ForeignNotional      optional.Float64 `json:"foreignNotional"`
	GrossExecCost        optional.Float64 `json:"grossExecCost"`
	GrossOpenCost        optional.Float64 `json:"grossOpenCost"`
	GrossOpenPremium     optional.Float64 `json:"grossOpenPremium"`
	HomeNotional         optional.Float64 `json:"homeNotional"`
	IndicativeTax        optional.Float64 `json:"indicativeTax"`
	IndicativeTaxRate    optional.Float64 `json:"indicativeTaxRate"`
	InitMargin           optional.Float64 `json:"initMargin"`
	InitMarginReq        optional.Float64 `json:"initMarginReq"`
	IsOpen               optional.Bool    `json:"isOpen"`
	LastPrice            optional.Float64 `json:"lastPrice"`
	LastValue            optional.Float64 `json:"lastValue"`
	Leverage             optional.Float64 `json:"leverage"`
	LiquidationPrice     optional.Float64 `json:"liquidationPrice"`
	LongBankrupt         optional.Float64 `json:"longBankrupt"`
	MaintMargin          optional.Float64 `json:"maintMargin"`
	MaintMarginReq       optional.Float64 `json:"maintMarginReq"`
	MarginCallPrice      optional.Float64 `json:"marginCallPrice"`
	MarkPrice            optional.Float64 `json:"markPrice"`
	MarkValue            optional.Float64 `json:"markValue"`
	OpenOrderBuyCost     optional.Float64 `json:"openOrderBuyCost"`
	OpenOrderBuyPremium  optional.Float64 `json:"openOrderBuyPremium"`
	OpenOrderBuyQty      optional.Float64 `json:"openOrderBuyQty"`
	OpenOrderSellCost    optional.Float64 `json:"openOrderSellCost"`
	OpenOrderSellPremium optional.Float64 `json:"openOrderSellPremium"`
	OpenOrderSellQty     optional.Float64 `json:"openOrderSellQty"`
	OpeningComm          optional.Float64 `json:"openingComm"`
	OpeningCost          optional.Float64 `json:"openingCost"`
	OpeningQty           optional.Float64 `json:"openingQty"`
	OpeningTimestamp     optional.Time    `json:"openingTimestamp"`
	PosAllowance         optional.Float64 `json:"posAllowance"`
	PosComm              optional.Float64 `json:"posComm"`
	PosCost              optional.Float64 `json:"posCost"`
	PosCost2             optional.Float64 `json:"posCost2"`
	PosCross             optional.Float64 `json:"posCross"`
	PosInit              optional.Float64 `json:"posInit"`
	PosLoss              optional.Float64 `json:"posLoss"`
	PosMaint             optional.Float64 `json:"posMaint"`
	PosMargin            optional.Float64 `json:"posMargin"`
	PosState             optional.String  `json:"posState"`
	PrevClosePrice       optional.Float64 `json:"prevClosePrice"`
	PrevRealisedPnl      optional.Float64 `json:"prevRealisedPnl"`
	PrevUnrealisedPnl    optional.Float64 `json:"prevUnrealisedPnl"`
	QuoteCurrency        optional.String  `json:"quoteCurrency"`
	RealisedCost         optional.Float64 `json:"realisedCost"`
	RealisedGrossPnl     optional.Float64 `json:"realisedGrossPnl"`
	RealisedPnl          optional.Float64 `json:"realisedPnl"`
	RealisedTax          optional.Float64 `json:"realisedTax"`
	RebalancedPnl        optional.Float64 `json:"rebalancedPnl"`
	RiskLimit            optional.Float64 `json:"riskLimit"`
	RiskValue            optional.Float64 `json:"riskValue"`
	SessionMargin        optional.Float64 `json:"sessionMargin"`
	ShortBankrupt        optional.Float64 `json:"shortBankrupt"`
	SimpleCost           optional.Float64 `json:"simpleCost"`
	SimplePnl            optional.Float64 `json:"simplePnl"`
	SimplePnlPcnt        optional.Float64 `json:"simplePnlPcnt"`
	SimpleQty            optional.Float64 `json:"simpleQty"`
	SimpleValue          optional.Float64 `json:"simpleValue"`
	Symbol               string           `json:"symbol"`
	TargetExcessMargin   optional.Float64 `json:"targetExcessMargin"`
	TaxBase              optional.Float64 `json:"taxBase"`
	TaxableMargin        optional.Float64 `json:"taxableMargin"`
	Timestamp            optional.Time    `json:"timestamp"`
	Underlying           optional.String  `json:"underlying"`
	UnrealisedCost       optional.Float64 `json:"unrealisedCost"`
	UnrealisedGrossPnl   optional.Float64 `json:"unrealisedGrossPnl"`
	UnrealisedPnl        optional.Float64 `json:"unrealisedPnl"`
	UnrealisedPnlPcnt    optional.Float64 `json:"unrealisedPnlPcnt"`
	UnrealisedRoePcnt    optional.Float64 `json:"unrealisedRoePcnt"`
	UnrealisedTax        optional.Float64 `json:"unrealisedTax"`
	VarMargin            optional.Float64 `json:"varMargin"`
}

// GetPositions --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetPositions(ctx context.Context, req ReqToGetPositions) (RespToGetPositions, error) {
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
	Count   int64  `json:"count,omitempty"`
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
func (c *RestClient) IsolateMargin(ctx context.Context, req ReqToIsolateMargin) (RespToIsolateMargin, error) {
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
func (c *RestClient) ChangeLeverage(ctx context.Context, req ReqToChangeLeverage) (RespToChangeLeverage, error) {
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
func (c *RestClient) ChangeRiskLimit(ctx context.Context, req ReqToChangeRiskLimit) (RespToChangeRiskLimit, error) {
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
func (c *RestClient) TransferMargin(ctx context.Context, req ReqToTransferMargin) (RespToTransferMargin, error) {
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
