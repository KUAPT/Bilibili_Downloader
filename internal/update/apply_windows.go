//go:build windows

package update

import (
	"Bilibili_Downloader/internal/paths"
	"Bilibili_Downloader/pkg/httpclient"
	"Bilibili_Downloader/pkg/toolkit"
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func ApplyUpdate(ctx context.Context, info UpdateInfo) (needRestart bool, newProgramName string, err error) {
	if info.WindowsAssetURL == "" {
		return false, "", ErrUpdateCheckFailed{Reason: "无可用的 Windows 下载资源"}
	}

	resp, err := httpclient.Get(ctx, info.WindowsAssetURL)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	p := paths.Get()
	tempPath := p.UpdateTempPath()
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return false, "", err
	}

	bar := pb.StartNew(int(resp.ContentLength))
	bar.Set(pb.SIBytesPrefix, true)
	err = toolkit.DownloadAndTrackProgress(resp.Body, tempFile, bar)
	if err != nil {
		tempFile.Close()
		return false, "", err
	}
	bar.Finish()

	if err := tempFile.Close(); err != nil {
		log.Println("close tempFile failed:", err)
		return false, "", err
	}

	newProgramName = "Bilibili_Downloader_" + info.LatestVersion + ".exe"
	newProgramPath := filepath.Join(p.Root, newProgramName)
	if err := os.Rename(tempPath, newProgramPath); err != nil {
		return false, "", err
	}

	return true, newProgramName, nil
}

func LaunchNewVersion(newProgramName string, oldProgramPath string) error {
	p := paths.Get()
	newProgram := filepath.Join(p.Root, newProgramName)
	cmd := exec.Command("cmd.exe", "/K", "start", "", newProgram, "--update", oldProgramPath)
	return cmd.Start()
}

func CleanupOldVersion(oldVersionPath string) error {
	if err := os.Remove(oldVersionPath); err != nil {
		log.Println("删除旧版本失败:", err)
		return err
	}
	log.Println("成功删除旧版本，完成程序更新")
	return nil
}
