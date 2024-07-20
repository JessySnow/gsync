package bptindex

import (
	"testing"
)

func Test_leafNodeKeyIterator_getNext(t *testing.T) {
	tree := newEnhanceBptree()
	_ = tree.Insert(1, nil)
	_ = tree.Insert(2, nil)
	_ = tree.Insert(3, nil)
	_ = tree.Insert(4, nil)
	_ = tree.Insert(5, nil)
	_ = tree.Insert(6, nil)
	_ = tree.Insert(7, nil)
	leaf, _ := tree.findLeaf(1)
	iterator, _ := leaf.newIterator()
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	for _, i := range nums {
		get, _ := iterator.getNext()
		if i != get {
			t.Errorf("iterator get failed, except: %d, get: %d", i, get)
			return
		}
	}

	if iterator.hasNext() {
		t.Errorf("iterator hasNext check failed, except: %v, get: %v", false, true)
	}
}

func Test_leafNodeKeyIterator_hasNext(t *testing.T) {
	tree := newEnhanceBptree()
	_ = tree.Insert(1, nil)
	_ = tree.Insert(2, nil)
	_ = tree.Insert(3, nil)
	_ = tree.Insert(4, nil)
	_ = tree.Insert(5, nil)
	_ = tree.Insert(6, nil)
	_ = tree.Insert(7, nil)
	_ = tree.Insert(8, nil)
	leaf, _ := tree.findLeaf(1)
	iterator, _ := leaf.newIterator()
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	for _, i := range nums {
		get, _ := iterator.getNext()
		if i != get {
			t.Errorf("iterator get failed, except: %d, get: %d", i, get)
			return
		}
	}

	if !iterator.hasNext() {
		t.Errorf("iterator hasNext check failed, except: %v, get: %v", false, true)
	}
	_, _ = iterator.getNext()
	if iterator.hasNext() {
		t.Errorf("iterator hasNext check failed, except: %v, get: %v", false, true)
	}
}
