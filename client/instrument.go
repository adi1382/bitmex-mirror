package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"4d63.com/optional"
	"github.com/google/go-querystring/query"
)

func (c *Client) GetInstruments() *reqToGetInstruments {
	return &reqToGetInstruments{
		c: c,
	}
}

func (c *Client) GetActiveInstruments() *reqToGetActiveInstruments {
	return &reqToGetActiveInstruments{
		c: c,
	}
}

func (c *Client) GetActiveInstrumentsAndIndices() *reqToGetActiveInstrumentsAndIndices {
	return &reqToGetActiveInstrumentsAndIndices{
		c: c,
	}
}

func (c *Client) GetActiveIntervals() *reqToGetActiveIntervals {
	return &reqToGetActiveIntervals{
		c: c,
	}
}

func (c *Client) GetCompositeIndex() *reqToGetCompositeIndex {
	return &reqToGetCompositeIndex{
		c: c,
	}
}

func (c *Client) GetAllIndices() *reqToGetAllIndices {
	return &reqToGetAllIndices{
		c: c,
	}
}

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

type reqToGetInstruments struct {
	c         *Client
	symbol    string
	filter    map[string]interface{}
	columns   string
	count     int
	start     int
	reverse   optional.Bool
	startTime time.Time
	endTime   time.Time
}

type RespToGetInstruments struct {
}

func (req *reqToGetInstruments) Symbol(symbol string) *reqToGetInstruments {
	req.symbol = symbol
	return req
}

func (req *reqToGetInstruments) Filter(filter map[string]interface{}) *reqToGetInstruments {
	req.filter = filter
	return req
}

func (req *reqToGetInstruments) Columns(columns ...string) *reqToGetInstruments {
	req.columns = strings.Join(columns, ",")
	return req
}

func (req *reqToGetInstruments) Count(count int) *reqToGetInstruments {
	req.count = count
	return req
}

func (req *reqToGetInstruments) Start(start int) *reqToGetInstruments {
	req.start = start
	return req
}

func (req *reqToGetInstruments) Reverse(reverse bool) *reqToGetInstruments {
	req.reverse = optional.Bool{reverse}
	return req
}

func (req *reqToGetInstruments) StartTime(startTime time.Time) *reqToGetInstruments {
	req.startTime = startTime
	return req
}

func (req *reqToGetInstruments) EndTime(endTime time.Time) *reqToGetInstruments {
	req.endTime = endTime
	return req
}

func (req *reqToGetInstruments) Do() (RespToGetInstruments, error) {
	var response RespToGetInstruments
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToGetInstruments) path() string {
	return fmt.Sprintf("/instrument")
}

func (req *reqToGetInstruments) method() string {
	return http.MethodGet
}

func (req *reqToGetInstruments) query() (string, error) {
	s, e := json.Marshal(req.filter)
	filterStr := string(s)
	if filterStr == "null" || e != nil {
		filterStr = ""
	}

	a := struct {
		C         *Client       `url:"-"`
		Symbol    string        `url:"symbol,omitempty"`
		Filter    string        `url:"filter,omitempty"`
		Columns   string        `url:"columns,omitempty"`
		Count     int           `url:"count,omitempty"`
		Start     int           `url:"start,omitempty"`
		Reverse   optional.Bool `url:"reverse,omitempty"`
		StartTime time.Time     `url:"startTime,omitempty"`
		EndTime   time.Time     `url:"endTime,omitempty"`
	}{
		req.c,
		req.symbol,
		filterStr,
		req.columns,
		req.count,
		req.start,
		req.reverse,
		req.startTime,
		req.endTime,
	}
	value, err := query.Values(&a)
	if err != nil {
		return "", err
	}

	return value.Encode(), nil
}

func (req *reqToGetInstruments) payload() (string, error) {
	return "", nil
}

type reqToGetActiveInstruments struct {
	c *Client
}

type RespToGetActiveInstruments []Instrument

func (req *reqToGetActiveInstruments) Do() (RespToGetActiveInstruments, error) {
	var response RespToGetActiveInstruments
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToGetActiveInstruments) path() string {
	return fmt.Sprintf("/instrument/active")
}

func (req *reqToGetActiveInstruments) method() string {
	return http.MethodGet
}

func (req *reqToGetActiveInstruments) query() (string, error) {
	return "", nil
}

func (req *reqToGetActiveInstruments) payload() (string, error) {
	return "", nil
}

type reqToGetActiveInstrumentsAndIndices struct {
	c *Client
}

type RespToGetActiveInstrumentsAndIndices []Instrument

