package ui

import (
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var ResolutionNames = map[int]string{
	6:   "240P",
	16:  "360P",
	32:  "480P",
	64:  "720P",
	74:  "720P60",
	80:  "1080P",
	112: "1080P+",
	116: "1080P60",
	120: "4K",
	125: "HDR",
	126: "杜比视界",
	127: "8K超高清",
}

var ResolutionPriority = []int{
	6, 16, 32, 64, 74, 80, 112, 116, 120, 126, 127,
}

func SelectBestResolution(available []int, skipHDR bool) int {
	if len(available) == 0 {
		return 0
	}

	avail := make(map[int]bool)
	for _, q := range available {
		avail[q] = true
	}

	for i := len(ResolutionPriority) - 1; i >= 0; i-- {
		q := ResolutionPriority[i]
		if avail[q] {
			return q
		}
	}

	if !skipHDR && avail[125] {
		return 125
	}

	if avail[125] {
		return 125
	}

	return available[0]
}

func printVideoInfo(info *data_struct.VideoInfoResponse) {
	fmt.Println("========== 视频信息 ==========")
	fmt.Println("标题:", info.Data.Title)
	fmt.Println("UP主:", info.Data.Owner.Name)
	fmt.Println("BV号:", info.Data.Bvid)
	fmt.Println("简介:", truncateString(info.Data.Desc, 100))
	fmt.Println("==============================")
}

func printUgcSeasonParts(info *data_struct.VideoInfoResponse) {
	fmt.Println("\n========== 合集分P列表 ==========")
	episodes := info.Data.UgcSeason.Sections[0].Episodes
	for i, ep := range episodes {
		fmt.Printf("%d. %s\n", i+1, ep.Title)
	}
	fmt.Println("================================")
}

func printPagesParts(info *data_struct.VideoInfoResponse) {
	fmt.Println("\n========== 视频分P列表 ==========")
	for i, page := range info.Data.Pages {
		fmt.Printf("%d. %s\n", i+1, page.Part)
	}
	fmt.Println("================================")
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
