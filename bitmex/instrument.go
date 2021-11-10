package bitmex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitmex-mirror/optional"
	"github.com/valyala/fasthttp"
)

type Instrument struct {
	Symbol                         string    `json:"symbol"`
	RootSymbol                     string    `json:"rootSymbol"`
	State                          string    `json:"state"`
	Typ                            string    `json:"typ"`
	Listing                        time.Time `json:"listing"`
	Front                          time.Time `json:"front"`
	Expiry                         time.Time `json:"expiry"`
	Settle                         time.Time `json:"settle"`
	ListedSettle                   time.Time `json:"listedSettle"`
	RelistInterval                 time.Time `json:"relistInterval"`
	InverseLeg                     string    `json:"inverseLeg"`
	SellLeg                        string    `json:"sellLeg"`
	BuyLeg                         string    `json:"buyLeg"`
	OptionStrikePcnt               int64     `json:"optionStrikePcnt"`
	OptionStrikeRound              int64     `json:"optionStrikeRound"`
	OptionStrikePrice              int64     `json:"optionStrikePrice"`
	OptionMultiplier               int64     `json:"optionMultiplier"`
	PositionCurrency               string    `json:"positionCurrency"`
	Underlying                     string    `json:"underlying"`
	QuoteCurrency                  string    `json:"quoteCurrency"`
	UnderlyingSymbol               string    `json:"underlyingSymbol"`
	Reference                      string    `json:"reference"`
	ReferenceSymbol                string    `json:"referenceSymbol"`
	CalcInterval                   time.Time `json:"calcInterval"`
	PublishInterval                time.Time `json:"publishInterval"`
	PublishTime                    time.Time `json:"publishTime"`
	MaxOrderQty                    int64     `json:"maxOrderQty"`
	MaxPrice                       int64     `json:"maxPrice"`
	LotSize                        int64     `json:"lotSize"`
	TickSize                       int64     `json:"tickSize"`
	Multiplier                     int64     `json:"multiplier"`
	SettlCurrency                  string    `json:"settlCurrency"`
	UnderlyingToPositionMultiplier int64     `json:"underlyingToPositionMultiplier"`
	UnderlyingToSettleMultiplier   int64     `json:"underlyingToSettleMultiplier"`
	QuoteToSettleMultiplier        int64     `json:"quoteToSettleMultiplier"`
	IsQuanto                       bool      `json:"isQuanto"`
	IsInverse                      bool      `json:"isInverse"`
	InitMargin                     int64     `json:"initMargin"`
	MaintMargin                    int64     `json:"maintMargin"`
	RiskLimit                      int64     `json:"riskLimit"`
	RiskStep                       int64     `json:"riskStep"`
	Limit                          int64     `json:"limit"`
	Capped                         bool      `json:"capped"`
	Taxed                          bool      `json:"taxed"`
	Deleverage                     bool      `json:"deleverage"`
	MakerFee                       int64     `json:"makerFee"`
	TakerFee                       int64     `json:"takerFee"`
	SettlementFee                  int64     `json:"settlementFee"`
	InsuranceFee                   int64     `json:"insuranceFee"`
	FundingBaseSymbol              string    `json:"fundingBaseSymbol"`
	FundingQuoteSymbol             string    `json:"fundingQuoteSymbol"`
	FundingPremiumSymbol           string    `json:"fundingPremiumSymbol"`
	FundingTimestamp               time.Time `json:"fundingTimestamp"`
	FundingInterval                time.Time `json:"fundingInterval"`
	FundingRate                    int64     `json:"fundingRate"`
	IndicativeFundingRate          int64     `json:"indicativeFundingRate"`
	RebalanceTimestamp             time.Time `json:"rebalanceTimestamp"`
	RebalanceInterval              time.Time `json:"rebalanceInterval"`
	OpeningTimestamp               time.Time `json:"openingTimestamp"`
	ClosingTimestamp               time.Time `json:"closingTimestamp"`
	SessionInterval                time.Time `json:"sessionInterval"`
	PrevClosePrice                 int64     `json:"prevClosePrice"`
	LimitDownPrice                 int64     `json:"limitDownPrice"`
	LimitUpPrice                   int64     `json:"limitUpPrice"`
	BankruptLimitDownPrice         int64     `json:"bankruptLimitDownPrice"`
	BankruptLimitUpPrice           int64     `json:"bankruptLimitUpPrice"`
	PrevTotalVolume                int64     `json:"prevTotalVolume"`
	TotalVolume                    int64     `json:"totalVolume"`
	Volume                         int64     `json:"volume"`
	Volume24H                      int64     `json:"volume24h"`
	PrevTotalTurnover              int64     `json:"prevTotalTurnover"`
	TotalTurnover                  int64     `json:"totalTurnover"`
	Turnover                       int64     `json:"turnover"`
	Turnover24H                    int64     `json:"turnover24h"`
	HomeNotional24H                int64     `json:"homeNotional24h"`
	ForeignNotional24H             int64     `json:"foreignNotional24h"`
	PrevPrice24H                   int64     `json:"prevPrice24h"`
	Vwap                           int64     `json:"vwap"`
	HighPrice                      int64     `json:"highPrice"`
	LowPrice                       int64     `json:"lowPrice"`
	LastPrice                      int64     `json:"lastPrice"`
	LastPriceProtected             int64     `json:"lastPriceProtected"`
	LastTickDirection              string    `json:"lastTickDirection"`
	LastChangePcnt                 int64     `json:"lastChangePcnt"`
	BidPrice                       int64     `json:"bidPrice"`
	MidPrice                       int64     `json:"midPrice"`
	AskPrice                       int64     `json:"askPrice"`
	ImpactBidPrice                 int64     `json:"impactBidPrice"`
	ImpactMidPrice                 int64     `json:"impactMidPrice"`
	ImpactAskPrice                 int64     `json:"impactAskPrice"`
	HasLiquidity                   bool      `json:"hasLiquidity"`
	OpenInterest                   int64     `json:"openInterest"`
	OpenValue                      int64     `json:"openValue"`
	FairMethod                     string    `json:"fairMethod"`
	FairBasisRate                  int64     `json:"fairBasisRate"`
	FairBasis                      int64     `json:"fairBasis"`
	FairPrice                      int64     `json:"fairPrice"`
	MarkMethod                     string    `json:"markMethod"`
	MarkPrice                      int64     `json:"markPrice"`
	IndicativeTaxRate              int64     `json:"indicativeTaxRate"`
	IndicativeSettlePrice          int64     `json:"indicativeSettlePrice"`
	OptionUnderlyingPrice          int64     `json:"optionUnderlyingPrice"`
	SettledPriceAdjustmentRate     int64     `json:"settledPriceAdjustmentRate"`
	SettledPrice                   int64     `json:"settledPrice"`
	Timestamp                      time.Time `json:"timestamp"`
}

