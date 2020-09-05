package main

import (
	"fmt"
	"goSearch/tree"
	"strings"
)

type indexTree struct {
	name string
	tree *tree.Tree
	len  int64
}

func (tree indexTree) addTreeIndex(k interface{}, v interface{}) {
	value := fmt.Sprintf("%v", v)
	fmt.Println("### Adding tree index...")
	// Break up string into individual words (if string)
	fmt.Println("### Splitting:", value)
	split := strings.Fields(value)
	fmt.Printf("### Split value: %q\n", split)
	for _, item := range split {
		key := stringToNum(item)
		fmt.Println("### Key: ", key)
		tree.tree.Insert(key, []byte(item))
		tree.len++
	}
}

type Forest struct {
	trees []indexTree
	name  string
}

// Jungle - A Forest is kind of like an app, or collection of indexes
type Jungle []Forest

func initTree(forest *Forest, name string) *indexTree {
	newTree := indexTree{name, tree.NewTree(), 0}
	forest.trees = append(forest.trees, newTree)
	return &newTree
}

// Create index from parsed json
func addIndex(parsed map[string]interface{}, forest *Forest) {
	fmt.Println("### Adding index...")
	// TODO: Store document
	// TODO: check if tree exists with name of every json key, if not create tree
	// For all of the k,v pairs in the json
	for k, v := range parsed {
		var treePointer *indexTree = nil
		for i := 0; i < len(forest.trees); i++ {
			if k == forest.trees[i].name { // Found existing tree
				fmt.Println("### Found existing tree!")
				treePointer = &forest.trees[i]
				break
			}
		}

		if treePointer == nil { // Create tree
			fmt.Println("### Creating new tree...")
			treePointer = initTree(forest, k)
		}

		// Add index to forest
		treePointer.addTreeIndex(k, v)
	}

}

func saveTree(tree indexTree) (err error) {

	return nil
}
