package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/armon/go-radix"
	"github.com/elliotchance/orderedmap"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// =============================================================================
// Globals
// =============================================================================

// Apps is the array of App Indexes
var Apps map[string]*appIndexes

// =============================================================================
// Structs
// =============================================================================

type indexMap struct {
	field          string
	index          *radix.Tree
	TotalDocuments int
}

type appIndexes struct {
	Indexes        []indexMap
	Name           string `json:"Name"`
	TotalDocuments int    `json:"TotalDocuments"`
}

type fuzzyItem struct {
	key   string
	value interface{}
}

type listItem struct {
	IndexName   string
	IndexValues []string
}

type DocumentObject struct {
	Score float32                `json:"score"`
	Data  map[string]interface{} `json:"data"`
	DocID uint64                 `json:"docID"`
}

type SearchResponse struct {
	Items      []DocumentObject `json:"items"`
	SearchTime time.Duration    `json:"searchTimeNS"`
	ScoreTime  time.Duration    `json:"scoreTimeNS"`
	LoadTime   time.Duration    `json:"loadTimeNS"`
}

func initApp(name string) *appIndexes {
	appindex := appIndexes{make([]indexMap, 0), name, 0}
	if Apps == nil {
		Apps = make(map[string]*appIndexes, 0)
	}
	Apps[name] = &appindex
	return &appindex
}

func initIndexMap(indexmap *indexMap, name string) *indexMap {
	newMap := indexMap{name, radix.New(), 0}
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
		os.Mkdir("./documents", os.FileMode(0755))
	}
}

