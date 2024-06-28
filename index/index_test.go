package index

import (
	"gync/client/github"
	config2 "gync/config"
	"reflect"
	"testing"
)

var repo1 = config2.Repo{Owner: "A", Name: "B"}
var release1 = github.Release{Name: "R1", Time: "2022-06-05T12:37:28Z"}
var repo2 = config2.Repo{Owner: "A", Name: "B"}
var release2 = github.Release{Name: "R1", Time: "2022-06-05T12:37:29Z"}
var node1 = DirNode{DirName: "A/B/R1"}

func TestGenerateReleaseKey(t *testing.T) {
	type args struct {
		repo    config2.Repo
		release github.Release
	}
	tests := []struct {
		name    string
		args    args
		wantKey int
		wantErr bool
	}{
		{args: args{repo: config2.Repo{Owner: "A", Name: "B"}, release: github.Release{Name: "R1", Time: "2022-06-05T12:37:28Z"}}, wantKey: 5084752807},
		{args: args{repo: config2.Repo{Owner: "A", Name: "B"}, release: github.Release{Name: "R1", Time: "2022-06-05T12:37:29Z"}}, wantKey: 5084752808},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := GenerateReleaseKey(tt.args.repo, tt.args.release)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateReleaseKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKey != tt.wantKey {
				t.Errorf("GenerateReleaseKey() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestGenerateReleaseDirName(t *testing.T) {
	type args struct {
		repo    config2.Repo
		release github.Release
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{args: args{repo: config2.Repo{Owner: "A", Name: "B"}, release: github.Release{Name: "R1", Time: "2022-06-05T12:37:28Z"}}, wantName: "A/B/R1"},
		{args: args{repo: config2.Repo{Owner: "A", Name: "B"}, release: github.Release{Name: "R2", Time: "2022-06-05T12:37:29Z"}}, wantName: "A/B/R2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, err := GenerateReleaseDirName(tt.args.repo, tt.args.release)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateReleaseDirName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("GenerateReleaseDirName() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}

func TestAddRelease(t *testing.T) {
	type args struct {
		repo    config2.Repo
		release github.Release
	}
	tests := []struct {
		name        string
		args        args
		wantNewNode *DirNode
		wantErr     bool
	}{
		{args: args{repo: repo1, release: release1}, wantNewNode: &node1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewNode, err := AddRelease(tt.args.repo, tt.args.release)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewNode, tt.wantNewNode) {
				t.Errorf("AddRelease() gotNewNode = %v, want %v", gotNewNode, tt.wantNewNode)
			}
		})
	}
}

func TestGetRelease(t *testing.T) {
	type args struct {
		repo    config2.Repo
		release github.Release
	}
	tests := []struct {
		name     string
		args     args
		wantNode *DirNode
		wantErr  bool
	}{
		{args: args{repo: repo1, release: release1}, wantNode: &node1},
	}
	_, _ = AddRelease(repo1, release1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNode, err := GetRelease(tt.args.repo, tt.args.release)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("GetRelease() gotNode = %v, want %v", gotNode, tt.wantNode)
			}
		})
	}
}
