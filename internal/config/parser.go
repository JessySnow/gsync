package config

import (
	"encoding/json"
	"fmt"
)

// Parse parse config
func Parse(bytes []byte) (config *Config, err error) {
	config = new(Config)
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("parse config filed: %v", err)
	}
	return
}
