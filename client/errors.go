package client

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	StatusCode int
	Message    string
	Name       string
}

// Error return the error Message
func (e APIError) Error() string {
	return fmt.Sprintf("APIError: StatusCode=%d, Message=%s, Name=%s", e.StatusCode, e.Message, e.Name)
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	var v map[string]map[string]string

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	e.Message = v["error"]["message"]
	e.Name = v["error"]["name"]

	return nil
}
