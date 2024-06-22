package github

import (
	"fmt"
	"net/http"
	"testing"
)

func TestListRelease(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    []Release
		wantErr bool
	}{
		{
			args: args{repo: "Xray-core", owner: "XTLS"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListRelease(tt.args.owner, tt.args.repo)
			if err != nil {
				panic(err)
			}
			fmt.Println(got)
		})
	}
}

func TestDownloadRelease(t *testing.T) {
	type args struct {
		downloadUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			args: args{downloadUrl: "https://github.com/XTLS/Xray-core/releases/download/v1.8.16/Xray-android-arm64-v8a.zip"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadRelease(tt.args.downloadUrl)
			if err != nil {
				panic(err)
			}
			fmt.Println(got.Body)
		})
	}
}