func LoadAppsFromDisk() (apps []*appIndexes) {
	start := time.Now()
	loadedApps := make([]*appIndexes, 0)
	if _, err := os.Stat("./apps"); os.IsNotExist(err) { // Make sure serialized folder exists
		os.Mkdir("./apps", os.FileMode(0755))
	}
	filepath.Walk("./apps", func(path string, info os.FileInfo, err error) error {
		appName := filepath.Base(path)
		if info == nil {
			log.Printf("Error: ./apps/%s does not exist\n", appName)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		appBytes, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		var newApp appIndexes
		json.Unmarshal(appBytes, &newApp)
		newApp.Indexes = make([]indexMap, 0)
		loadedApps = append(loadedApps, &newApp)
		return nil
	})
	end := time.Now()
	log.Printf("### Loaded serialized apps in %v\n", end.Sub(start))
	return loadedApps
}

func createOrderMapKey(docID uint64, score float32) string {
	return fmt.Sprintf("%v#%v", score, docID)
}

func parseOrderMapKey(compoundKey string) (docID uint64, score float32) {
	arr := strings.Split(compoundKey, "#")
	theScore, _ := strconv.ParseFloat(arr[0], 32)
	theDocID, _ := strconv.ParseUint(arr[1], 10, 64)
	return theDocID, float32(theScore)
}

func (app *appIndexes) LoadIndexesFromDisk() { // TODO: Change to search folders and load based on app
	start := time.Now()
	if _, err := os.Stat("./serialized"); os.IsNotExist(err) { // Make sure serialized folder exists
		os.Mkdir("./serialized", os.FileMode(0755))
	}
	if _, err := os.Stat(fmt.Sprintf("./serialized/%s", app.Name)); os.IsNotExist(err) { // Make sure app folder exists
		os.Mkdir(fmt.Sprintf("./serialized/%s", app.Name), os.FileMode(0755))
	}
	filepath.Walk(fmt.Sprintf("./serialized/%s", app.Name), func(path string, info os.FileInfo, err error) error {
		if info == nil {
			fmt.Printf("Error: ./serialized/%s does not exist\n", app.Name)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fieldName := filepath.Base(path)
		indexInd := len(app.Indexes)
		app.addIndexMap(fieldName)
		decodeFile, err := os.Open(path)
		defer decodeFile.Close()
		d := gob.NewDecoder(decodeFile)
		decoded := make(map[string]*orderedmap.OrderedMap)
		err = d.Decode(&decoded)

		converted := make(map[string]interface{})

		for key, value := range decoded {
			converted[key] = value
		}

		app.Indexes[indexInd].index = radix.NewFromMap(converted)
		return nil
	})
	end := time.Now()
	log.Printf("### Loaded serialized indexes in %v\n", end.Sub(start))
}

func fetchDocument(docID uint64) string {
	dat, _ := ioutil.ReadFile(fmt.Sprintf("./documents/%v", docID))
	return string(dat)
}

func loadDocuments(docObjs *SearchResponse) {
	for idx, docObj := range docObjs.Items {
		// Determine precision of document
		// Determine recall of document
		docString := fetchDocument(docObj.DocID)
		docObjs.Items[idx].Data, _ = parseArbJSON(docString)
	}
}

// ListApps lists the names of current apps
func ListApps() []string {
	appArr := make([]string, 0)
	for k, _ := range Apps {
		appArr = append(appArr, k)
	}
	return appArr
}

func GetApp(appName string) (*appIndexes, error) {
	app := Apps[appName]
	if app == nil {
		return nil, errors.New("App does not exist!")
	} else {
		return app, nil
	}
}

// ########################################################################
// ######################## appIndexes functions ##########################
// ########################################################################

func (appIndex *appIndexes) listIndexItems() []listItem {
	var output []listItem
	for _, i := range appIndex.Indexes {
		newItem := listItem{i.field, make([]string, 0)}
		i.index.Walk(func(k string, value interface{}) bool {
			newItem.IndexValues = append(newItem.IndexValues, k)
			return false
		})
		output = append(output, newItem)
	}
	return output
}

func (appIndex *appIndexes) listIndexes() []string {
	var output []string
	for _, i := range appIndex.Indexes {
		output = append(output, i.field)
	}
	return output
}

// addIndexMap creates a new index map (field) to be indexes
func (appindex *appIndexes) addIndexMap(name string) *indexMap {
	newIndexMap := indexMap{name, radix.New(), 0}
	oldLen := len(appindex.Indexes)
	appindex.Indexes = append(appindex.Indexes, newIndexMap)
	return &appindex.Indexes[oldLen]
}

func (appindex *appIndexes) addIndexFromDisk(parsed map[string]interface{}, filename string) (documentID uint64) {
	// log.Println("### Adding index...")
	// Format the input
	rand.Seed(time.Now().UnixNano())
	id, _ := strconv.ParseUint(filename, 10, 64)
	// log.Println("### ID:", id)
	var totalDocLen int
	for _, v := range parsed {
		totalDocLen += len(strings.FieldsFunc(v.(string), func(r rune) bool {
			return r == ' ' || r == '"'
		}))
	}
	for k, v := range parsed {
		// Don't index ID
		if strings.ToLower(k) == "docid" {
			continue
		}
		// Find if indexMap already exists
		var indexMapPointer *indexMap = nil
		for i := 0; i < len(appindex.Indexes); i++ {
			if k == appindex.Indexes[i].field {
				indexMapPointer = &appindex.Indexes[i]
				break
			}
		}

		if indexMapPointer == nil { // Create indexMap
			indexMapPointer = appindex.addIndexMap(k)
			// log.Println("### Creating new indexMap")
		}

		// Add index to indexMap
		indexMapPointer.addIndex(id, fmt.Sprintf("%v", v), totalDocLen)

	}
	// Write indexes document to disk
	fmt.Sprintf("Writing out to: %s\n", fmt.Sprintf("./documents/%v", id))
	sendback, _ := stringIndex(parsed)
	ioutil.WriteFile(fmt.Sprintf("./documents/%v", id), []byte(sendback), os.FileMode(0660))
	return id

	// TODO: Store document
	// TODO: check if tree exists with name of every json key, if not create tree

}

func (appindex *appIndexes) addIndex(parsed map[string]interface{}) (documentID uint64) {
	// log.Println("### Adding index...")
	// Format the input
	rand.Seed(time.Now().UnixNano())
	var id uint64
	if parsed["docID"] != nil {
		pre, _ := parsed["docID"].(json.Number).Int64()
		id = uint64(pre)
	} else {
		id = rand.Uint64()
	}
	// log.Println("### ID:", id)
	var totalDocLen int
	for _, v := range parsed {
		totalDocLen += len(strings.FieldsFunc(v.(string), func(r rune) bool {
			return r == ' ' || r == '"'
		}))
	}
	for k, v := range parsed {
		// Don't index ID
		if strings.ToLower(k) == "docid" {
			continue
		}
		// Find if indexMap already exists
		var indexMapPointer *indexMap
		for i := 0; i < len(appindex.Indexes); i++ {
			if k == appindex.Indexes[i].field {
				indexMapPointer = &appindex.Indexes[i]
				break
			}
		}

		if indexMapPointer == nil { // Create indexMap
			indexMapPointer = appindex.addIndexMap(k)
			// log.Println("### Creating new indexMap")
		}

		// Add index to indexMap
		indexMapPointer.addIndex(id, fmt.Sprintf("%v", v), totalDocLen)
	}
	// Remove docID field
	delete(parsed, "docID")
	appindex.TotalDocuments++ // Increase document count
	// Write indexes document to disk
	// fmt.Printf("Writing out to: %s\n", fmt.Sprintf("./documents/%v", id))
	sendback, _ := stringIndex(parsed)
	ioutil.WriteFile(fmt.Sprintf("./documents/%v", id), []byte(sendback), os.FileMode(0660))
	return id
}

func (appindex *appIndexes) search(input string, fields []string, bw bool) (documentIDs []uint64, response SearchResponse) {
	output := make(map[uint64][]float32, 0) // docID: [score, docLen]
	var avgDocLen float32
	// Tokenize input
	start := time.Now()
	searchTokens := lowercaseTokens(tokenizeString(input))
	for _, token := range searchTokens { // Initial scoring pass, covers TF weighting and IDF weighting
		// Check fields
		if len(fields) == 0 { // check all
			// log.Println("### No fields given, searching all fields...")
			for _, indexmap := range appindex.Indexes {
				// log.Println("### Searching index:", indexmap.field, "for", token)
				if bw {
					searchItems := indexmap.beginsWithSearch(token)
					for _, searchItem := range searchItems {
						if searchItem != nil {
							// FIXME: This is a hacky solution for getting number of docs with the term in it
							numDocsWithTerm := 0
							for el := searchItem.Front(); el != nil || numDocsWithTerm >= 100; el = el.Next() {
								numDocsWithTerm++
							}
							// Only first 100 items in orderedmap
							iterCount := 0
							for el := searchItem.Front(); el != nil || iterCount >= 100; el = el.Next() {
								docID, tfWeighting := parseOrderMapKey(el.Key.(string))
								idfWeighting := math.Log(float64(1+indexmap.TotalDocuments) / float64(numDocsWithTerm))
								fieldLen := float32(el.Value.(int))
								newVal := make([]float32, 2)
								if val, ok := output[docID]; ok { // Exists
									prevScore := val[0]
									prevLen := val[1]
									newVal[0] = prevScore + (tfWeighting * float32(idfWeighting))
									newVal[1] = prevLen + fieldLen
									avgDocLen += fieldLen
									output[docID] = newVal
								} else {
									newVal[0] = tfWeighting
									newVal[1] = fieldLen
									avgDocLen += fieldLen
									output[docID] = newVal
								}
								iterCount++
							}
						}
					}
				} else {
					searchItems := indexmap.search(token)
					if searchItems != nil {
						// FIXME: This is a hacky solution for getting number of docs with the term in it
						numDocsWithTerm := 0
						for el := searchItems.Front(); el != nil || numDocsWithTerm >= 100; el = el.Next() {
							numDocsWithTerm++
						}
						// Only first 100 items in orderedmap
						iterCount := 0
						for el := searchItems.Front(); el != nil || iterCount >= 100; el = el.Next() {
							docID, tfWeighting := parseOrderMapKey(el.Key.(string))
							idfWeighting := math.Log(float64(1+indexmap.TotalDocuments) / float64(numDocsWithTerm))
							fieldLen := float32(el.Value.(int))
							newVal := make([]float32, 2)
							if val, ok := output[docID]; ok { // Exists
								prevScore := val[0]
								prevLen := val[1]
								newVal[0] = prevScore + (tfWeighting * float32(idfWeighting))
								newVal[1] = prevLen
								output[docID] = newVal
							} else {
								newVal[0] = tfWeighting
								newVal[1] = fieldLen
								fmt.Println()
								avgDocLen += fieldLen
								output[docID] = newVal
							}
							iterCount++
						}
					}
				}
			}
		} else { // check given fields
			for _, field := range fields {
				// log.Println("### Searching index:", field, "for", token)
				searchItems := appindex.searchByField(token, field, bw)
				for _, searchItem := range searchItems {
					if searchItem != nil {
						// FIXME: This is a hacky solution for getting number of docs with the term in it
						numDocsWithTerm := 0
						for el := searchItem.Front(); el != nil || numDocsWithTerm >= 100; el = el.Next() {
							numDocsWithTerm++
						}
						// Only first 100 items in orderedmap
						iterCount := 0
						for el := searchItem.Front(); el != nil || iterCount >= 100; el = el.Next() {
							docID, tfWeighting := parseOrderMapKey(el.Key.(string))
							// FIXME: Currently does not cover the IDF weighting
							fieldLen := float32(el.Value.(int))
							newVal := make([]float32, 2)
							if val, ok := output[docID]; ok { // Exists
								prevScore := val[0]
								prevLen := val[1]
								newVal[0] = prevScore + tfWeighting
								newVal[1] = prevLen + float32(fieldLen)
								avgDocLen += float32(fieldLen)
								output[docID] = newVal
							} else {
								newVal[0] = tfWeighting
								newVal[1] = float32(fieldLen)
								avgDocLen += float32(fieldLen)
								output[docID] = newVal
							}
							iterCount++
						}
					}
				}
			}
		}
	}
	responseObj := SearchResponse{
		Items: make([]DocumentObject, 0),
	}
	end := time.Now()
	diff := end.Sub(start)
	responseObj.SearchTime = diff

	// -------------------
	// Scoring and Sorting
	// -------------------

	start = time.Now()

	// Calculate avg doc len
	avgDocLen = avgDocLen / float32(len(output)) // Get average

	// Now we have docID: [score, docLen]
	for docID, scoreLen := range output {
		termsInQueryLen := len(strings.Fields(input))
		fmt.Println(docID, scoreLen[1], avgDocLen)
		finalDocScore := (scoreLen[0] / (scoreLen[1] / avgDocLen)) * float32(termsInQueryLen)
		responseObj.Items = append(responseObj.Items, DocumentObject{
			Data:  nil,
			Score: finalDocScore,
			DocID: docID,
		})
	}
	sort.Slice(responseObj.Items, func(i int, j int) bool {
		return responseObj.Items[i].Score > responseObj.Items[j].Score
	})
	// TODO: Only take first 100 documents for now
	if len(responseObj.Items) > 100 {
		responseObj.Items = responseObj.Items[:100]
	}
	end = time.Now()
	responseObj.ScoreTime = end.Sub(start)

	// -----------------
	// Loading Documents
	// -----------------

	start = time.Now()
	loadDocuments(&responseObj)
	end = time.Now()
	responseObj.LoadTime = end.Sub(start)
	yeye := make([]uint64, 0) // FIXME: This is a temp solution obviously
	return yeye, responseObj
}

func (appindex *appIndexes) searchByField(input string, field string, bw bool) (documents []*orderedmap.OrderedMap) {
	// Check if field exists
	output := make([]*orderedmap.OrderedMap, 0)
	for _, indexmap := range appindex.Indexes {
		if indexmap.field == field {
			if bw {
				searchItems := indexmap.beginsWithSearch(input)
				output = append(output, searchItems...)
			} else {
				searchItems := indexmap.search(input)
				output = append(output, searchItems)
			}
			break
		}
	}
	return output
}

func (appindex *appIndexes) SerializeApp() {
	log.Printf("### Serializing App: %s\n", appindex.Name)
	if _, err := os.Stat("./apps"); os.IsNotExist(err) { // Make sure apps folder exists
		os.Mkdir("./apps", os.FileMode(0755))
	}
	// Wipe tree since serializing tree elsewhere
	// appindex.Indexes = nil
	serializedApp, err := json.Marshal(appindex)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(fmt.Sprintf("./apps/%s", appindex.Name), serializedApp, os.FileMode(0755))
	log.Printf("### Successfully Serialized App %s!\n", appindex.Name)
}

func (appindex *appIndexes) SerializeIndex() {
	log.Printf("### Serializing %s Indexes...\n", appindex.Name)
	if _, err := os.Stat("./serialized"); os.IsNotExist(err) { // Make sure serialized folder exists
		os.Mkdir("./serialized", os.FileMode(0755))
	}
	if _, err := os.Stat(fmt.Sprintf("./serialized/%s", appindex.Name)); os.IsNotExist(err) { // Make sure app folder exists
		os.Mkdir(fmt.Sprintf("./serialized/%s", appindex.Name), os.FileMode(0755))
	}
	for _, i := range appindex.Indexes {
		serializedTree := i.index.ToMap()
		encodeFile, err := os.Create(fmt.Sprintf("./serialized/%s/%s", appindex.Name, i.field))
		if err != nil {
			panic(err)
		}
		e := gob.NewEncoder(encodeFile)
		converted := make(map[string]*orderedmap.OrderedMap)
		for key, value := range serializedTree {
			converted[key] = value.(*orderedmap.OrderedMap)
		}
		err = e.Encode(converted)
		encodeFile.Close()
	}
	log.Printf("### Successfully Serialized %s Indexes!\n", appindex.Name)
}

func (appindex *appIndexes) deleteIndex(docID uint64) error {
	var err error = nil
	// Find document on disk
	docPath := fmt.Sprintf("./documents/%v", docID)
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		log.Printf("DocID: %v does not exist\n", docID)
		return err
	}
	docData, err := ioutil.ReadFile(docPath)
	docString := string(docData)
	docJSON, err := parseArbJSON(docString)
	if err != nil {
		return err
	}
	for k, v := range docJSON {
		// Don't index ID
		input := fmt.Sprintf("%v", v)
		if k == "docID" {
			continue
		}
		// Find the field index
		for i := 0; i < len(appindex.Indexes); i++ {
			if k == appindex.Indexes[i].field {
				indexmap := appindex.Indexes[i]
				// Tokenize field
				for _, token := range lowercaseTokens(tokenizeString(input)) {
					tokenMap, _ := indexmap.index.Get(token)
					if tokenMap == nil { // somehow not indexed
						log.Printf("### Error: field not indexes from disk document (deleteIndex)\n")
						continue
					} else { // update node
						docMap := tokenMap.(*orderedmap.OrderedMap)
						documentTermScore := calculateTokenScoreByField(input, token, indexmap.TotalDocuments, docMap.Len())
						deleted := docMap.Delete(createOrderMapKey(docID, documentTermScore))
						if deleted != true {
							log.Printf("### Error: Document entry not deleted for token %s\n", token)
						}
						if docMap.Len() == 0 {
							_, deleted := indexmap.index.Delete(token) // this should delete the node if it is empty
							if !deleted {
								log.Printf("Node was not deleted\n")
							}
						} else {
							_, updated := indexmap.index.Insert(token, docMap)
							if !updated {
								log.Printf("### SOMEHOW DIDN'T UPDATE WHEN DELETING INDEX ###\n")
								continue
							}
						}
					}
				}
				indexmap.TotalDocuments--
				break
			}
		}
	}
	appindex.TotalDocuments--
	// Delete document on disk
	err = os.Remove(docPath)
	if err != nil {
		return err
	}
	return nil
}

func (appindex *appIndexes) updateIndex(parsed map[string]interface{}) (errrr error, created bool) {
	var err error = nil
	// Validate docID
	var docID uint64
	pre := parsed["docID"].(json.Number)
	docID, err = strconv.ParseUint(pre.String(), 10, 64)
	if err != nil {
		return err, false
	}
	err = appindex.deleteIndex(docID)
	if os.IsNotExist(err) {
		appindex.addIndex(parsed)
		err = nil
		return err, true
	} else if err == nil {
		appindex.addIndex(parsed)
	}
	return err, false
}

func FuzzySearch(key string, t *radix.Tree) []fuzzyItem {
	output := make([]fuzzyItem, 0)
	split := strings.Split(key, "")
	t.WalkPrefix(split[0], func(k string, value interface{}) bool {
		// log.Println("### Checking", k)
		includesAll := true
		for _, char := range split {
			// log.Println("### DOES", k, "include", char)
			if !strings.Contains(k, char) {
				// log.Println("NOPE")
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

func calculateTokenScoreByField(fieldValue string, tokenValue string, totalDocuments int, orderMapLen int) (tokenScore float32) {
	// BEGIN Precision word match
	tokenPrecisionFreq := 0
	docStringWord := strings.FieldsFunc(fieldValue, func(r rune) bool {
		return r == ' ' || r == '"'
	})
	for _, word := range docStringWord {
		if strings.ToLower(word) == tokenValue {
			tokenPrecisionFreq++
		}
	}
	termFreqWeight := 1 + (math.Log(1 + math.Log(1+float64(tokenPrecisionFreq))))
	// END Precision word match
	return float32(termFreqWeight)
}

func (indexmap *indexMap) addIndex(id uint64, value string, docLen int) {
	CheckDocumentsFolder()
	// Tokenize
	for _, token := range lowercaseTokens(tokenizeString(value)) {
		// log.Println("### INDEXING:", token)
		// Check if index already exists
		prenode, _ := indexmap.index.Get(token)
		var ids *orderedmap.OrderedMap
		if prenode == nil { // create new node
			ids = orderedmap.NewOrderedMap()
			// Calculate token score here
			documentTermScore := calculateTokenScoreByField(value, token, indexmap.TotalDocuments, ids.Len())
			ids.Set(createOrderMapKey(id, documentTermScore), docLen)
			_, updated := indexmap.index.Insert(token, ids)
			if updated {
				fmt.Errorf("### SOMEHOW UPDATED WHEN INSERTING NEW ###\n")
			}
			// log.Printf("### ADDED %v WITH IDS: %v ###\n", token, ids)
		} else { // update node
			node := prenode.(*orderedmap.OrderedMap)
			ids = node
			documentTermScore := calculateTokenScoreByField(value, token, indexmap.TotalDocuments, ids.Len())
			node.Set(createOrderMapKey(id, documentTermScore), docLen)
			_, updated := indexmap.index.Insert(token, node)
			if !updated {
				fmt.Errorf("### SOMEHOW DIDN'T UPDATE WHEN UPDATING INDEX ###\n")
			}
			// log.Printf("### UPDATED %v WITH IDS: %v ###\n", token, ids)
		}
	}
	indexmap.TotalDocuments++
}

func (indexmap *indexMap) search(input string) (documents *orderedmap.OrderedMap) {
	search, _ := indexmap.index.Get(input)
	if search == nil {
		return nil
	}
	node := search.(*orderedmap.OrderedMap)
	return node
}

func (indexmap *indexMap) beginsWithSearch(input string) (documents []*orderedmap.OrderedMap) {
	output := make([]*orderedmap.OrderedMap, 0)
	count := 0
	indexmap.index.WalkPrefix(input, func(key string, value interface{}) bool {
		if count >= 100 {
			return true
		}
		node := value.(*orderedmap.OrderedMap)
		output = append(output, node)
		return false
	})
	return output
}
