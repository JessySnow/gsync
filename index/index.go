package index

import (
	"encoding/json"
	"fmt"
	"github.com/collinglass/bptree"
	"gync/client/github"
	config2 "gync/config"
	"hash/fnv"
	"time"
)

var indexRoot = bptree.NewTree()
var config config2.Config

func GetRelease(repo config2.Repo, release github.Release) (node *DirNode, err error) {
	key, err := GenerateReleaseKey(repo, release)
	if err != nil {
		return nil, fmt.Errorf("generate release key failed: %v", err)
	}

	record, err := indexRoot.Find(key, false)
	if err != nil {
		return nil, fmt.Errorf("find release key record in bptree failed: %v", err)
	}
	node = new(DirNode)
	err = json.Unmarshal(record.Value, node)
	if err != nil {
		return nil, fmt.Errorf("unmarshal prtree record'value to release failed")
	}

	return
}

// AddRelease add new release to bptree
func AddRelease(repo config2.Repo, release github.Release) (newNode *DirNode, err error) {
	key, err := GenerateReleaseKey(repo, release)
	if err != nil {
		return nil, fmt.Errorf("generate release key failed: %v", err)
	}

	find, err := indexRoot.Find(key, false)
	if find != nil {
		return nil, fmt.Errorf("release already in the index")
	}

	newNode = new(DirNode)
	name, _ := GenerateReleaseDirName(repo, release)
	newNode.DirName = name
	bytes, err := json.Marshal(newNode)
	if err != nil {
		return nil, fmt.Errorf("marshal release node failed: %v", err)
	}
	err = indexRoot.Insert(key, bytes)
	if err != nil {
		return nil, fmt.Errorf("insert release node to bptree failed: %v", err)
	}

	return
}

// GenerateReleaseDirName of a release
func GenerateReleaseDirName(repo config2.Repo, release github.Release) (name string, err error) {
	if len(repo.Owner) == 0 || len(repo.Name) == 0 || len(release.Name) == 0 {
		return "", fmt.Errorf("repo's owner or repo's name or release's name is empty")
	}
	name = fmt.Sprintf("%s/%s/%s", repo.Owner, repo.Name, release.Name)
	return
}

// GenerateReleaseKey generate a key of repo's release
func GenerateReleaseKey(repo config2.Repo, release github.Release) (key int, err error) {
	if len(release.Time) == 0 {
		return -1, fmt.Errorf("release's time is empty")
	}
	h := fnv.New32()
	name, err := GenerateReleaseDirName(repo, release)
	_, _ = h.Write([]byte(name))
	key = int(h.Sum32())

	t, err := time.Parse(time.RFC3339, release.Time)
	if err != nil {
		return -1, fmt.Errorf("parse release time failed: %v", err)
	}

	key += int(t.Unix())
	return
}
