package config

import (
	"encoding/json"
	"fmt"
)

// Config config of gsync
type Config struct {
	RootDir      string `json:"root_dir"`
	SyncInterval int    `json:"sync_interval"`
	SyncGap      int    `json:"sync_gap"`
	Repos        []Repo `json:"repos"`
}

// Repo repo to sync
type Repo struct {
	Owner string `json:"owner"`
	Name  string `json:"repoName"`
}

// Parse parse config
func Parse(bytes []byte) (config *Config, err error) {
	config = new(Config)
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("parse config filed: %v", err)
	}
	return
}
