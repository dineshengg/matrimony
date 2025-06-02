package utility

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONErr string

func (e JSONErr) Error() string {
	return string(e)
}

func ParseJSONPayload(r *http.Request) (map[string]interface{}, error) {
	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON payload: %v", err)
	}
	return payload, nil
}
