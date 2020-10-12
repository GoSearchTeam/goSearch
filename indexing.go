package main

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/armon/go-radix"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type indexMap struct {
	field string
	index *radix.Tree
}

type appIndexes struct {
	indexes []indexMap
	name    string
}

type fuzzyItem struct {
	key   string
	value interface{}
}

type listItem struct {
	IndexName   string
	IndexValues []string
}

func initApp(name string) *appIndexes {
	appindex := appIndexes{make([]indexMap, 0), name}
	return &appindex
}

func initIndexMap(indexmap *indexMap, name string) *indexMap {
	newMap := indexMap{name, radix.New()}
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

func CheckDocumentsFolder() {
	if _, err := os.Stat("./documents"); os.IsNotExist(err) {
		os.Mkdir("./documents", os.ModePerm)
	}
}

func LoadIndexesFromDisk(app *appIndexes) {
	// TODO: Change to load in serialized
	files := make([]string, 0)
	start := time.Now()
	filepath.Walk("./documents", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		// Load document into index
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		doc, _ := parseArbJSON(string(dat))
		app.addIndex(doc)
		return nil
	})
	end := time.Now()
	fmt.Printf("### Loaded %d file(s) in %v ###\n", len(files), end.Sub(start))
}

// ########################################################################
// ######################## appIndexes functions ##########################
// ########################################################################

func (appIndex *appIndexes) listIndexItems() []listItem {
	var output []listItem
	for _, i := range appIndex.indexes {
		newItem := listItem{i.field, make([]string, 0)}
		i.index.Walk(func(k string, value interface{}) bool {
			// v := value.(roaring.Bitmap)
			fmt.Println(k)
			newItem.IndexValues = append(newItem.IndexValues, k)
			return false
		})
		output = append(output, newItem)
	}
	return output
}

func (appIndex *appIndexes) listIndexes() []string {
	var output []string
	for _, i := range appIndex.indexes {
		output = append(output, i.field)
	}
	return output
}

func (appindex *appIndexes) addIndexMap(name string) *indexMap {
	newIndexMap := indexMap{name, radix.New()}
	appindex.indexes = append(appindex.indexes, newIndexMap)
	return &newIndexMap
}

func (appindex *appIndexes) addIndex(parsed map[string]interface{}) (documentID uint32) {
	// fmt.Println("### Adding index...")
	// Format the input
	id := rand.Uint32()
	// fmt.Println("### ID:", id)
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
			// fmt.Println("### Creating new indexMap")
		}

		// Add index to indexMap
		indexMapPointer.addIndex(id, fmt.Sprintf("%v", v))
	}
	// Write indexes document to disk
	fmt.Sprintf("Writing out to: %s\n", fmt.Sprintf("./documents/%v", id))
	sendback, _ := stringIndex(parsed)
	ioutil.WriteFile(fmt.Sprintf("./documents/%v", id), []byte(sendback), os.ModePerm)
	return id

	// TODO: Store document
	// TODO: check if tree exists with name of every json key, if not create tree

}

func (appindex *appIndexes) search(input string, fields []string) (documentIDs []uint32) {
	var output []uint32
	// Tokenize input
	for _, token := range lowercaseTokens(tokenizeString(input)) {
		// Check fields
		if len(fields) == 0 { // check all
			// fmt.Println("### No fields given, searching all fields...")
			for _, indexmap := range appindex.indexes {
				// fmt.Println("### Searching index:", indexmap.field, "for", token)
				output = append(output, indexmap.search(token)...)
			}
		} else { // check given fields
			for _, field := range fields {
				// fmt.Println("### Searching index:", field, "for", token)
				docIDs := appindex.searchByField(token, field)
				if docIDs == nil {
					// fmt.Println("### Field doesn't exist:", field)
					continue
				}
				output = append(output, docIDs...)
			}
		}
	}
	return output
}

func (appindex *appIndexes) searchByField(input string, field string) (documentIDs []uint32) {
	// Check if field exists
	var output []uint32
	for _, indexmap := range appindex.indexes {
		if indexmap.field == field {
			output = append(output, indexmap.search(input)...)
			break
		}
	}
	return output
}

// FuzzySearch performs a fuzzy search on a given tree - WARNING: FAR LESS EFFICIENT
func FuzzySearch(key string, t *radix.Tree) []fuzzyItem {
	output := make([]fuzzyItem, 0)
	split := strings.Split(key, "")
	t.WalkPrefix(split[0], func(k string, value interface{}) bool {
		// fmt.Println("### Checking", k)
		includesAll := true
		for _, char := range split {
			// fmt.Println("### DOES", k, "include", char)
			if !strings.Contains(k, char) {
				// fmt.Println("NOPE")
				includesAll = false
				break
			}
		}
		if includesAll {
			output = append(output, fuzzyItem{k, value})
		}
		return false
	})
	return output
}

// ########################################################################
// ######################### indexMap functions ###########################
// ########################################################################

func (indexmap *indexMap) addIndex(id uint32, value string) {
	CheckDocumentsFolder()
	// Tokenize
	for _, token := range lowercaseTokens(tokenizeString(value)) {
		// fmt.Println("### INDEXING:", token)
		// Check if index already exists
		prenode, _ := indexmap.index.Get(token)
		var ids *roaring.Bitmap
		if prenode == nil { // create new node
			ids = roaring.BitmapOf(id)
			_, updated := indexmap.index.Insert(token, ids)
			if updated {
				fmt.Errorf("### SOMEHOW UPDATED WHEN INSERTING NEW ###\n")
			}
			// fmt.Printf("### ADDED %v WITH IDS: %v ###\n", token, ids)
		} else { // update node
			node := prenode.(*roaring.Bitmap)
			newid := rand.Uint32()
			node.Add(newid)
			ids = node
			_, updated := indexmap.index.Insert(token, ids)
			if !updated {
				fmt.Errorf("### SOMEHOW DIDN'T UPDATE WHEN UPDATING INDEX ###\n")
			}
			// fmt.Printf("### UPDATED %v WITH IDS: %v ###\n", token, ids)
		}
		// if indexmap.index[token] != nil {
		// 	var found bool = false
		// 	for _, docID := range indexmap.index[token] {
		// 		// fmt.Println("### Found token, Checking if doc exists...")
		// 		if docID == id {
		// 			// fmt.Println("### Skip to avoid duplicates")
		// 			found = true
		// 			break
		// 		}
		// 	}
		// 	if found {
		// 		continue
		// 	}
		// }
		// indexmap.index[token] = append(indexmap.index[token], id)
	}
}

func (indexmap *indexMap) search(input string) (documentIDs []uint32) {
	var output []uint32
	search, _ := indexmap.index.Get(input)
	if search == nil {
		return output
	}
	node := search.(*roaring.Bitmap)
	output = append(output, node.ToArray()...)
	return output
}
