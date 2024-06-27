package config

import (
	"reflect"
	"testing"
)

var testConfig = Config{10, 5, []Repo{{"owner1", "repo1"}, {"owner2", "repo2"}}}

func TestParse(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *Config
		wantErr    bool
	}{
		{wantConfig: &testConfig, args: args{bytes: []byte("{\n  \"sync_interval\": 10,\n  \"sync_gap\": 5,\n  \"repos\": [\n    {\n      \"owner\": \"owner1\",\n      \"repoName\": \"repo1\"\n    },\n    {\n      \"owner\": \"owner2\",\n      \"repoName\": \"repo2\"\n    }\n  ]\n}")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepos, err := Parse(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRepos, tt.wantConfig) {
				t.Errorf("Parse() gotRepos = %v, want %v", gotRepos, tt.wantConfig)
			}
		})
	}
}
