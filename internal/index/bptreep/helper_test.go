package bptreep

import (
	"github.com/collinglass/bptree"
	"testing"
)

func TestKeyIterator(t *testing.T) {
	tree := bptree.NewTree()
	_ = tree.Insert(1, nil)
	_ = tree.Insert(2, nil)
	_ = tree.Insert(3, nil)
	_ = tree.Insert(4, nil)
	_ = tree.Insert(5, nil)
	_ = tree.Insert(6, nil)
	_ = tree.Insert(7, nil)

	leaf, _ := findLeaf(tree.Root, 1)
	iterator := LeafNodeKeyIterator{node: leaf, index: 0}
	ints := []int{1, 2, 3, 4, 5, 6, 7}
	for _, i := range ints {
		get := *iterator.getNext()
		if i != get {
			t.Errorf("iterator get failed, except: %d, get: %d", i, get)
			return
		}
	}

	if iterator.hasNext() {
		t.Errorf("iterator hasNext check failed, except: %v, get: %v", false, true)
	}
}

func findLeaf(node *bptree.Node, key int) (leaf *bptree.Node, err error) {
	i := 0
	leaf = node
	if leaf == nil {
		return nil, nil
	}
	for !leaf.IsLeaf {
		i = 0
		for i < leaf.NumKeys {
			if key >= leaf.Keys[i] {
				i += 1
			} else {
				break
			}
		}
		leaf, _ = leaf.Pointers[i].(*bptree.Node)
	}
	return leaf, nil
}
