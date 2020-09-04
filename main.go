package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Starting GoSearch")
	log.Println("Creating Forest...")
	forest := Forest{}
	initTree(&forest, "testTree")
	fmt.Println("### forest:")
	fmt.Println(forest)
	fmt.Println("### Going over Forest:")
	for i := 0; i < len(forest.trees); i++ {
		fmt.Println(forest.trees[i])
	}
	parsed, _ := parseArbJSON(`{"thing1": "first thing", "thing2": "second thing"}`)
	addIndex(parsed, &forest)
	fmt.Println("### Going over Forest:")
	for i := 0; i < len(forest.trees); i++ {
		fmt.Println(forest.trees[i])
		val, _ := forest.trees[i].tree.Find(538, false)
		if val != nil {
			fmt.Println("VALUES: ", string(val.Value))
		}
	}
}
