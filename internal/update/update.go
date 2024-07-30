package update

import (
	"Bilibili_Downloader/internal/toolkit"
	"Bilibili_Downloader/pkg/config"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/cheggaaa/pb/v3"
	"log"
	"net/http"
	"os"
	"strings"
)

// checkForUpdate 检查更新
func checkForUpdate(VersionUpdateApi string, currentVersion string) (string, string, error) {
	resp, err := http.Get(VersionUpdateApi)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("close body failed")
		}
	}()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to get latest release: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	latestVersionStr := result["tag_name"].(string)
	downloadURL := result["assets"].([]interface{})[0].(map[string]interface{})["browser_download_url"].(string)

	// 使用 semver 库比较版本
	currentSemver, err := semver.NewVersion(strings.TrimPrefix(currentVersion, "v"))
	if err != nil {
		return "", "", err
	}
	latestSemver, err := semver.NewVersion(strings.TrimPrefix(latestVersionStr, "v"))
	if err != nil {
		return "", "", err
	}

	if latestSemver.GreaterThan(currentSemver) {
		return downloadURL, latestVersionStr, nil // 有更新
	}
	return "", "", nil // 没有更新
}

// 下载更新
func downloadUpdate(downloadURL string, downloadVersion string) (string, error) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("close body failed")
		}
	}()

	tempFile, err := os.Create("update_temp")
	if err != nil {
		return "", err
	}

	bar := pb.StartNew(int(resp.ContentLength))
	bar.Set(pb.SIBytesPrefix, true)
	err = toolkit.DownloadAndTrackProgress(resp.Body, tempFile, bar)
	if err != nil {
		return "", err
	}
	bar.Finish()

	if err := tempFile.Close(); err != nil {
		log.Println("close tempFile failed")
		return "", err
	}

	newProgramName := "BiliBili_Downloader_" + downloadVersion + ".exe"
	err = os.Rename("update_temp", newProgramName)
	if err != nil {
		return "", err
	}

	return newProgramName, nil
}

// 处理更新替换
func update(oldVersionPath string) {
	if err := os.Remove(oldVersionPath); err != nil {
		log.Println("删除旧版本失败:", err)
		return
	}
	log.Println("成功删除旧版本，完成程序更新")
	fmt.Println("完成程序更新")
}

// CheckAndUpdate 返回值 -1：err；0; 正常无更新或跳过更新; 1：正常且更新
func CheckAndUpdate() (int, string) {
	if len(os.Args) > 1 && os.Args[1] == "--update" {
		// 被老版本启动用于更新
		if len(os.Args) < 3 {
			log.Println("更新参数缺失")
			fmt.Println("无法正确删除旧版本，可手动删除旧版本")
		}
		update(os.Args[2])
	}

	fmt.Println("读取配置...")
	currentConfig, err := config.ReadConfig()
	if err != nil {
		log.Println("读取配置时发生错误:", err)
		fmt.Println("读取配置时发生错误:", err)
		return -1, ""
	}

	fmt.Println("检查更新...")
	downloadURL, latestVersion, err := checkForUpdate(currentConfig.VersionUpdateApi, currentConfig.CurrentVersion)
	if err != nil {
		log.Println("检查更新时发生错误:", err)
		fmt.Println("检查更新时发生错误:", err)
		return -1, ""
	}

	if downloadURL != "" {
		fmt.Printf("发现新版本，当前版本: %s，最新版本: %s\n", currentConfig.CurrentVersion, latestVersion)
		fmt.Print("是否下载更新? (y/n): ")
		if toolkit.YesOrNo() {
			if newProgramName, err := downloadUpdate(downloadURL, latestVersion); err != nil {
				log.Println("更新失败:", err)
				fmt.Println("更新失败:", err)
			} else {
				fmt.Println("更新成功，请重新启动程序！")
				return 1, newProgramName
			}
		} else {
			fmt.Println("跳过更新。")
			return 0, ""
		}
	} else {
		fmt.Println("当前已是最新版本！")
		return 0, ""
	}
	return -1, ""
}