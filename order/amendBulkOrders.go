package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToAmendBulkOrders struct {
	Orders []ReqToAmendOrder `url:"orders,omitempty" json:"orders,omitempty"`
}

type RespToAmendBulkOrders []Order

func (req *ReqToAmendBulkOrders) Path() string {
	return fmt.Sprintf("/order/bulk")
}

func (req *ReqToAmendBulkOrders) Method() string {
	return http.MethodPut
}

func (req *ReqToAmendBulkOrders) Query() string {
	return ""
}

func (req *ReqToAmendBulkOrders) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
