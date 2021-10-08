package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToCancelOrders struct {
	OrderIDs []string `url:"ordID,omitempty" json:"orderID,omitempty"`
	ClOrdIDs []string `url:"clOrdID,omitempty" json:"clOrdID,omitempty"`
	Text     string   `url:"text,omitempty" json:"text,omitempty"`
}

type RespToCancelOrders Order

func (req *ReqToCancelOrders) Path() string {
	return fmt.Sprintf("/order")
}

func (req *ReqToCancelOrders) Method() string {
	return http.MethodDelete
}

func (req *ReqToCancelOrders) Query() string {
	return ""
}

func (req *ReqToCancelOrders) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
