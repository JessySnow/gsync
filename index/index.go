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

func GetDir(dirName string) (err error, find bool, node *DirNode) {
	h := fnv.New64()
	hash, err := h.Write([]byte(dirName))
	if err != nil {
		return fmt.Errorf("hash dirName filed: %v", err), false, nil
	}

	record, err := indexRoot.Find(hash, false)
	if err != nil {
		return
	}
	if nil == record {
		return nil, false, nil
	}

	node = new(DirNode)
	err = json.Unmarshal(record.Value, node)
	if err != nil {
		return fmt.Errorf("unmarshal bptree record value filed: %v", err), false, nil
	}

	return nil, true, node
}

//
//func UpdateDir(dirName string, node *DirNode) (err error) {
//	h := fnv.New64()
//
//	hash, err := h.Write([]byte(dirName))
//	if err != nil {
//		return fmt.Errorf("hash dirName filed: %v", err)
//	}
//
//	record, err := indexRoot.Find(hash, false)
//	if err != nil {
//		return
//	}
//	if record == nil {
//		return fmt.Errorf("dir node doesn't exists")
//	}
//}

//func MakeDir(dirName string) (err error, node *DirNode) {
//}

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
