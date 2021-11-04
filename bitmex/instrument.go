package bitmex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"4d63.com/optional"
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
	OptionStrikePcnt               int       `json:"optionStrikePcnt"`
	OptionStrikeRound              int       `json:"optionStrikeRound"`
	OptionStrikePrice              int       `json:"optionStrikePrice"`
	OptionMultiplier               int       `json:"optionMultiplier"`
	PositionCurrency               string    `json:"positionCurrency"`
	Underlying                     string    `json:"underlying"`
	QuoteCurrency                  string    `json:"quoteCurrency"`
	UnderlyingSymbol               string    `json:"underlyingSymbol"`
	Reference                      string    `json:"reference"`
	ReferenceSymbol                string    `json:"referenceSymbol"`
	CalcInterval                   time.Time `json:"calcInterval"`
	PublishInterval                time.Time `json:"publishInterval"`
	PublishTime                    time.Time `json:"publishTime"`
	MaxOrderQty                    int       `json:"maxOrderQty"`
	MaxPrice                       int       `json:"maxPrice"`
	LotSize                        int       `json:"lotSize"`
	TickSize                       int       `json:"tickSize"`
	Multiplier                     int       `json:"multiplier"`
	SettlCurrency                  string    `json:"settlCurrency"`
	UnderlyingToPositionMultiplier int       `json:"underlyingToPositionMultiplier"`
	UnderlyingToSettleMultiplier   int       `json:"underlyingToSettleMultiplier"`
	QuoteToSettleMultiplier        int       `json:"quoteToSettleMultiplier"`
	IsQuanto                       bool      `json:"isQuanto"`
	IsInverse                      bool      `json:"isInverse"`
	InitMargin                     int       `json:"initMargin"`
	MaintMargin                    int       `json:"maintMargin"`
	RiskLimit                      int       `json:"riskLimit"`
	RiskStep                       int       `json:"riskStep"`
	Limit                          int       `json:"limit"`
	Capped                         bool      `json:"capped"`
	Taxed                          bool      `json:"taxed"`
	Deleverage                     bool      `json:"deleverage"`
	MakerFee                       int       `json:"makerFee"`
	TakerFee                       int       `json:"takerFee"`
	SettlementFee                  int       `json:"settlementFee"`
	InsuranceFee                   int       `json:"insuranceFee"`
	FundingBaseSymbol              string    `json:"fundingBaseSymbol"`
	FundingQuoteSymbol             string    `json:"fundingQuoteSymbol"`
	FundingPremiumSymbol           string    `json:"fundingPremiumSymbol"`
	FundingTimestamp               time.Time `json:"fundingTimestamp"`
	FundingInterval                time.Time `json:"fundingInterval"`
	FundingRate                    int       `json:"fundingRate"`
	IndicativeFundingRate          int       `json:"indicativeFundingRate"`
	RebalanceTimestamp             time.Time `json:"rebalanceTimestamp"`
	RebalanceInterval              time.Time `json:"rebalanceInterval"`
	OpeningTimestamp               time.Time `json:"openingTimestamp"`
	ClosingTimestamp               time.Time `json:"closingTimestamp"`
	SessionInterval                time.Time `json:"sessionInterval"`
	PrevClosePrice                 int       `json:"prevClosePrice"`
	LimitDownPrice                 int       `json:"limitDownPrice"`
	LimitUpPrice                   int       `json:"limitUpPrice"`
	BankruptLimitDownPrice         int       `json:"bankruptLimitDownPrice"`
	BankruptLimitUpPrice           int       `json:"bankruptLimitUpPrice"`
	PrevTotalVolume                int       `json:"prevTotalVolume"`
	TotalVolume                    int       `json:"totalVolume"`
	Volume                         int       `json:"volume"`
	Volume24H                      int       `json:"volume24h"`
	PrevTotalTurnover              int       `json:"prevTotalTurnover"`
	TotalTurnover                  int       `json:"totalTurnover"`
	Turnover                       int       `json:"turnover"`
	Turnover24H                    int       `json:"turnover24h"`
	HomeNotional24H                int       `json:"homeNotional24h"`
	ForeignNotional24H             int       `json:"foreignNotional24h"`
	PrevPrice24H                   int       `json:"prevPrice24h"`
	Vwap                           int       `json:"vwap"`
	HighPrice                      int       `json:"highPrice"`
	LowPrice                       int       `json:"lowPrice"`
	LastPrice                      int       `json:"lastPrice"`
	LastPriceProtected             int       `json:"lastPriceProtected"`
	LastTickDirection              string    `json:"lastTickDirection"`
	LastChangePcnt                 int       `json:"lastChangePcnt"`
	BidPrice                       int       `json:"bidPrice"`
	MidPrice                       int       `json:"midPrice"`
	AskPrice                       int       `json:"askPrice"`
	ImpactBidPrice                 int       `json:"impactBidPrice"`
	ImpactMidPrice                 int       `json:"impactMidPrice"`
	ImpactAskPrice                 int       `json:"impactAskPrice"`
	HasLiquidity                   bool      `json:"hasLiquidity"`
	OpenInterest                   int       `json:"openInterest"`
	OpenValue                      int       `json:"openValue"`
	FairMethod                     string    `json:"fairMethod"`
	FairBasisRate                  int       `json:"fairBasisRate"`
	FairBasis                      int       `json:"fairBasis"`
	FairPrice                      int       `json:"fairPrice"`
	MarkMethod                     string    `json:"markMethod"`
	MarkPrice                      int       `json:"markPrice"`
	IndicativeTaxRate              int       `json:"indicativeTaxRate"`
	IndicativeSettlePrice          int       `json:"indicativeSettlePrice"`
	OptionUnderlyingPrice          int       `json:"optionUnderlyingPrice"`
	SettledPriceAdjustmentRate     int       `json:"settledPriceAdjustmentRate"`
	SettledPrice                   int       `json:"settledPrice"`
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
	IndexMultiplier      int       `json:"indexMultiplier"`
	Reference            string    `json:"reference"`
	LastPrice            int       `json:"lastPrice"`
	SourcePrice          int       `json:"sourcePrice"`
	ConversionIndex      string    `json:"conversionIndex"`
	ConversionIndexPrice int       `json:"conversionIndexPrice"`
	Weight               int       `json:"weight"`
	Logged               time.Time `json:"logged"`
}

// GetInstruments --
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////
func (c *restClient) GetInstruments(ctx context.Context, req ReqToGetInstruments) (RespToGetInstruments, error) {
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
	Count     int                    `json:"count,omitempty"`
	Start     int                    `json:"start,omitempty"`
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
func (c *restClient) GetActiveInstruments(ctx context.Context,
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
func (c *restClient) GetActiveInstrumentsAndIndices(ctx context.Context,
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
func (c *restClient) GetActiveIntervals(ctx context.Context,
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
func (c *restClient) GetCompositeIndex(ctx context.Context,
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
	Count     int                    `json:"count,omitempty"`
	Start     int                    `json:"start,omitempty"`
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
func (c *restClient) GetAllIndices(ctx context.Context, req ReqToGetAllIndices) (RespToGetAllIndices, error) {
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
