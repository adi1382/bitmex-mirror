package bitmex

import "github.com/bitmex-mirror/optional"

//type Margin struct {
//	Account            int64     `json:"account"`
//	Action             string    `json:"action"`
//	Amount             int64     `json:"amount"`
//	AvailableMargin    int64     `json:"availableMargin"`
//	Commission         float64   `json:"commission"`
//	ConfirmedDebit     int64     `json:"confirmedDebit"`
//	Currency           string    `json:"currency"`
//	ExcessMargin       int64     `json:"excessMargin"`
//	ExcessMarginPcnt   float64   `json:"excessMarginPcnt"`
//	GrossComm          int64     `json:"grossComm"`
//	GrossExecCost      int64     `json:"grossExecCost"`
//	GrossLastValue     int64     `json:"grossLastValue"`
//	GrossMarkValue     int64     `json:"grossMarkValue"`
//	GrossOpenCost      int64     `json:"grossOpenCost"`
//	GrossOpenPremium   int64     `json:"grossOpenPremium"`
//	IndicativeTax      int64     `json:"indicativeTax"`
//	InitMargin         int64     `json:"initMargin"`
//	MaintMargin        int64     `json:"maintMargin"`
//	MakerFeeDiscount   float64   `json:"makerFeeDiscount"`
//	MarginBalance      int64     `json:"marginBalance"`
//	MarginBalancePcnt  float64   `json:"marginBalancePcnt"`
//	MarginLeverage     float64   `json:"marginLeverage"`
//	MarginUsedPcnt     float64   `json:"marginUsedPcnt"`
//	PendingCredit      int64     `json:"pendingCredit"`
//	PendingDebit       int64     `json:"pendingDebit"`
//	PrevRealisedPnl    int64     `json:"prevRealisedPnl"`
//	PrevState          string    `json:"prevState"`
//	PrevUnrealisedPnl  int64     `json:"prevUnrealisedPnl"`
//	RealisedPnl        int64     `json:"realisedPnl"`
//	RiskLimit          int64     `json:"riskLimit"`
//	RiskValue          int64     `json:"riskValue"`
//	SessionMargin      int64     `json:"sessionMargin"`
//	State              string    `json:"state"`
//	SyntheticMargin    int64     `json:"syntheticMargin"`
//	TakerFeeDiscount   float64   `json:"takerFeeDiscount"`
//	TargetExcessMargin int64     `json:"targetExcessMargin"`
//	TaxableMargin      int64     `json:"taxableMargin"`
//	Timestamp          time.Time `json:"timestamp"`
//	UnrealisedPnl      int64     `json:"unrealisedPnl"`
//	UnrealisedProfit   int64     `json:"unrealisedProfit"`
//	VarMargin          int64     `json:"varMargin"`
//	WalletBalance      int64     `json:"walletBalance"`
//	WithdrawableMargin int64     `json:"withdrawableMargin"`
//}

type Margin struct {
	Account            int64            `json:"account"`
	Action             optional.String  `json:"action"`
	Amount             optional.Float64 `json:"amount"`
	AvailableMargin    optional.Float64 `json:"availableMargin"`
	Commission         optional.Float64 `json:"commission"`
	ConfirmedDebit     optional.Float64 `json:"confirmedDebit"`
	Currency           string           `json:"currency"`
	ExcessMargin       optional.Float64 `json:"excessMargin"`
	ExcessMarginPcnt   optional.Float64 `json:"excessMarginPcnt"`
	GrossComm          optional.Float64 `json:"grossComm"`
	GrossExecCost      optional.Float64 `json:"grossExecCost"`
	GrossLastValue     optional.Float64 `json:"grossLastValue"`
	GrossMarkValue     optional.Float64 `json:"grossMarkValue"`
	GrossOpenCost      optional.Float64 `json:"grossOpenCost"`
	GrossOpenPremium   optional.Float64 `json:"grossOpenPremium"`
	IndicativeTax      optional.Float64 `json:"indicativeTax"`
	InitMargin         optional.Float64 `json:"initMargin"`
	MaintMargin        optional.Float64 `json:"maintMargin"`
	MakerFeeDiscount   optional.Float64 `json:"makerFeeDiscount"`
	MarginBalance      optional.Float64 `json:"marginBalance"`
	MarginBalancePcnt  optional.Float64 `json:"marginBalancePcnt"`
	MarginLeverage     optional.Float64 `json:"marginLeverage"`
	MarginUsedPcnt     optional.Float64 `json:"marginUsedPcnt"`
	PendingCredit      optional.Float64 `json:"pendingCredit"`
	PendingDebit       optional.Float64 `json:"pendingDebit"`
	PrevRealisedPnl    optional.Float64 `json:"prevRealisedPnl"`
	PrevState          optional.String  `json:"prevState"`
	PrevUnrealisedPnl  optional.Float64 `json:"prevUnrealisedPnl"`
	RealisedPnl        optional.Float64 `json:"realisedPnl"`
	RiskLimit          optional.Float64 `json:"riskLimit"`
	RiskValue          optional.Float64 `json:"riskValue"`
	SessionMargin      optional.Float64 `json:"sessionMargin"`
	State              optional.String  `json:"state"`
	SyntheticMargin    optional.Float64 `json:"syntheticMargin"`
	TakerFeeDiscount   optional.Float64 `json:"takerFeeDiscount"`
	TargetExcessMargin optional.Float64 `json:"targetExcessMargin"`
	TaxableMargin      optional.Float64 `json:"taxableMargin"`
	Timestamp          optional.Time    `json:"timestamp"`
	UnrealisedPnl      optional.Float64 `json:"unrealisedPnl"`
	UnrealisedProfit   optional.Float64 `json:"unrealisedProfit"`
	VarMargin          optional.Float64 `json:"varMargin"`
	WalletBalance      optional.Float64 `json:"walletBalance"`
	WithdrawableMargin optional.Float64 `json:"withdrawableMargin"`
}
