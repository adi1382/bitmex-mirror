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

func (req *ReqToChangeRiskLimit) path() string {
	return fmt.Sprintf("/position/riskLimit")
}

func (req *ReqToChangeRiskLimit) method() string {
	return http.MethodPost
}

func (req *ReqToChangeRiskLimit) query() string {
	return ""
}

func (req *ReqToChangeRiskLimit) payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
