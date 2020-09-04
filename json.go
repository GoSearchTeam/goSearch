package main

import (
	"encoding/json"
)

func parseArbJSON(body string) (parsed map[string]interface{}, errr error) {
	b := []byte(body)
	jsonMap := make(map[string](interface{}))
	err := json.Unmarshal([]byte(b), &jsonMap)
	return jsonMap, err
}

func stringIndex(parsed interface{}) (stringed string, err error) {
	jsonObj, e := json.Marshal(parsed)
	return string(jsonObj), e
}
