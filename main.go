package main

import (
	"fmt"
	"log"
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
}
