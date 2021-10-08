package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToCancelAllOrders struct {
	Symbol string                 `url:"symbol,omitempty" json:"symbol,omitempty"`
	Filter map[string]interface{} `url:"filter,omitempty" json:"filter,omitempty"`
	Text   string                 `url:"text,omitempty" json:"text,omitempty"`
}

type RespToCancelAllOrders []Order

func (req *ReqToCancelAllOrders) Path() string {
	return fmt.Sprintf("/order/all")
}

func (req *ReqToCancelAllOrders) Method() string {
	return http.MethodDelete
}

func (req *ReqToCancelAllOrders) Query() string {
	return ""
}

func (req *ReqToCancelAllOrders) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
