package position

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
	"time"
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

type ReqToGetPositions struct {
	Filter  map[string]string `url:"filter,omitempty"`
	Columns string            `url:"columns,omitempty"`
	Count   int               `url:"count,omitempty"`
}

type RespToGetPositions []Position

func (req *ReqToGetPositions) path() string {
	return fmt.Sprintf("/position")
}

func (req *ReqToGetPositions) Method() string {
	return http.MethodPost
}

func (req *ReqToGetPositions) Query() string {
	value, _ := query.Values(req)
	return value.Encode()
}

func (req *ReqToGetPositions) Payload() string {
	return ""
}
