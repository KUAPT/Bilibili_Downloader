package video_processing

import (
	"Bilibili_Downloader/internal/downloader"
	"Bilibili_Downloader/internal/paths"
	"Bilibili_Downloader/pkg/toolkit"
	"context"
	"fmt"
	"log"
)

func Transcoding(videoName string, definition string) {
	p := paths.Get()
	outputRel := toolkit.CheckAndCleanFileName(fmt.Sprintf("%s(%s)", videoName, definition)) + ".mp4"
	input := downloader.Result{AudioPath: p.AudioCachePath(), VideoPath: p.VideoCachePath()}
	_, err := Transcode(context.Background(), p, input, outputRel)
	if err != nil {
		log.Println("视频文件转码失败:", err)
		return
	}
	log.Println("视频文件转码成功")
}