func (req *reqToGetActiveInstrumentsAndIndices) Do() (RespToGetActiveInstrumentsAndIndices, error) {
	var response RespToGetActiveInstrumentsAndIndices
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToGetActiveInstrumentsAndIndices) path() string {
	return fmt.Sprintf("/instrument/activeAndIndices")
}

func (req *reqToGetActiveInstrumentsAndIndices) method() string {
	return http.MethodGet
}

func (req *reqToGetActiveInstrumentsAndIndices) query() (string, error) {
	return "", nil
}

func (req *reqToGetActiveInstrumentsAndIndices) payload() (string, error) {
	return "", nil
}

type reqToGetActiveIntervals struct {
	c *Client
}

type RespToGetActiveIntervals Intervals

func (req *reqToGetActiveIntervals) Do() (RespToGetActiveIntervals, error) {
	var response RespToGetActiveIntervals
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToGetActiveIntervals) path() string {
	return fmt.Sprintf("/instrument/activeIntervals")
}

func (req *reqToGetActiveIntervals) method() string {
	return http.MethodGet
}

func (req *reqToGetActiveIntervals) query() (string, error) {
	return "", nil
}

func (req *reqToGetActiveIntervals) payload() (string, error) {
	return "", nil
}

type reqToGetCompositeIndex struct {
	c         *Client
	symbol    string
	filter    map[string]interface{}
	columns   string
	count     int
	start     int
	reverse   optional.Bool
	startTime time.Time
	endTime   time.Time
}

type RespToGetCompositeIndex []CompositeIndex

func (req *reqToGetCompositeIndex) Symbol(symbol string) *reqToGetCompositeIndex {
	req.symbol = symbol
	return req
}

func (req *reqToGetCompositeIndex) Filter(filter map[string]interface{}) *reqToGetCompositeIndex {
	req.filter = filter
	return req
}

func (req *reqToGetCompositeIndex) Columns(columns ...string) *reqToGetCompositeIndex {
	req.columns = strings.Join(columns, ",")
	return req
}

func (req *reqToGetCompositeIndex) Count(count int) *reqToGetCompositeIndex {
	req.count = count
	return req
}

func (req *reqToGetCompositeIndex) Start(start int) *reqToGetCompositeIndex {
	req.start = start
	return req
}

func (req *reqToGetCompositeIndex) Reverse(reverse bool) *reqToGetCompositeIndex {
	req.reverse = optional.Bool{reverse}
	return req
}

func (req *reqToGetCompositeIndex) StartTime(startTime time.Time) *reqToGetCompositeIndex {
	req.startTime = startTime
	return req
}

func (req *reqToGetCompositeIndex) EndTime(endTime time.Time) *reqToGetCompositeIndex {
	req.endTime = endTime
	return req
}

func (req *reqToGetCompositeIndex) Do() (RespToGetCompositeIndex, error) {
	var response RespToGetCompositeIndex
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToGetCompositeIndex) path() string {
	return fmt.Sprintf("/instrument/compositeIndex")
}

func (req *reqToGetCompositeIndex) method() string {
	return http.MethodGet
}

func (req *reqToGetCompositeIndex) query() (string, error) {
	s, e := json.Marshal(req.filter)
	filterStr := string(s)
	if filterStr == "null" || e != nil {
		filterStr = ""
	}

	a := struct {
		C         *Client       `url:"-"`
		Symbol    string        `url:"symbol,omitempty"`
		Filter    string        `url:"filter,omitempty"`
		Columns   string        `url:"columns,omitempty"`
		Count     int           `url:"count,omitempty"`
		Start     int           `url:"start,omitempty"`
		Reverse   optional.Bool `url:"reverse,omitempty"`
		StartTime time.Time     `url:"startTime,omitempty"`
		EndTime   time.Time     `url:"endTime,omitempty"`
	}{
		req.c,
		req.symbol,
		filterStr,
		req.columns,
		req.count,
		req.start,
		req.reverse,
		req.startTime,
		req.endTime,
	}
	value, err := query.Values(&a)
	if err != nil {
		return "", err
	}

	return value.Encode(), nil
}

func (req *reqToGetCompositeIndex) payload() (string, error) {
	return "", nil
}

type reqToGetAllIndices struct {
	c *Client
}

type RespToGetAllIndices []Instrument

func (req *reqToGetAllIndices) Do() (RespToGetAllIndices, error) {
	var response RespToGetAllIndices
	err := req.c.request(req, &response)
	return response, err
}

func (req *reqToGetAllIndices) path() string {
	return fmt.Sprintf("/instrument/indices")
}

func (req *reqToGetAllIndices) method() string {
	return http.MethodGet
}

func (req *reqToGetAllIndices) query() (string, error) {
	return "", nil
}

func (req *reqToGetAllIndices) payload() (string, error) {
	return "", nil
}
