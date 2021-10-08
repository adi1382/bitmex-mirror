package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToPlaceOrder struct {
	Symbol         string  `url:"symbol" json:"symbol"`
	Side           string  `url:"side,omitempty" json:"side,omitempty"`
	OrderQty       float64 `url:"orderQty,omitempty" json:"orderQty,omitempty"`
	Price          float64 `url:"price,omitempty" json:"price,omitempty"`
	DisplayQty     float64 `url:"displayQty,omitempty" json:"displayQty,omitempty"`
	StopPx         float64 `url:"stopPx,omitempty" json:"stopPx,omitempty"`
	ClOrdID        string  `url:"clOrdID,omitempty" json:"clOrdID,omitempty"`
	PegOffsetValue float64 `url:"pefOffsetValue,omitempty" json:"pefOffsetValue,omitempty"`
	PegPriceType   string  `url:"pegPriceType,omitempty" json:"pegPriceType,omitempty"`
	OrdTyp         string  `url:"ordType,omitempty" json:"ordType,omitempty"`
	TimeInForce    string  `url:"timeInForce,omitempty" json:"timeInForce,omitempty"`
	ExecInst       string  `url:"execInst,omitempty" json:"execInst,omitempty"`
	Text           string  `url:"text,omitempty" json:"text,omitempty"`
}

type RespToPlaceOrders Order

func (req *ReqToPlaceOrder) Path() string {
	return fmt.Sprintf("/order")
}

func (req *ReqToPlaceOrder) Method() string {
	return http.MethodPost
}

func (req *ReqToPlaceOrder) Query() string {
	return ""
}

func (req *ReqToPlaceOrder) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
