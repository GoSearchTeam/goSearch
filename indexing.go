package main

import (
	"github.com/collinglass/bptree"
)

type indexTree struct {
	name string
	tree *bptree.Tree
	len  int64
}

func (tree indexTree) addTreeIndex(key interface{}, value interface{}) {
	// TODO: Break up string into individual words if string
}

type Forest struct {
	trees []indexTree
	name  string
}

// Jungle - A Forest is kind of like an app, or collection of indexes
type Jungle []Forest

func initTree(forest *Forest, name string) (tree *indexTree) {
	newTree := indexTree{name, bptree.NewTree(), 0}
	forest.trees = append(forest.trees, newTree)
	return &newTree
}

// Create index from parsed json
func addIndex(parsed map[string]interface{}, forest *Forest) {
	// TODO: Store document
	// TODO: check if tree exists with name of every json key, if not create tree
	// For all of the k,v pairs in the json
	for k, v := range parsed {
		var treePointer *indexTree = nil
		for i := 0; i < len(forest.trees); i++ {
			if k == forest.trees[i].name { // Found existing tree
				treePointer = &forest.trees[i]
				break
			}
		}

		if treePointer == nil { // Create tree
			treePointer = initTree(forest, k)
		}

		// Add index to forest
		treePointer.addTreeIndex(k, v)
	}

}

func saveTree(tree indexTree) (err error) {

	return nil
}
