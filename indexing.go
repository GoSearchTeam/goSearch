package main

import (
	"fmt"
	"github.com/segmentio/ksuid"
	"strings"
)

type indexMap struct {
	field string
	index map[string][]string
}

type appIndexes struct {
	indexes []indexMap
	name    string
}

func initApp(name string) *appIndexes {
	appindex := appIndexes{make([]indexMap, 0), name}
	return &appindex
}

func initIndexMap(indexmap *indexMap, name string) *indexMap {
	newMap := indexMap{name, make(map[string][]string)}
	return &newMap
}

func tokenizeString(input string) []string {
	return strings.Fields(input)
}

func lowercaseTokens(tokens []string) []string {
	output := make([]string, len(tokens))
	for i, token := range tokens {
		output[i] = strings.ToLower(token)
	}
	return output
}

// ########################################################################
// ######################## appIndexes functions ##########################
// ########################################################################

func (appindex *appIndexes) addIndexMap(name string) *indexMap {
	newIndexMap := indexMap{name, make(map[string][]string)}
	appindex.indexes = append(appindex.indexes, newIndexMap)
	return &newIndexMap
}

func (appindex *appIndexes) addIndex(parsed map[string]interface{}) {
	fmt.Println("### Adding index...")
	// Format the input
	var id string = fmt.Sprintf("%v", parsed["id"])
	if parsed["id"] == nil {
		fmt.Println("### No id found")
		id = ksuid.New().String()
	}
	fmt.Println("### ID:", id)
	for k, v := range parsed {
		// Don't index ID
		if strings.ToLower(k) == "id" {
			continue
		}
		// Find if indexMap already exists
		var indexMapPointer *indexMap = nil
		for i := 0; i < len(appindex.indexes); i++ {
			if k == appindex.indexes[i].field {
				indexMapPointer = &appindex.indexes[i]
				break
			}
		}

		if indexMapPointer == nil { // Create indexMap
			indexMapPointer = appindex.addIndexMap(k)
			fmt.Println("### Creating new indexMap")
		}

		// Add index to indexMap
		indexMapPointer.addIndex(id, fmt.Sprintf("%v", v))
	}

	// TODO: Store document
	// TODO: check if tree exists with name of every json key, if not create tree

}

// ########################################################################
// ######################### indexMap functions ###########################
// ########################################################################

func (indexmap *indexMap) addIndex(id string, value string) {
	// Tokenize
	for _, token := range lowercaseTokens(tokenizeString(value)) {
		fmt.Println("### INDEXING:", token)
		// Check if index already exists
		if indexmap.index[token] != nil {
			var found bool = false
			for _, docID := range indexmap.index[token] {
				fmt.Println("### Found token, Checking if doc exists...")
				if docID == id {
					fmt.Println("### Skip to avoid duplicates")
					found = true
					break
				}
			}
			if found {
				continue
			}
		}
		indexmap.index[token] = append(indexmap.index[token], id)
	}
}
