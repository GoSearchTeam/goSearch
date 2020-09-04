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
}
