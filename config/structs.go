package config

// Config config of gsync
type Config struct {
	SyncInterval int    `json:"sync_interval"`
	SyncGap      int    `json:"sync_gap"`
	Repos        []Repo `json:"repos"`
}

type Repo struct {
	Owner    string `json:"owner"`
	RepoName string `json:"repoName"`
}
