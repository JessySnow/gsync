package index

import (
	"encoding/json"
	"fmt"
	"github.com/collinglass/bptree"
	"gync/internal/client/github"
	config2 "gync/internal/config"
	"hash/fnv"
	"os"
	"path"
	"time"
)

var indexRoot *bptree.Tree
var context *config2.Config

// init bp tree from root dir
func init() {
	indexRoot = bptree.NewTree()
}

// InitIndex build tree and init config
func InitIndex(config *config2.Config) (err error) {
	context = config
	err = bulkLoad()
	return
}

func GetRelease(repo *config2.Repo, release *github.Release) (node *DirNode, err error) {
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
func AddRelease(repo *config2.Repo, release *github.Release) (newNode *DirNode, err error) {
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

func UpdateRelease(repo *config2.Repo, release *github.Release, newNode *DirNode) (err error) {
	key, err := GenerateReleaseKey(repo, release)
	if err != nil {
		return fmt.Errorf("generate release key failed: %v", err)
	}

	origin, err := indexRoot.Find(key, false)
	if err != nil {
		return fmt.Errorf("find release in bptree failed: %v", err)
	}

	newValue, err := json.Marshal(*newNode)
	if err != nil {
		return fmt.Errorf("marshal dirNode failed: %v", err)
	}
	origin.Value = newValue
	return
}

// GenerateReleaseDirName of a release
func GenerateReleaseDirName(repo *config2.Repo, release *github.Release) (name string, err error) {
	if len(repo.Owner) == 0 || len(repo.Name) == 0 || len(release.Name) == 0 {
		return "", fmt.Errorf("repo's owner or repo's name or release's name is empty")
	}
	name = fmt.Sprintf("%s/%s/%s", repo.Owner, repo.Name, release.Name)
	return
}

// GenerateReleaseKey generate a key of repo's release
func GenerateReleaseKey(repo *config2.Repo, release *github.Release) (key int, err error) {
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

// bulkLoad bulk load dir to bptree
func bulkLoad() (err error) {
	if context == nil || len(context.RootDir) == 0 {
		return fmt.Errorf("config check failed: %v", context)
	}

	entries, err := os.ReadDir(context.RootDir)
	if err != nil {
		return fmt.Errorf("open root dirctionary filed: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		owner := entry.Name()
		subEntries, err := os.ReadDir(path.Join(context.RootDir, owner))
		if err != nil {
			return err
		}

		for _, subEntry := range subEntries {
			if !subEntry.IsDir() {
				continue
			}

			repo := &config2.Repo{Name: subEntry.Name(), Owner: owner}
			subSubEntries, err := os.ReadDir(path.Join(context.RootDir, owner, subEntry.Name()))
			if err != nil {
				return err
			}

			for _, subSubEntry := range subSubEntries {
				if !subSubEntry.IsDir() {
					continue
				}

				meta, err := os.ReadFile(path.Join(context.RootDir, owner, subEntry.Name(), subSubEntry.Name(), ".meta"))
				if err != nil {
					return fmt.Errorf("failed to read release metadata: %v", err)
				}

				release := new(github.Release)
				err = json.Unmarshal(meta, release)
				if err != nil {
					return fmt.Errorf("unable to parse release's metadata")
				}

				_, err = AddRelease(repo, release)
				if err != nil {
					return fmt.Errorf("unable to bulkLoad index: %v", err)
				}
			}
		}
	}
	return
}

func findLeaf(key int) (leaf *bptree.Node, err error) {
	i := 0
	leaf = indexRoot.Root
	if leaf == nil {
		return nil, fmt.Errorf("index root is nil")
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
