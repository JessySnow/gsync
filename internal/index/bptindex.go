package index

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sort"
	"time"
)

type BptreeReleaseDirIndexer struct {
	indexTree *enhanceBptree
}

func (i *BptreeReleaseDirIndexer) Locate(rr *RepoRelease) (node *DirNode, err error) {
	key, err := generateKey(rr)
	if err != nil {
		return nil, fmt.Errorf("generate release key failed: %v", err)
	}

	record, err := i.indexTree.Find(key, false)
	if err != nil {
		return nil, fmt.Errorf("find key record in bptree failed: %v", err)
	}
	node = new(DirNode)
	err = json.Unmarshal(record.Value, node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal prtree record'value to release failed")
	}

	return
}

func (i *BptreeReleaseDirIndexer) Add(rr *RepoRelease) (node *DirNode, err error) {
	key, err := generateKey(rr)
	if err != nil {
		return nil, fmt.Errorf("generate release key failed: %v", err)
	}

	find, err := i.indexTree.Find(key, false)
	if find != nil {
		return nil, fmt.Errorf("key already record in bptree : %v", err)
	}

	newNode := new(DirNode)
	dirName, err := rr.GenerateReleaseDirName()
	if err != nil {
		return
	}

	newNode.DirName = dirName
	bytes, err := json.Marshal(newNode)
	if err != nil {
		return nil, fmt.Errorf("marshal dir node failed: %v", err)
	}
	err = i.indexTree.Insert(key, bytes)
	if err != nil {
		return nil, fmt.Errorf("insert dir node to bptree failed: %v", err)
	}

	return newNode, nil
}

func (i *BptreeReleaseDirIndexer) Update(rr *RepoRelease, node *DirNode) (err error) {
	key, err := generateKey(rr)
	if err != nil {
		return fmt.Errorf("generate release key failed: %v", err)
	}

	origin, err := i.indexTree.Find(key, false)
	if err != nil {
		return fmt.Errorf("find release in bptree failed: %v", err)
	}

	newValue, err := json.Marshal(*node)
	if err != nil {
		return fmt.Errorf("marshal dirNode failed: %v", err)
	}
	origin.Value = newValue
	return
}

func (i *BptreeReleaseDirIndexer) GetAbsent(rrs []RepoRelease) (absent []RepoRelease, err error) {
	if 0 == len(rrs) {
		return rrs, nil
	}

	// generate key for each RepoRelease and save all key to a slice
	keys := make([]int, 0)
	key2RepoRelease := make(map[int]RepoRelease)
	for _, rr := range rrs {
		key, err := generateKey(&rr)
		if err != nil {
			return nil, fmt.Errorf("generate key failed: %v", err)
		}
		key2RepoRelease[key] = rr
		keys = append(keys, key)
	}

	// sort keys
	sort.Ints(keys)

	// find key in bptree index
	leaf, err := i.indexTree.findLeaf(keys[0])
	if err != nil {
		return rrs, nil
	}
	iterator, err := leaf.newIterator()
	if err != nil {
		return nil, fmt.Errorf("failed to get leaf iterator")
	}

	// traverse through keys
	absent = *new([]RepoRelease)
	for _, key := range keys {
		find := false
		for iterator.hasNext() {
			current, _ := iterator.getNext()
			if current == key {
				find = true
				break
			}
			if current > key {
				break
			}
		}

		if !find {
			iterator.reset()
			absent = append(absent, key2RepoRelease[key])
		} else {
			iterator.snapshot()
		}
	}

	return absent, nil
}

func New() (indexer *BptreeReleaseDirIndexer, err error) {
	indexer = &BptreeReleaseDirIndexer{indexTree: newEnhanceBptree()}
	return indexer, nil
}

func generateKey(rr *RepoRelease) (key int, err error) {
	if len(rr.RepoName) == 0 || len(rr.ReleaseTime) == 0 {
		return -1, fmt.Errorf("enpty repo name or releaseTime: %v", rr)
	}
	h := fnv.New32()
	_, _ = h.Write([]byte(rr.RepoName))
	key = int(h.Sum32())

	t, err := time.Parse(time.RFC3339, rr.ReleaseTime)
	if err != nil {
		return -1, fmt.Errorf("parse release time failed: %v", err)
	}

	key += int(t.Unix())
	return
}
