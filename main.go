package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Starting GoSearch")
	log.Println("Creating indexmap...")
	app := initApp("test app")
	fmt.Println("")
	input, _ := parseArbJSON(`{"example": "hey", "ho": "yo hey hey hey hey"}`)
	app.addIndex(input)
	app.addIndex(input)
	fmt.Println(app.indexes)
}
