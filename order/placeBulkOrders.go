package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToPlaceBulkOrder struct {
	Orders []ReqToPlaceOrder `url:"orders,omitempty" json:"orders,omitempty"`
}

type RespToPlaceBulkOrders []Order

func (req *ReqToPlaceBulkOrder) Path() string {
	return fmt.Sprintf("/order/bulk")
}

func (req *ReqToPlaceBulkOrder) Method() string {
	return http.MethodPost
}

func (req *ReqToPlaceBulkOrder) Query() string {
	return ""
}

func (req *ReqToPlaceBulkOrder) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
