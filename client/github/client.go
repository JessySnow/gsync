package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	releaseUrlPrefix = "https://api.github.com/repos"
	releaseUrlSuffix = "releases"
)

// ListRelease 拉取 Release 列表
func ListRelease(owner string, repo string) (*[]Release, error) {
	// 0. 组装 url
	url := fmt.Sprintf("%s/%s/%s/%s", releaseUrlPrefix, owner, repo, releaseUrlSuffix)

	// 1. 向 GitHub 发起拉取 Release 的请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch release failed: %v", err)
	}

	// 2. 判断 HTTP 响应
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed with status: %s", resp.Status)
	}

	// 3.反序列化
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read github release response failed: %v", err)
	}
	releases := make([]Release, 0)
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, fmt.Errorf("unmarshal github release response failed: %v", err)
	}

	return &releases, nil
}

// DownloadRelease 下载具体的 Release 内容
func DownloadRelease(downloadUrl string) (*http.Response, error) {
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return nil, fmt.Errorf("download release failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download release failed with status: %s", resp.Status)
	}

	return resp, nil
}
