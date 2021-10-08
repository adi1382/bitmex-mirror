package position

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToChangeLeverage struct {
	Symbol   string  `url:"symbol" json:"symbol"`
	Leverage float64 `url:"leverage" json:"leverage"`
}

type RespToChangeLeverage Position

func (req *ReqToChangeLeverage) Path() string {
	return fmt.Sprintf("/position/leverage")
}

func (req *ReqToChangeLeverage) Method() string {
	return http.MethodPost
}

func (req *ReqToChangeLeverage) Query() string {
	return ""
}

func (req *ReqToChangeLeverage) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
