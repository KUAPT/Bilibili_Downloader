package app

import (
	"Bilibili_Downloader/internal"
	"Bilibili_Downloader/internal/downloader"
	"Bilibili_Downloader/internal/paths"
	"Bilibili_Downloader/internal/sso"
	"Bilibili_Downloader/internal/ui"
	"Bilibili_Downloader/internal/update"
	"Bilibili_Downloader/internal/video_processing"
	"Bilibili_Downloader/pkg/httpclient"
	"Bilibili_Downloader/pkg/toolkit"
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	BVID       string
	NoUpdate   bool
	NonInter   bool
	OutputPath string
}

type App struct {
	ui     ui.UI
	config Config
}

func New(u ui.UI, cfg Config) *App {
	return &App{
		ui:     u,
		config: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	if !a.config.NoUpdate && a.ui.IsInteractive() {
		a.handleUpdate(ctx)
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		bvid, err := a.getBVID()
		if err != nil {
			return err
		}
		if bvid == "" {
			break
		}

		if err := a.processVideo(ctx, bvid); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			log.Printf("处理视频失败: %v\n", err)
		}

		if ctx.Err() != nil {
			break
		}

		if !a.ui.IsInteractive() || !a.ui.PromptContinue() {
			break
		}

		if stdioUI, ok := a.ui.(*ui.StdIOUI); ok {
			stdioUI.ResetHDRState()
		}
		a.ui.ClearScreen()
	}

	return nil
}

func (a *App) getBVID() (string, error) {
	if a.config.BVID != "" {
		return a.config.BVID, nil
	}

	if !a.ui.IsInteractive() {
		return "", fmt.Errorf("non-interactive mode requires --bvid")
	}

	return a.ui.PromptBVID()
}

func (a *App) processVideo(ctx context.Context, bvid string) error {
	videoInfoURL := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvid)

	videoInfo, err := fetchVideoInfo(videoInfoURL)
	if err != nil {
		return err
	}

	a.ui.ConfirmVideo(videoInfo)

	targets, err := a.selectDownloadTargets(videoInfo)
	if err != nil {
		return err
	}

	defaultResMode := int64(-2)

	for _, target := range targets {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := a.downloadAndTranscode(ctx, target, &defaultResMode, len(targets) > 1); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			log.Printf("下载失败: %s - %v\n", target.Title, err)
			a.ui.ShowError("请求下载失败，请检查网络连接或前往log文件查看详情.")
			a.ui.ShowMessage("跳过当前视频，3s后继续下载下一个视频...")
			time.Sleep(3 * time.Second)
		}
	}

	return nil
}

func (a *App) selectDownloadTargets(info *data_struct.VideoInfoResponse) ([]ui.DownloadTarget, error) {
	hasMultiPart := info.Data.UgcSeason.Sections != nil || len(info.Data.Pages) > 1

	if !hasMultiPart {
		return []ui.DownloadTarget{{
			Title: info.Data.Title,
			BVID:  info.Data.Bvid,
			CID:   info.Data.Cid,
		}}, nil
	}

	if !a.ui.IsInteractive() || a.ui.PromptMultiPart() {
		var kind int64 = 2
		if len(info.Data.Pages) <= 1 {
			kind = 1
		}
		return a.ui.SelectParts(info, kind)
	}

	return []ui.DownloadTarget{{
		Title: info.Data.Title,
		BVID:  info.Data.Bvid,
		CID:   info.Data.Cid,
	}}, nil
}

func (a *App) downloadAndTranscode(ctx context.Context, target ui.DownloadTarget, defaultResMode *int64, multiOutput bool) error {
	downloadURL := fmt.Sprintf(
		"https://api.bilibili.com/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=4048",
		target.BVID, target.CID,
	)

	data, err := internal.CatchData(downloadURL)
	if err != nil {
		log.Printf("获取下载信息数据发生错误: %v\n\n", err)
		return fmt.Errorf("获取下载信息失败: %w", err)
	}

	resp, err := internal.ProcessResponse(data, 1)
	if err != nil {
		log.Printf("处理下载信息发生错误: %v\n\n", err)
		return fmt.Errorf("处理下载信息失败: %w", err)
	}

	downloadInfo, ok := resp.(*data_struct.DownloadInfoResponse)
	if !ok {
		return fmt.Errorf("下载信息类型断言失败")
	}

	if len(downloadInfo.Data.AcceptDescription) == 0 || downloadInfo.Data.AcceptDescription[0] == "试看" {
		a.ui.ShowNoPermission()
		return nil
	}

	if *defaultResMode == -2 && a.ui.IsInteractive() {
		*defaultResMode = a.ui.PromptDefaultResolution()
	} else if *defaultResMode == -2 {
		*defaultResMode = 1
	}

	videoIdx, _, resDesc := a.ui.SelectResolution(*defaultResMode, target.Title, downloadInfo)

	videoURL := downloadInfo.Data.Dash.Video[videoIdx].BackupURL[0]
	audioURL := downloadInfo.Data.Dash.Audio[0].BackupURL[0]

	progressWriter := io.Discard
	if a.ui.IsInteractive() {
		progressWriter = os.Stdout
	}

	dlResult, err := downloader.Download(ctx, httpclient.GetClient(), paths.Get(), audioURL, videoURL, downloader.Options{ProgressWriter: progressWriter})
	if err != nil {
		return err
	}

	defaultName := toolkit.CheckAndCleanFileName(fmt.Sprintf("%s(%s)", target.Title, resDesc)) + ".mp4"
	outputRel, err := resolveOutputRel(a.config.OutputPath, defaultName, multiOutput)
	if err != nil {
		return err
	}

	a.ui.ShowMessage("开始视频转码：\n")
	if _, err := video_processing.Transcode(ctx, paths.Get(), dlResult, outputRel); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)

	return nil
}

