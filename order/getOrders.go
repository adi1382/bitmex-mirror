package order

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

type Order struct {
	OrderID               string    `json:"orderID"`
	ClOrdID               string    `json:"clOrdID"`
	ClOrdLinkID           string    `json:"clOrdLinkID"`
	Account               int       `json:"account"`
	Symbol                string    `json:"symbol"`
	Side                  string    `json:"side"`
	SimpleOrderQty        int       `json:"simpleOrderQty"`
	OrderQty              int       `json:"orderQty"`
	Price                 int       `json:"price"`
	DisplayQty            int       `json:"displayQty"`
	StopPx                int       `json:"stopPx"`
	PegOffsetValue        int       `json:"pegOffsetValue"`
	PegPriceType          string    `json:"pegPriceType"`
	Currency              string    `json:"currency"`
	SettlCurrency         string    `json:"settlCurrency"`
	OrdType               string    `json:"ordType"`
	TimeInForce           string    `json:"timeInForce"`
	ExecInst              string    `json:"execInst"`
	ContingencyType       string    `json:"contingencyType"`
	ExDestination         string    `json:"exDestination"`
	OrdStatus             string    `json:"ordStatus"`
	Triggered             string    `json:"triggered"`
	WorkingIndicator      bool      `json:"workingIndicator"`
	OrdRejReason          string    `json:"ordRejReason"`
	SimpleLeavesQty       int       `json:"simpleLeavesQty"`
	LeavesQty             int       `json:"leavesQty"`
	SimpleCumQty          int       `json:"simpleCumQty"`
	CumQty                int       `json:"cumQty"`
	AvgPx                 int       `json:"avgPx"`
	MultiLegReportingType string    `json:"multiLegReportingType"`
	Text                  string    `json:"text"`
	TransactTime          time.Time `json:"transactTime"`
	Timestamp             time.Time `json:"timestamp"`
}

type ReqToGetOrders struct {
	Symbol    string            `url:"symbol,omitempty"`
	Filter    map[string]string `url:"filter,omitempty"`
	Columns   string            `url:"columns,omitempty"`
	Count     int               `url:"count,omitempty"`
	Start     int               `url:"start,omitempty"`
	Reverse   bool              `url:"reverse"`
	StartTime time.Time         `url:"startTime,omitempty"`
	EndTime   time.Time         `url:"endTime,omitempty"`
}

type RespToGetOrders []Order

func (req *ReqToGetOrders) Path() string {
	return fmt.Sprintf("/order")
}

func (req *ReqToGetOrders) Method() string {
	return http.MethodGet
}

func (req *ReqToGetOrders) Query() string {
	value, _ := query.Values(req)
	return value.Encode()
}

func (req *ReqToGetOrders) Payload() string {
	return ""
}
