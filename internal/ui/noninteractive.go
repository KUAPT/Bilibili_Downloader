package ui

import (
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"errors"
	"fmt"
	"log"
)

var _ UI = (*NonInteractiveUI)(nil)

var ErrNoInteraction = errors.New("non-interactive mode does not support user input")

type NonInteractiveUI struct {
	skipHDR bool
}

func NewNonInteractiveUI() *NonInteractiveUI {
	return &NonInteractiveUI{
		skipHDR: true,
	}
}

func (u *NonInteractiveUI) IsInteractive() bool {
	return false
}

func (u *NonInteractiveUI) PromptBVID() (string, error) {
	return "", ErrNoInteraction
}

func (u *NonInteractiveUI) ConfirmVideo(info *data_struct.VideoInfoResponse) bool {
	printVideoInfo(info)
	return true
}

func (u *NonInteractiveUI) PromptMultiPart() bool {
	return true
}

func (u *NonInteractiveUI) SelectParts(info *data_struct.VideoInfoResponse, kind int64) ([]DownloadTarget, error) {
	var targets []DownloadTarget

	if kind == 1 {
		episodes := info.Data.UgcSeason.Sections[0].Episodes
		for _, ep := range episodes {
			targets = append(targets, DownloadTarget{
				Title: ep.Title,
				BVID:  ep.Bvid,
				CID:   ep.Page.Cid,
			})
		}
	} else {
		for _, page := range info.Data.Pages {
			targets = append(targets, DownloadTarget{
				Title: page.Part,
				BVID:  info.Data.Bvid,
				CID:   page.Cid,
			})
		}
	}

	return targets, nil
}

func (u *NonInteractiveUI) SelectResolution(defaultMode int64, title string, downloadInfo *data_struct.DownloadInfoResponse) (int, int, string) {
	fmt.Println("当前下载的视频为：", title)

	effectiveQuality := extractEffectiveQualities(downloadInfo)
	bestQuality := SelectBestResolution(effectiveQuality, u.skipHDR)

	videoIndex := findVideoIndex(downloadInfo, bestQuality)
	desc := ResolutionNames[bestQuality]

	log.Printf("自动选择清晰度: %s (%d)\n", desc, bestQuality)
	fmt.Printf("自动选择清晰度: %s\n", desc)

	return videoIndex, bestQuality, desc
}

func (u *NonInteractiveUI) PromptDefaultResolution() int64 {
	return 1
}

func (u *NonInteractiveUI) PromptHDRSkip() bool {
	return true
}

func (u *NonInteractiveUI) PromptContinue() bool {
	return false
}

func (u *NonInteractiveUI) PromptUpdateDownload(currentVersion, latestVersion string) bool {
	return false
}

func (u *NonInteractiveUI) ShowMessage(msg string) {
	fmt.Println(msg)
}

func (u *NonInteractiveUI) ShowError(msg string) {
	fmt.Println(msg)
}

func (u *NonInteractiveUI) WaitForExit() {
}

func (u *NonInteractiveUI) ClearScreen() {
}

func (u *NonInteractiveUI) ShowNoPermission() {
	log.Println("当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址")
	fmt.Printf("\n当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址\n")
}
