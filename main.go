package main

import (
	"Bilibili_Downloader/detail"
	"Bilibili_Downloader/httpclient"
	"Bilibili_Downloader/sso"
	"Bilibili_Downloader/tool"
	"Bilibili_Downloader/video_processing"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	defer func() {
		if err := tool.RemoveCacheDir(); err != nil {
			fmt.Println("缓存目录清理失败，确认需清理时可手动清理或重新运行程序.")
			log.Println("缓存目录清理失败:", err)
		}
	}()

	logFile, err := os.OpenFile("detail.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer logFile.Close()
	// 将日志输出设置到文件
	log.SetOutput(logFile)
	// 设置日志前缀和格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("程序运行，开始日志记录")

	if success := httpclient.Init(); !success {
		err := sso.HandleQRCodeLogin()
		if err != nil {
			fmt.Println("处理二维码登录失败:", err)
			os.Exit(1)
		}
	}

	//正则对BV号进行基本检查
	BVcheak := `^BV[1-9A-HJ-NP-Za-km-z]{10}$`
	cheak, err := regexp.Compile(BVcheak)
	var BV_id string
	for {
		fmt.Printf("请输入你需要下载视频的BV号：")
		if _, err := fmt.Scanln(&BV_id); err != nil {
			tool.ClearScreen()
			fmt.Println("输入读取错误，请重试！")
			continue
		}

		if cheak.MatchString(BV_id) {
			break
		} else {
			tool.ClearScreen()
			fmt.Println("BV号输入错误，请重试！")
		}
	}
	VideoInfoUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", BV_id)

	// 获取数据
	data, err := detail.CatchData(VideoInfoUrl)
	if err != nil {
		log.Fatalf("获取视频信息数据错误: %v\n", err)
		return
	}

	// 处理数据
	Response, err := detail.ProcessResponse(data, 0)
	if err != nil {
		log.Fatalf("视频信息处理发生错误: %v\n", err)
		return
	}
	videoInfoResponse := Response.(*detail.VideoInfoResponse)

	DownloadURL := fmt.Sprintf("https://api.bilibili.com/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=80&fnver=0&fourk=1&qn=0", videoInfoResponse.Data.Bvid, videoInfoResponse.Data.Cid)

	data, err = detail.CatchData(DownloadURL)
	if err != nil {
		log.Fatalf("获取下载信息数据错误: %v\n", err)
		return
	}

	newResponse, err := detail.ProcessResponse(data, 1)
	if err != nil {
		log.Fatalf("下载信息处理发生错误: %v\n", err)
		return
	}

	downloadInfoResponse := newResponse.(*detail.DownloadInfoResponse)

	if err := detail.DownloadFile(downloadInfoResponse.Data.Dash.Video[0].BackupURL[0], downloadInfoResponse.Data.Dash.Audio[0].BackupURL[0], ""); err != nil {
		log.Printf("请求下载失败：%s\n", err)
	}

	fmt.Println("开始视频转码：\n")
	video_processing.Transcoding(videoInfoResponse.Data.Title)

	fmt.Printf("程序执行完毕，请按Enter键退出...")
	// 创建一个新的读取器
	reader := bufio.NewReader(os.Stdin)
	// 读取一个字符
	_, _ = reader.ReadString('\n')
	log.Println("程序执行完毕，正常退出")
}
