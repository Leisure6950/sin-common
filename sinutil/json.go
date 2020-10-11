package sinutil

import (
	"encoding/json"
)

func IsValidJson(str string) bool {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(str), &js)
	return err == nil
}
