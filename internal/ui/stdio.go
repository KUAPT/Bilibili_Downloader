package ui

import (
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var _ UI = (*StdIOUI)(nil)

type StdIOUI struct {
	reader         *bufio.Reader
	ignoreHDR      bool
	ignoreHDRAsked bool
}

func NewStdIOUI() *StdIOUI {
	return &StdIOUI{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (u *StdIOUI) IsInteractive() bool {
	return true
}

func (u *StdIOUI) PromptBVID() (string, error) {
	bvPattern := regexp.MustCompile(`^BV[1-9A-HJ-NP-Za-km-z]{10}$`)
	for {
		fmt.Printf("请输入需要下载视频的BV号：")
		var bvid string
		if _, err := fmt.Scanln(&bvid); err != nil {
			u.ClearScreen()
			fmt.Println("输入读取错误，请重试！")
			log.Println("读取输入错误：", err)
			continue
		}
		if bvPattern.MatchString(bvid) {
			return bvid, nil
		}
		u.ClearScreen()
		fmt.Println("BV号格式错误，请检查格式后重试！")
	}
}

func (u *StdIOUI) ConfirmVideo(info *data_struct.VideoInfoResponse) bool {
	u.ClearScreen()
	printVideoInfo(info)
	return true
}

func (u *StdIOUI) PromptMultiPart() bool {
	fmt.Printf("检测到视频含有分P，是否使用多分P选择/连续下载功能？(Y/n):")
	return yesOrNo(u.reader)
}

func (u *StdIOUI) SelectParts(info *data_struct.VideoInfoResponse, kind int64) ([]DownloadTarget, error) {
	if kind == 1 {
		printUgcSeasonParts(info)
	} else {
		printPagesParts(info)
	}

	fmt.Println("Tip：若需下载所有分P，可直接输入序号0")
	fmt.Printf("请输入用逗号分隔的整数序号（支持中英逗号）：")

	input, _ := u.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "，", ",")

	return parsePartSelection(info, kind, input)
}

func (u *StdIOUI) SelectResolution(defaultMode int64, title string, downloadInfo *data_struct.DownloadInfoResponse) (int, int, string) {
	effectiveQuality := extractEffectiveQualities(downloadInfo)

	var chosenIdx int
	if defaultMode != 1 {
		chosenIdx = u.promptResolutionChoice(title, effectiveQuality, downloadInfo)
	} else {
		fmt.Println("当前下载的视频为：", title)
		chosenIdx = 0
	}

	chosenIdx = u.handleHDRChoice(defaultMode, chosenIdx, effectiveQuality)

	videoIndex := findVideoIndex(downloadInfo, effectiveQuality[chosenIdx])
	videoCode := effectiveQuality[chosenIdx]
	desc := ResolutionNames[videoCode]

	return videoIndex, videoCode, desc
}

func (u *StdIOUI) promptResolutionChoice(title string, qualities []int, info *data_struct.DownloadInfoResponse) int {
	qualityNames := buildQualityNameMap(info)

	for {
		fmt.Println("\n当前下载的视频为：", title)
		fmt.Println("\n请选择想要下载的分辨率：(ps:此处仅显示当前登录账号有权获取的所有分辨率选项)")
		for i, q := range qualities {
			fmt.Printf("%d %s\n", i+1, qualityNames[q])
		}
		fmt.Printf("请输入分辨率前的序号(单个数字)：")

		var choice int
		if _, err := fmt.Scanln(&choice); err != nil {
			u.ClearScreen()
			log.Println("读取输入发生错误")
			fmt.Println("读取输入发生错误,请检查输入格式后重试，若问题依旧，请携带日志log文件向开发者反馈！")
			continue
		}
		if choice < 1 || choice > len(qualities) {
			u.ClearScreen()
			fmt.Println("输入错误，请检查输入后重试！")
			continue
		}
		return choice - 1
	}
}

func (u *StdIOUI) handleHDRChoice(defaultMode int64, chosenIdx int, qualities []int) int {
	if defaultMode != 1 {
		return chosenIdx
	}

	chosenQuality := qualities[chosenIdx]
	if ResolutionNames[chosenQuality] != "HDR" {
		return chosenIdx
	}

	if u.ignoreHDRAsked {
		if u.ignoreHDR {
			return findNonHDRIndex(qualities, chosenIdx)
		}
		return chosenIdx
	}

	fmt.Printf("\n--------【请注意：HDR画质在不受支持的设备上播放将会有显著偏色现象】--------\n")
	fmt.Printf("是否忽略HDR画质?(Y/n):")

	u.ignoreHDRAsked = true
	if yesOrNo(u.reader) {
		u.ignoreHDR = true
		return findNonHDRIndex(qualities, chosenIdx)
	}
	u.ignoreHDR = false
	return chosenIdx
}

func (u *StdIOUI) PromptDefaultResolution() int64 {
	fmt.Printf("当前为多分P下载模式，是否默认下载可获取的最高分辨率?(Y/n)")
	if yesOrNo(u.reader) {
		return 1
	}
	return -1
}

func (u *StdIOUI) PromptHDRSkip() bool {
	fmt.Printf("\n--------【请注意：HDR画质在不受支持的设备上播放将会有显著偏色现象】--------\n")
	fmt.Printf("是否忽略HDR画质?(Y/n):")
	return yesOrNo(u.reader)
}

func (u *StdIOUI) PromptContinue() bool {
	fmt.Printf("是否继续下载其他视频？(Y/n):")
	return yesOrNo(u.reader)
}

func (u *StdIOUI) PromptUpdateDownload(currentVersion, latestVersion string) bool {
	fmt.Printf("发现新版本，当前版本: %s，最新版本: %s\n", currentVersion, latestVersion)
	fmt.Print("是否下载更新? (Y/n): ")
	return yesOrNo(u.reader)
}

func (u *StdIOUI) ShowMessage(msg string) {
	fmt.Println(msg)
}

func (u *StdIOUI) ShowError(msg string) {
	fmt.Println(msg)
}

func (u *StdIOUI) WaitForExit() {
	fmt.Printf("程序执行完毕，请按Enter键退出...")
	_, _ = u.reader.ReadString('\n')
}

func (u *StdIOUI) ClearScreen() {
	clearScreen()
}

func (u *StdIOUI) ShowNoPermission() {
	log.Println("当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址")
	fmt.Printf("\n当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址\n")
}

func (u *StdIOUI) ResetHDRState() {
	u.ignoreHDR = false
	u.ignoreHDRAsked = false
}

func yesOrNo(reader *bufio.Reader) bool {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "" || input == "y" || input == "yes"
}

func parsePartSelection(info *data_struct.VideoInfoResponse, kind int64, input string) ([]DownloadTarget, error) {
	parts := strings.Split(input, ",")
	selectedIndices := make(map[int]bool)
	selectAll := false

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if num, err := strconv.Atoi(part); err == nil {
			if num == 0 {
				selectAll = true
				break
			}
			selectedIndices[num-1] = true
		}
	}

	var targets []DownloadTarget

	if kind == 1 {
		episodes := info.Data.UgcSeason.Sections[0].Episodes
		for i, ep := range episodes {
			if selectAll || selectedIndices[i] {
				targets = append(targets, DownloadTarget{
					Title: ep.Title,
					BVID:  ep.Bvid,
					CID:   ep.Page.Cid,
				})
			}
		}
	} else {
		pages := info.Data.Pages
		for i, page := range pages {
			if selectAll || selectedIndices[i] {
				targets = append(targets, DownloadTarget{
					Title: page.Part,
					BVID:  info.Data.Bvid,
					CID:   page.Cid,
				})
			}
		}
	}

	return targets, nil
}

func extractEffectiveQualities(info *data_struct.DownloadInfoResponse) []int {
	seen := make(map[int]bool)
	var result []int
	for _, v := range info.Data.Dash.Video {
		if !seen[v.ID] {
			seen[v.ID] = true
			result = append(result, v.ID)
		}
	}
	return result
}

func buildQualityNameMap(info *data_struct.DownloadInfoResponse) map[int]string {
	m := make(map[int]string)
	for i, desc := range info.Data.AcceptDescription {
		if i < len(info.Data.AcceptQuality) {
			m[info.Data.AcceptQuality[i]] = desc
		}
	}
	return m
}

func findVideoIndex(info *data_struct.DownloadInfoResponse, quality int) int {
	for i, v := range info.Data.Dash.Video {
		if v.ID == quality {
			return i
		}
	}
	return 0
}

func findNonHDRIndex(qualities []int, current int) int {
	for i, q := range qualities {
		if ResolutionNames[q] != "HDR" {
			return i
		}
	}
	fmt.Println("未找到其他非HDR画质选项，继续使用HDR画质")
	return current
}
