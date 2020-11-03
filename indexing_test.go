package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestAddIndex(t *testing.T) {
	app := initApp("TestAddIndex")
	data, _ := parseArbJSON(`{
		"index": "Test Add"
	}`)
	id := app.addIndex(data)

	doc := fmt.Sprintf("./documents/%d", id)
	if _, err := os.Stat(doc); os.IsNotExist(err) {
		t.Errorf("%s not created", doc)
	}
	ids := []uint64{id}
	documentsCleanup(ids)
}

func TestDeleteIndex(t *testing.T) {
	app := initApp("TestDeleteIndex")
	data, _ := parseArbJSON(`{
		"index":"Test Delete"	
	}`)

	docID := app.addIndex(data)
	doc := fmt.Sprintf("./documents/%d", docID)
	err := app.deleteIndex(docID)

	if err != nil {
		t.Errorf("Delete failed:\n%v", err)
	}
	if _, err := os.Stat(doc); os.IsExist(err) {
		t.Error("Document still exists on disk")
	}
}

func TestSmallSearch(t *testing.T) {
	app := initApp("TestSmallSearch")
	fields := []string{"name"}
	tacoData, _ := parseArbJSON(`{
		"name": "tacos"
	}`)
	tacoID := app.addIndex(tacoData)

	burgerData, _ := parseArbJSON(`{
		"name": "burger"
	}`)
	burgerID := app.addIndex(burgerData)

	tacoSearchID, tacoSearchData := app.search("tacos", fields, false)
	burgerSearchID, burgerSearchData := app.search("burger", fields, false)

	documentsCleanup([]uint64{tacoID, burgerID})

	if len(tacoSearchID) != 1 || len(tacoSearchData) != 1 || len(burgerSearchID) != 1 || len(burgerSearchData) != 1 {
		t.Error("Returned wrong number of restults")
	}
	if tacoSearchID[0] != tacoID || burgerSearchID[0] != burgerID {
		t.Error("ID's didnt match")
	}
	if tacoSearchData[0] != `{"name":"tacos"}` || burgerSearchData[0] != `{"name":"burger"}` {
		t.Error("Documents didnt match")
	}
}

func TestUpdateIndex(t *testing.T) {
	app := initApp("TestUpdateIndex")
	data, _ := parseArbJSON(`{
		"index":"Test Update"
	}`)

	docID := app.addIndex(data)

	d := json.NewDecoder(strings.NewReader(fmt.Sprintf(`{
		"docID":%v,
		"index2":"Test Update 2"
	}`, docID)))
	d.UseNumber()
	var uData map[string]interface{}
	if err := d.Decode(&uData); err != nil {
		t.Error(err)
	}

	app.updateIndex(uData)
	updatedID, updatedData := app.search("Test Update 2", []string{"index2"}, false)

	if len(updatedID) != 1 || len(updatedData) != 1 {
		t.Error("Search returned wrong number of results")
	}
	if updatedID[0] != docID {
		t.Error("IDs dont match")
	}
	if updatedData[0] != `{"index2":"value2"}` {
		t.Error("Data doesnt match")
	}

	documentsCleanup([]uint64{docID})
}

func TestFullCRUD(t *testing.T) {
	app := initApp("TestFullCRUD")
	data, _ := parseArbJSON(`{
		"index":"Full CRUD"
	}`)

	docID := app.addIndex(data)
	doc := fmt.Sprintf("./documents/%d", docID)
	if _, err := os.Stat(doc); os.IsNotExist(err) {
		t.Error("Document not created")
	}

	searchedID, searchedData := app.search("Full CRUD", []string{"index"}, false)

	if len(searchedID) != 1 || len(searchedData) != 1 {
		t.Error("SEARCH: returned incorrect number of results")
	}
	if searchedID[0] != docID {
		t.Error("SEARCH: IDs dont match")
	}
	if searchedData[0] != `{"index":"Full CRUD"}` {
		t.Error("SEARCH: Data doesnt match")
	}

	d := json.NewDecoder(strings.NewReader(fmt.Sprintf(`{
		"docID":%v,
		"index2":"Full CRUD 2"
	}`, docID)))
	d.UseNumber()
	var uData map[string]interface{}
	if err := d.Decode(&uData); err != nil {
		t.Error(err)
	}

	app.updateIndex(uData)
	updatedID, updatedData := app.search("value2", []string{"Full CRUD 2"}, false)

	if len(updatedID) != 1 || len(updatedData) != 1 {
		t.Error("UPDATE: search returned incorrect number of results")
	}
	if updatedID[0] != docID {
		t.Error("UPDATE: IDs dont match")
	}
	if updatedData[0] != `{"index2":"Full CRUD 2"}` {
		t.Error("UPDATE: Data doesnt match")
	}

	app.deleteIndex(docID)
	if _, err := os.Stat(doc); os.IsExist(err) {
		t.Error("DELETE: Document still on disk")
		documentsCleanup([]uint64{docID})
	}
}

func documentsCleanup(docs []uint64) {
	for _, id := range docs {
		doc := fmt.Sprintf("./documents/%d", id)
		os.Remove(doc)
	}
}
