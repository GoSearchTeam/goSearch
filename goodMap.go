package main

import ()
import "fmt"

type GoodList struct {
	first *GoodNode
	last  *GoodNode
}

type GoodNode struct {
	id     uint64
	value  uint64
	prev   *GoodNode
	next   *GoodNode
	mapRef *GoodList
}

func NewList() *GoodList {
	newList := GoodList{nil, nil}
	return &newList
}

// AddItem inserts the node into the GoodList with lazy sorting
func (goodList *GoodList) AddItem(item uint64) (newValue uint64, err error) {
	fmt.Println("Adding", item)
	newNode := GoodNode{
		id:     item,
		prev:   nil,
		next:   nil,
		value:  0,
		mapRef: goodList,
	}
	// Find insert or update
	if goodList.first == nil { // First item
		newNode.value++
		goodList.first = &newNode
		goodList.last = &newNode
		return newNode.value, nil
	} else {
		tempNode := goodList.first
		if tempNode.id == item { // Check if first node
			tempNode.value++
			return tempNode.value, nil
		}

		higherNode := tempNode                                 // Keep track of last spot where value is higher if updating to move up
		for tempNode.next != nil || tempNode.next.id != item { // Leave on end of list or next is item
			tempNode = tempNode.next
			if tempNode.value < higherNode.value { // The node we will need to place after
				higherNode = tempNode.prev
			}
		}
		updateNode := tempNode.next

		if tempNode.next == nil { // Append to end of list
			tempNode.next = &newNode
			newNode.prev = tempNode
		} else {
			tempNode.next = updateNode.next
			// Check if last node
			if updateNode.next != nil {
				updateNode.next.prev = tempNode
			}
			// Check if higherNode same as tempNode

		}
		newNode.value++
		return newNode.value, nil

		for tempNode.next != nil {
			if tempNode.value < higherNode.value { // Update higher counter node pointer
				higherNode = tempNode
			}
			if tempNode.next.id == item { // Update
				updateNode := tempNode.next
				updateNode.value++

				// Special Cases --- IM GETTING MYSELF SO CONFUSED IM HARD CODING SOME CASES

				// First nodes
				if tempNode.id == goodList.first.id {
					fmt.Println("first nodes")
					goodList.first = updateNode
					if updateNode.next == nil { // only 2 nodes in list
						fmt.Println("only 2 in list")
						goodList.last = tempNode
					}
					tempNode.next = updateNode.next
					tempNode.prev = updateNode
					updateNode.prev = nil
					updateNode.next = tempNode
					return updateNode.value, nil
				}

				tempNode.next = updateNode.next

				if higherNode == tempNode { // Switch spots
					updateNode.prev = tempNode.prev
					tempNode.prev = updateNode
					// Last nodes
					if updateNode.next == nil {
						fmt.Println("last node")
						goodList.last = tempNode
					} else {
						updateNode.next.prev = tempNode
					}
					updateNode.next = tempNode
					updateNode.prev.next = updateNode
					return updateNode.value, nil
				}

				// Move behind higher count node
				// Remove from spot
				fmt.Println("step out")
				updateNode.next.prev = tempNode
				// Insert in new spot
				updateNode.next = higherNode.next
				if higherNode.id != tempNode.id { // Not same node
					higherNode.next.prev = updateNode
				} else if higherNode.id == goodList.first.id { // Update first node
					goodList.first = updateNode
				}
				higherNode.next = updateNode
				updateNode.prev = higherNode

				return tempNode.next.value, nil
			}
			tempNode = tempNode.next
		}
		// New node
		tempNode.next = &newNode
		newNode.prev = tempNode
		newNode.value = 1
		return 1, nil
	}
}

func GoodListFromArray(arr []uint64) (*GoodList, error) {
	newList := GoodList{nil, nil}
	for _, item := range arr {
		_, err := newList.AddItem(item)
		if err != nil {
			return nil, err
		}
	}
	return &newList, nil
}

// WalkFn walks the list. Return whether to keep walking.
func (goodList *GoodList) WalkFn(walkFn func(id uint64, value uint64) bool) {
	tempNode := goodList.first
	if goodList.first == nil {
		return
	}
	for tempNode != nil {
		keepGoing := walkFn(tempNode.id, tempNode.value)
		if keepGoing == false { // If false returned
			return
		}
		tempNode = tempNode.next
	}
}

// ToString returns a string of "(id, value), (id, value)..."
func (goodList *GoodList) ToString() string {
	theString := ""
	goodList.WalkFn(func(id uint64, value uint64) bool {
		theString = fmt.Sprintf("%v(%v, %v), ", theString, id, value)
		return true
	})
	return theString
}
