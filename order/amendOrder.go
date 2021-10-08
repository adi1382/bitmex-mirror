package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToAmendOrder struct {
	OrderID        string  `url:"orderID,omitempty" json:"orderID,omitempty"`
	OrigClOrdID    string  `url:"origClOrdID,omitempty" json:"origClOrdID,omitempty"`
	ClOrdID        string  `url:"clOrdID,omitempty" json:"clOrdID,omitempty"`
	OrderQty       float64 `url:"orderQty,omitempty" json:"orderQty,omitempty"`
	LeavesQty      float64 `url:"leavesQty,omitempty" json:"leavesQty,omitempty"`
	Price          float64 `url:"price,omitempty" json:"price,omitempty"`
	StopPx         float64 `url:"stopPx,omitempty" json:"stopPx,omitempty"`
	PegOffsetValue float64 `url:"pefOffsetValue,omitempty" json:"pefOffsetValue,omitempty"`
	Text           string  `url:"text,omitempty" json:"text,omitempty"`
}

type RespToAmendOrder Order

func (req *ReqToAmendOrder) Path() string {
	return fmt.Sprintf("/order")
}

func (req *ReqToAmendOrder) Method() string {
	return http.MethodPut
}

func (req *ReqToAmendOrder) Query() string {
	return ""
}

func (req *ReqToAmendOrder) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
