package index

import "fmt"

// DirNode an abstract representation of a path in a file system.
type DirNode struct {
	DirName      string   `json:"dir_name"`       // dirname from config's RootDir
	Locked       bool     `json:"locked"`         // indicates a dir been locked
	SubFileNames []string `json:"sub_file_names"` // sub files name
}

// RepoRelease an abstract of repo and one of it's release
type RepoRelease struct {
	RepoOwner   string
	RepoName    string
	ReleaseName string
	ReleaseTime string
}

func (rr *RepoRelease) GenerateReleaseDirName() (dirName string, err error) {
	if len(rr.RepoOwner) == 0 || len(rr.RepoName) == 0 || len(rr.ReleaseName) == 0 {
		return "", fmt.Errorf("repo's owner or repo's name or release's name is empty")
	}

	dirName = fmt.Sprintf("%s/%s/%s", rr.RepoOwner, rr.RepoName, rr.ReleaseName)
	return
}

type ReleaseDirIndex interface {
	Locate(rr *RepoRelease) (node *DirNode, err error)
	Add(rr *RepoRelease) (node *DirNode, err error)
	Update(rr *RepoRelease, node *DirNode) (err error)
	GetAbsent(rrs []RepoRelease) (absent []RepoRelease, err error)
}
