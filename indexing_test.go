package main

import (
	"fmt"
	"os"
	"testing"
)

func TestAddIndex(t *testing.T) {
	app := initApp("TestAddIndex")
	data, _ := parseArbJSON(`{
		"name": "Test",
		"food": "Tacos",
		"time": "High noon"
	}`)
	id := app.addIndex(data)

	doc := fmt.Sprintf("./documents/%d", id)
	if _, err := os.Stat(doc); os.IsNotExist(err) {
		t.Errorf("%s not created", doc)
	}
	ids := []uint64{id}
	documentsCleanup(ids)
}

func TestSearch(t *testing.T) {
	app := initApp("TestAddIndex")
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

	ids := []uint64{tacoID, burgerID}
	documentsCleanup(ids)

	if tacoSearchID[0] != tacoID || burgerSearchID[0] != burgerID {
		t.Error("ID's didnt match")
	}
	if tacoSearchData[0] != `{"name":"tacos"}` || burgerSearchData[0] != `{"name":"burger"}` {
		t.Error("Documents didnt match")
	}

}

func documentsCleanup(docs []uint64) {
	for _, id := range docs {
		doc := fmt.Sprintf("./documents/%d", id)
		os.Remove(doc)
	}
}
