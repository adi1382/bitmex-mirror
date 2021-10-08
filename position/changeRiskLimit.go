package position

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToChangeRiskLimit struct {
	Symbol    string  `url:"symbol" json:"symbol"`
	RiskLimit float64 `url:"riskLimit" json:"riskLimit"`
}

type RespToChangeRiskLimit Position

func (req *ReqToChangeRiskLimit) Path() string {
	return fmt.Sprintf("/position/riskLimit")
}

func (req *ReqToChangeRiskLimit) Method() string {
	return http.MethodPost
}

func (req *ReqToChangeRiskLimit) Query() string {
	return ""
}

func (req *ReqToChangeRiskLimit) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
