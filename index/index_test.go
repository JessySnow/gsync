package index

import (
	"gync/client/github"
	config2 "gync/config"
	"testing"
)

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
