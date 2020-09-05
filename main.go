package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("Starting GoSearch")
	log.Println("Creating indexmap...")
	app := initApp("test app")
	fmt.Println("")
	input, _ := parseArbJSON(`{"example": "hey oi", "ho": "yo hey hey hey hey"}`)
	input2, _ := parseArbJSON(`{"example": "hey no", "ho": "yo hey hey hey hey"}`)
	app.addIndex(input)
	app.addIndex(input2)
	fmt.Println(app.indexes)
	fmt.Println("### SEARCHING...")
	start := time.Now()
	search := app.search("oi no", make([]string, 0))
	end := time.Now()
	fmt.Println("### SEARCH RESULT:", search)
	fmt.Println("### SEARCH TIME:", end.Sub(start))
	fmt.Println("### TESTING WITH LARGE DATA SET...")
	app = initApp("Movie Data")
	jsonFile, err := os.Open("./moviedata.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("### Successfully Opened users.json...")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	start = time.Now()
	for num, item := range result {
		id := app.addIndex(item)
		if num == 1701 {
			fmt.Println("doc", id)
		}
	}
	end = time.Now()
	fmt.Println("### Added", len(result), "records in", end.Sub(start))

	start = time.Now()
	search = app.search("Deux", make([]string, 0))
	end = time.Now()
	fmt.Println(search)
	fmt.Println("### All fields search done in", end.Sub(start))
	fields := make([]string, 1)
	fields[0] = "year" // notice how it auto casts the int into a string
	start = time.Now()
	search = app.search("19938", fields)
	end = time.Now()
	fmt.Println(search)
	fmt.Println("### Specific fields search done in", end.Sub(start))
	burnfields := make([]string, 3)
	burnfields[0] = "year"
	burnfields[1] = "info"
	burnfields[2] = "title"
	start = time.Now()
	search = app.search("Deux", burnfields)
	end = time.Now()
	fmt.Println(search)
	fmt.Println("### Many fields search done in", end.Sub(start))
}
