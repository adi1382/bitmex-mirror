package position

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqToIsolateMargin struct {
	Symbol  string `url:"symbol" json:"symbol"`
	Enabled bool   `url:"enabled" json:"enabled"`
}

type RespToIsolateMargin Position

func (req *ReqToIsolateMargin) Path() string {
	return fmt.Sprintf("/position/isolate")
}

func (req *ReqToIsolateMargin) Method() string {
	return http.MethodPost
}

func (req *ReqToIsolateMargin) Query() string {
	return ""
}

func (req *ReqToIsolateMargin) Payload() string {
	b, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(b)
}
