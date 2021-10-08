package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToCancelAllAfter struct {
	Timeout float64 `url:"timeout,omitempty" json:"timeout,omitempty"`
}

type RespToCancelAllAfter struct{}

func (req *ReqToCancelAllAfter) Path() string {
	return fmt.Sprintf("/order/cancelAllAfter")
}

func (req *ReqToCancelAllAfter) Method() string {
	return http.MethodPost
}

func (req *ReqToCancelAllAfter) Query() string {
	return ""
}

func (req *ReqToCancelAllAfter) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
