package index

import (
	"errors"
	"fmt"
	"github.com/collinglass/bptree"
)

// bptreeOrder bptree.node's order see bptree.tree.defaultOrder
const bptreeOrder = 4

var (
	bptreeRootNil = errors.New("bptree's root node is nil")
	noMoreElement = errors.New("no more element to return")
)

type enhanceBptree struct {
	*bptree.Tree
}

type enhanceBptreeLeafNode struct {
	*bptree.Node
}

type leafNodeKeyIterator struct {
	index         int
	node          *enhanceBptreeLeafNode
	snapshotIndex int
	snapshotNode  *enhanceBptreeLeafNode
}

func newEnhanceBptree() *enhanceBptree {
	innerTree := bptree.NewTree()
	return &enhanceBptree{innerTree}
}

// findLeaf get a leaf node which may contains key and wrap it in enhanceBptreeLeafNode
func (e *enhanceBptree) findLeaf(key int) (leaf *enhanceBptreeLeafNode, err error) {
	i := 0
	if e.Root == nil {
		return nil, bptreeRootNil
	}
	leaf = &enhanceBptreeLeafNode{e.Root}

	for !leaf.IsLeaf {
		i = 0
		for i < leaf.NumKeys {
			if key >= leaf.Keys[i] {
				i += 1
			} else {
				break
			}
		}
		leaf = &enhanceBptreeLeafNode{leaf.Pointers[i].(*bptree.Node)}
	}
	return leaf, nil
}

func (e *enhanceBptreeLeafNode) newIterator() (iterator *leafNodeKeyIterator, err error) {
	if !e.IsLeaf {
		return nil, fmt.Errorf("not a leaf node")
	}

	return &leafNodeKeyIterator{0, e, 0, e}, nil
}

func (it *leafNodeKeyIterator) hasNext() bool {
	if nil == it.node {
		return false
	}

	if it.index >= it.node.NumKeys && nil == it.node.Pointers[bptreeOrder-1] {
		return false
	} else {
		return true
	}
}

func (it *leafNodeKeyIterator) getNext() (key int, err error) {
	if it.hasNext() {
		if it.index < it.node.NumKeys {
			key = it.node.Keys[it.index]
			it.index += 1
			return key, nil
		} else {
			it.index = 0
			it.node = &enhanceBptreeLeafNode{it.node.Pointers[bptreeOrder-1].(*bptree.Node)}
			return it.getNext()
		}
	} else {
		return key, noMoreElement
	}
}

func (it *leafNodeKeyIterator) snapshot() {
	it.snapshotIndex = it.index
	it.snapshotNode = it.node
}

func (it *leafNodeKeyIterator) reset() {
	it.node = it.snapshotNode
	it.index = it.snapshotIndex
}
