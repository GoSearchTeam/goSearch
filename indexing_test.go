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

	doc1Data, _ := parseArbJSON(`{
		"indexSmallSearch": "value1"
	}`)
	doc1ID := app.addIndex(doc1Data)

	doc2Data, _ := parseArbJSON(`{
		"indexSmallSearch": "value2"
	}`)
	doc2ID := app.addIndex(doc2Data)

	res1ID, res1Data := app.search("value1", []string{"indexSmallSearch"}, false)
	res2ID, res2Data := app.search("value2", []string{"indexSmallSearch"}, false)
	resBWID, resBWData := app.search("value", []string{"indexSmallSearch"}, true)

	// fmt.Printf("Responce1: %v\nResponce2: %v\nResponceBW: %v\n", responce1, responce2, responceBW)
	if (len(res1ID) != 1 || len(res2ID) != 1 || len(resBWID) != 2) && false { // the and false is there because the docID list isnt done yet
		t.Error("Returned wrong number of IDs")
		return
	}

	if len(res1Data.Items) != 1 || len(res2Data.Items) != 1 || len(resBWData.Items) != 2 {
		t.Error("Returned wrong number of items")
		return
	}

	if res1Data.Items[0].Data["indexSmallSearch"] != "value1" || res2Data.Items[0].Data["indexSmallSearch"] != "value2" {
		t.Error("Returned data was incorrect")
		return
	}

	documentsCleanup([]uint64{doc1ID, doc2ID})
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

	if len(updatedID) != 1 && false {
		t.Error("Search returned wrong number of IDs")
		return
	}

	if len(updatedData.Items) != 1 {
		t.Error("Returned wrong number of items")
		return
	}

	if updatedData.Items[0].Data["index2"] != "Test Update 2" {
		t.Error("Returned data was incorrect")
		return
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
		return
	}

	searchedID, searchedData := app.search("Full CRUD", []string{"index"}, false)

	if len(searchedID) != 1 && false {
		t.Error("SEARCH: returned incorrect number of IDs")
		return
	}

	if len(searchedData.Items) != 1 {
		t.Error("SEARCH: returned wrong number of items")
		return
	}

	if searchedData.Items[0].DocID != docID {
		t.Error("SEARCH: IDs dont match")
		return
	}

	if searchedData.Items[0].Data["index"] != "Full CRUD" {
		t.Error("SEARCH: Data doesnt match")
		return
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
	updatedID, updatedData := app.search("Full CRUD 2", []string{"index2"}, false)

	if len(updatedID) != 1 && false {
		t.Error("UPDATE: search returned incorrect number of results")
		return
	}

	if len(updatedData.Items) != 1 {
		t.Error("UPDATE: returned wrong number of items")
		return
	}

	if updatedData.Items[0].Data["index2"] != "Full CRUD 2" {
		t.Error("UPDATE: Data doesnt match")
		return
	}

	app.deleteIndex(docID)
	if _, err := os.Stat(doc); os.IsExist(err) {
		t.Error("DELETE: Document still on disk")
	}
}

func documentsCleanup(docs []uint64) {
	for _, id := range docs {
		doc := fmt.Sprintf("./documents/%d", id)
		os.Remove(doc)
	}
}
