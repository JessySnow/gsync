package index

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

type ReleaseDirIndex interface {
	Locate(rr *RepoRelease) (node *DirNode, err error)
	Add(rr *RepoRelease) (node *DirNode, err error)
	Update(rr *RepoRelease, node *DirNode) (err error)
}
