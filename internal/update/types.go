package update

type ReleaseResponse struct {
	TagName string  `json:"tag_name"`
	HTMLURL string  `json:"html_url"`
	Message string  `json:"message"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type UpdateInfo struct {
	LatestVersion   string
	ReleasePageURL  string
	WindowsAssetURL string
	HasUpdate       bool
}

type ErrUpdateCheckFailed struct {
	Reason string
}

func (e ErrUpdateCheckFailed) Error() string {
	return "更新检查失败: " + e.Reason
}

type ErrUpdateAvailable struct {
	LatestVersion  string
	ReleasePageURL string
}

func (e ErrUpdateAvailable) Error() string {
	return "发现新版本: " + e.LatestVersion
}