type Intervals struct {
	Intervals []string `json:"intervals"`
	Symbols   []string `json:"symbols"`
}

type CompositeIndex struct {
	Timestamp            time.Time `json:"timestamp"`
	Symbol               string    `json:"symbol"`
	IndexSymbol          string    `json:"indexSymbol"`
	IndexMultiplier      int64     `json:"indexMultiplier"`
	Reference            string    `json:"reference"`
	LastPrice            int64     `json:"lastPrice"`
	SourcePrice          int64     `json:"sourcePrice"`
	ConversionIndex      string    `json:"conversionIndex"`
	ConversionIndexPrice int64     `json:"conversionIndexPrice"`
	Weight               int64     `json:"weight"`
	Logged               time.Time `json:"logged"`
}

// GetInstruments --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetInstruments(ctx context.Context, req ReqToGetInstruments) (RespToGetInstruments, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetInstruments{}, err
	}
	var response RespToGetInstruments
	err = c.request(req, &response)
	return response, err
}

type ReqToGetInstruments struct {
	Symbol    string                 `json:"symbol,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	Columns   string                 `json:"columns,omitempty"`
	Count     int64                  `json:"count,omitempty"`
	Start     int64                  `json:"start,omitempty"`
	Reverse   optional.Bool          `json:"reverse,omitempty"`
	StartTime optional.Time          `json:"startTime,omitempty"`
	EndTime   optional.Time          `json:"endTime,omitempty"`
}

type RespToGetInstruments []Instrument

func (req ReqToGetInstruments) path() string {
	return fmt.Sprintf("/instrument")
}

func (req ReqToGetInstruments) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetInstruments) query() (string, error) {
	return "", nil
}

func (req ReqToGetInstruments) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// GetActiveInstruments --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetActiveInstruments(ctx context.Context,
	req ReqToGetActiveInstruments) (RespToGetActiveInstruments, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetActiveInstruments{}, err
	}
	var response RespToGetActiveInstruments
	err = c.request(req, &response)
	return response, err
}

type ReqToGetActiveInstruments struct{}

type RespToGetActiveInstruments []Instrument

func (req ReqToGetActiveInstruments) path() string {
	return fmt.Sprintf("/instrument/active")
}

func (req ReqToGetActiveInstruments) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetActiveInstruments) query() (string, error) {
	return "", nil
}

func (req ReqToGetActiveInstruments) payload() (string, error) {
	return "", nil
}

// GetActiveInstrumentsAndIndices --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetActiveInstrumentsAndIndices(ctx context.Context,
	req ReqToGetActiveInstrumentsAndIndices) (RespToGetActiveInstrumentsAndIndices, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetActiveInstrumentsAndIndices{}, err
	}
	var response RespToGetActiveInstrumentsAndIndices
	err = c.request(req, &response)

	return response, err
}

type ReqToGetActiveInstrumentsAndIndices struct{}

type RespToGetActiveInstrumentsAndIndices []Instrument

func (req ReqToGetActiveInstrumentsAndIndices) path() string {
	return fmt.Sprintf("/instrument/activeAndIndices")
}

func (req ReqToGetActiveInstrumentsAndIndices) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetActiveInstrumentsAndIndices) query() (string, error) {
	return "", nil
}

func (req ReqToGetActiveInstrumentsAndIndices) payload() (string, error) {
	return "", nil
}

// GetActiveIntervals --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetActiveIntervals(ctx context.Context,
	req ReqToGetActiveIntervals) (RespToGetActiveIntervals, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetActiveIntervals{}, err
	}
	var response RespToGetActiveIntervals
	err = c.request(req, &response)
	return response, err

}

type ReqToGetActiveIntervals struct{}

type RespToGetActiveIntervals Intervals

func (req ReqToGetActiveIntervals) path() string {
	return fmt.Sprintf("/instrument/activeIntervals")
}

func (req ReqToGetActiveIntervals) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetActiveIntervals) query() (string, error) {
	return "", nil
}

func (req ReqToGetActiveIntervals) payload() (string, error) {
	return "", nil
}

// GetCompositeIndex --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetCompositeIndex(ctx context.Context,
	req ReqToGetCompositeIndex) (RespToGetCompositeIndex, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetCompositeIndex{}, err
	}
	var response RespToGetCompositeIndex
	err = c.request(req, &response)
	return response, err
}

type ReqToGetCompositeIndex struct {
	Symbol    string                 `json:"symbol,omitempty"`
	Filter    map[string]interface{} `json:"filter,omitempty"`
	Columns   string                 `json:"columns,omitempty"`
	Count     int64                  `json:"count,omitempty"`
	Start     int64                  `json:"start,omitempty"`
	Reverse   optional.Bool          `json:"reverse,omitempty"`
	StartTime time.Time              `json:"startTime,omitempty"`
	EndTime   time.Time              `json:"endTime,omitempty"`
}

type RespToGetCompositeIndex []CompositeIndex

func (req ReqToGetCompositeIndex) path() string {
	return fmt.Sprintf("/instrument/compositeIndex")
}

func (req ReqToGetCompositeIndex) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetCompositeIndex) query() (string, error) {
	return "", nil
}

func (req ReqToGetCompositeIndex) payload() (string, error) {
	b, err := json.Marshal(&req)
	return string(b), err
}

// GetAllIndices --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *RestClient) GetAllIndices(ctx context.Context, req ReqToGetAllIndices) (RespToGetAllIndices, error) {
	err := c.bucketM.Wait(ctx)
	if err != nil {
		return RespToGetAllIndices{}, err
	}
	var response RespToGetAllIndices
	err = c.request(req, response)
	return response, err
}

type ReqToGetAllIndices struct{}

type RespToGetAllIndices []Instrument

func (req ReqToGetAllIndices) path() string {
	return fmt.Sprintf("/instrument/indices")
}

func (req ReqToGetAllIndices) method() string {
	return fasthttp.MethodGet
}

func (req ReqToGetAllIndices) query() (string, error) {
	return "", nil
}

func (req ReqToGetAllIndices) payload() (string, error) {
	return "", nil
}
