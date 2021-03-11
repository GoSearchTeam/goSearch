package main

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func flattenJSON(original map[string]interface{}) (map[string]interface{}, error) {
	flattened := make(map[string]interface{})
	for key, value := range original {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			v, ok := value.(map[string]interface{})
			if !ok {
				panic("Failed to convert map to map[string]{interface}!")
			} else {
				flat, _ := flattenJSON(v)
				for nkey, nval := range flat {
					flattened[fmt.Sprintf("%v_%v", key, nkey)] = nval
				}
			}
		} else {
			flattened[key] = value
		}
	}
	return flattened, nil
}
