package update

import (
	"Bilibili_Downloader/pkg/httpclient"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func CheckForUpdate(ctx context.Context, apiURL, currentVersion string) (UpdateInfo, error) {
	resp, err := httpclient.Get(ctx, apiURL)
	if err != nil {
		return UpdateInfo{}, ErrUpdateCheckFailed{Reason: fmt.Sprintf("网络请求失败: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateInfo{}, ErrUpdateCheckFailed{Reason: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	var release ReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return UpdateInfo{}, ErrUpdateCheckFailed{Reason: fmt.Sprintf("JSON 解析失败: %v", err)}
	}

	if release.Message != "" {
		return UpdateInfo{}, ErrUpdateCheckFailed{Reason: release.Message}
	}

	if release.TagName == "" {
		return UpdateInfo{}, ErrUpdateCheckFailed{Reason: "响应缺少 tag_name"}
	}

	currentSemver, err := semver.NewVersion(strings.TrimPrefix(currentVersion, "v"))
	if err != nil {
		return UpdateInfo{}, ErrUpdateCheckFailed{Reason: fmt.Sprintf("当前版本解析失败: %v", err)}
	}

	latestSemver, err := semver.NewVersion(strings.TrimPrefix(release.TagName, "v"))
	if err != nil {
		return UpdateInfo{}, ErrUpdateCheckFailed{Reason: fmt.Sprintf("最新版本解析失败: %v", err)}
	}

	info := UpdateInfo{
		LatestVersion:  release.TagName,
		ReleasePageURL: release.HTMLURL,
		HasUpdate:      latestSemver.GreaterThan(currentSemver),
	}

	info.WindowsAssetURL = selectWindowsAsset(release.Assets)

	return info, nil
}

func selectWindowsAsset(assets []Asset) string {
	var fallbackExe string

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)

		if name == "bilibili_downloader.exe" {
			return asset.BrowserDownloadURL
		}

		if strings.HasSuffix(name, ".exe") {
			if strings.Contains(name, "bilibili_downloader") {
				return asset.BrowserDownloadURL
			}
			if fallbackExe == "" {
				fallbackExe = asset.BrowserDownloadURL
			}
		}
	}

	return fallbackExe
}
