package bptreep

import (
	"github.com/collinglass/bptree"
)

type LeafNodeKeyIterator struct {
	index int
	node  *bptree.Node
}

func (it *LeafNodeKeyIterator) hasNext() bool {
	if nil == it.node {
		return false
	}

	if it.index >= it.node.NumKeys && nil == it.node.Pointers[3] {
		return false
	} else {
		return true
	}
}

func (it *LeafNodeKeyIterator) getNext() *int {
	if it.hasNext() {
		if it.index < it.node.NumKeys {
			key := it.node.Keys[it.index]
			it.index += 1
			return &key
		} else {
			it.index = 0
			it.node = it.node.Pointers[3].(*bptree.Node)
			return it.getNext()
		}
	} else {
		return nil
	}
}
