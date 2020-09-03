package utils

import (
	"encoding/json"
)

func JsonToMap(s []byte) (map[string]interface{}, error) {
	var raw map[string]interface{}

	if err := json.Unmarshal(s, &raw); err != nil {
		return nil, err
	}

	return raw, nil
}
