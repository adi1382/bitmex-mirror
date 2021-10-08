package position

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToTransferMargin struct {
	Symbol string  `url:"symbol" json:"symbol"`
	Amount float64 `url:"amount" json:"amount"`
}

type RespToTransferMargin Position

func (req *ReqToTransferMargin) Path() string {
	return fmt.Sprintf("/position/transferMargin")
}

func (req *ReqToTransferMargin) Method() string {
	return http.MethodPost
}

func (req *ReqToTransferMargin) Query() string {
	return ""
}

func (req *ReqToTransferMargin) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
