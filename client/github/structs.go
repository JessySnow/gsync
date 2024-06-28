package github

// Release GitHub Release 结构体
type Release struct {
	Time   string  `json:"published_at"`
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}