func fetchVideoInfo(url string) (*data_struct.VideoInfoResponse, error) {
	data, err := internal.CatchData(url)
	if err != nil {
		log.Printf("获取视频信息数据错误: %v\n\n", err)
		return nil, fmt.Errorf("视频信息数据获取异常: %w", err)
	}

	resp, err := internal.ProcessResponse(data, 0)
	if err != nil {
		log.Printf("处理视频详情发生错误: %v\n\n", err)
		return nil, fmt.Errorf("视频详情数据处理发生错误: %w", err)
	}

	videoInfo, ok := resp.(*data_struct.VideoInfoResponse)
	if !ok {
		return nil, fmt.Errorf("视频详情数据类型断言失败")
	}

	return videoInfo, nil
}

func InitClient() bool {
	return httpclient.Init()
}

func HandleLogin() error {
	return sso.HandleQRCodeLogin()
}

func resolveOutputRel(flagValue string, defaultName string, multiOutput bool) (string, error) {
	if flagValue == "" {
		return defaultName, nil
	}

	clean := filepath.Clean(flagValue)
	if clean == "." || filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid --output path")
	}

	lower := strings.ToLower(clean)
	if multiOutput {
		if strings.HasSuffix(lower, ".mp4") {
			clean = strings.TrimSuffix(clean, filepath.Ext(clean))
			clean = filepath.Clean(clean)
			if clean == "." || strings.HasPrefix(clean, "..") {
				return "", fmt.Errorf("invalid --output path")
			}
		}
		return filepath.Join(clean, defaultName), nil
	}

	if !strings.HasSuffix(lower, ".mp4") {
		clean += ".mp4"
	}

	return clean, nil
}

func (a *App) handleUpdate(ctx context.Context) {
	update.HandleUpdateFlag()

	if err := update.EnsureConfigVersion(); err != nil {
		log.Println("配置版本检查失败:", err)
		return
	}

	apiURL, err := update.GetUpdateAPIURL()
	if err != nil {
		log.Println("获取更新API失败:", err)
		return
	}

	a.ui.ShowMessage("检查更新...")
	info, err := update.CheckForUpdate(ctx, apiURL, update.GetCurrentVersion())
	if err != nil {
		var checkErr update.ErrUpdateCheckFailed
		if errors.As(err, &checkErr) {
			log.Println("更新检查失败:", checkErr.Reason)
			a.ui.ShowMessage("检查更新失败，请检查网络环境\n")
		}
		return
	}

	if !info.HasUpdate {
		a.ui.ShowMessage("当前已是最新版本！\n\n")
		log.Println("完成检查更新过程")
		return
	}

	if !a.ui.PromptUpdateDownload(update.GetCurrentVersion(), info.LatestVersion) {
		a.ui.ShowMessage("跳过更新\n\n")
		return
	}

	needRestart, newProgramName, err := update.ApplyUpdate(ctx, info)
	if err != nil {
		var availableErr update.ErrUpdateAvailable
		if errors.As(err, &availableErr) {
			a.ui.ShowMessage(fmt.Sprintf("发现新版本 %s，请前往以下地址手动下载更新：\n%s\n\n",
				availableErr.LatestVersion, availableErr.ReleasePageURL))
			return
		}
		log.Println("应用更新失败:", err)
		a.ui.ShowError("更新失败: " + err.Error())
		return
	}

	if needRestart {
		a.ui.ShowMessage("更新成功，正在重启...")
		if err := update.LaunchNewVersion(newProgramName, os.Args[0]); err != nil {
			log.Println("启动新程序失败:", err)
			a.ui.ShowError("启动更新程序失败")
			return
		}
		os.Exit(0)
	}
}
