package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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
					flattened[fmt.Sprintf("%v\\.%v", key, nkey)] = nval
				}
			}
		} else {
			flattened[key] = value
		}
	}
	return flattened, nil
}

func nestJSON(flattened map[string]interface{}) (map[string]interface{}, error) {
	nested := make(map[string]interface{})
	for key, value := range flattened {
		if strings.Contains(key, "\\.") {
			splitKeys := strings.Split(key, "\\.")
			parentMap := &nested
			for i := 0; i < len(splitKeys); i++ {
				if i == len(splitKeys)-1 {
					(*parentMap)[splitKeys[i]] = value
				} else {
					if existing, found := (*parentMap)[splitKeys[i]]; found {
						worked, _ := existing.(map[string]interface{})
						parentMap = &worked
						continue
					}
					currentMap := make(map[string]interface{})
					(*parentMap)[splitKeys[i]] = currentMap
					parentMap = &currentMap
				}
			}
		} else {
			nested[key] = value
		}
	}
	return nested, nil
}
