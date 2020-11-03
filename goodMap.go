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

type WalkFn func(id uint64, value uint64) bool

func NewList() *GoodList {
	newList := GoodList{nil, nil}
	return &newList
}

// AddItem inserts the node into the GoodList with lazy sorting
func (goodList *GoodList) AddItem(item uint64) (newValue uint64, err error) {
	fmt.Println("Adding", item)
	// Find insert or update
	newNode := GoodNode{
		id:     item,
		prev:   nil,
		next:   nil,
		value:  0,
		mapRef: goodList,
	}
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

		higherNode := tempNode     // Keep track of last spot where value is higher if updating to move up
		for tempNode.next != nil { // Leave on end of list or next is item
			if tempNode.next.id == item {
				break
			}
			// TODO: Break when I get to value 1 so we can insert there and not walk to end
			tempNode = tempNode.next
			if tempNode.value < higherNode.value { // The node we will need to place after
				higherNode = tempNode.prev
			}
		}

		if tempNode.next == nil { // Append to end of list
			tempNode.next = &newNode
			newNode.prev = tempNode
			goodList.last = &newNode
			newNode.value++
			return newNode.value, nil
		} else {
			updateNode := tempNode.next
			// Check if higherNode same as tempNode
			// Manual adjust
			if tempNode.value > updateNode.value && higherNode.value >= tempNode.value {
				higherNode = tempNode
			}
			if higherNode != tempNode {
				tempNode.next = updateNode.next
				// Check if last node
				if updateNode.next != nil {
					updateNode.next.prev = tempNode
				}
				higherNode.next.prev = updateNode
				updateNode.next = higherNode.next
				updateNode.prev = higherNode
				higherNode.next = updateNode
			} else {
				if higherNode.value < updateNode.value+1 {
					higherNode.next = updateNode.next
					if updateNode.next != nil {
						updateNode.next.prev = higherNode
					}
					updateNode.next = higherNode
					if higherNode.prev != nil {
						updateNode.prev = higherNode.prev
					} else { // First node in list
						goodList.first = updateNode
						updateNode.prev = nil
					}
					higherNode.prev = updateNode
				}
			}
			updateNode.value++
			return updateNode.value, nil
		}
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
func (goodList *GoodList) Walk(walkFn WalkFn) {
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
	theString := "-"
	goodList.Walk(func(id uint64, value uint64) bool {
		theString = fmt.Sprintf("%v(%v, %v)-", theString, id, value)
		return true
	})
	return theString
}
