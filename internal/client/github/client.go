package github

import (
	"encoding/json"
	"fmt"
	"gync/internal/llog"
	"io"
	"net/http"
)

const (
	releaseUrlPrefix = "https://api.github.com/repos"
	releaseUrlSuffix = "releases"
)

// Release represent a release
type Release struct {
	Time   string  `json:"published_at"`
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

// Asset represent a asset in release
type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

// ListRelease fetch Release list
func ListRelease(owner string, repo string) (*[]Release, error) {
	// 0. concat url
	url := fmt.Sprintf("%s/%s/%s/%s", releaseUrlPrefix, owner, repo, releaseUrlSuffix)
	llog.Infof("list release: owner: %s, repo: %s", owner, repo)

	// 1. http request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch release failed: %v", err)
	}

	// 2. unmarshal if request's code is StatusOk
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed with status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read github release response failed: %v", err)
	}
	releases := new([]Release)
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, fmt.Errorf("unmarshal github release response failed: %v", err)
	}

	return releases, nil
}

// DownloadRelease download a release asset return http.Response struct
func DownloadRelease(downloadUrl string) (*http.Response, error) {
	resp, err := http.Get(downloadUrl)
	llog.Infof("download release assrt, download url: %s", downloadUrl)

	if err != nil {
		return nil, fmt.Errorf("download release failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download release failed with status: %s", resp.Status)
	}

	return resp, nil
}
