package main

import (
	"Bilibili_Downloader/internal"
	"Bilibili_Downloader/internal/sso"
	"Bilibili_Downloader/internal/tool"
	"Bilibili_Downloader/internal/video_processing"
	"Bilibili_Downloader/pkg/httpclient"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Println("log文件close失败:", err)
		}
	}()
	// 将日志输出设置到文件
	log.SetOutput(logFile)
	// 设置日志前缀和格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("程序运行，开始日志记录")

	defer func() {
		fmt.Printf("程序执行完毕，请按Enter键退出...")
		// 创建一个新的读取器
		reader := bufio.NewReader(os.Stdin)
		// 读取一个字符
		_, _ = reader.ReadString('\n')
		log.Println("程序执行完毕，正常退出")
	}()

	defer func() {
		if err := tool.RemoveCacheDir(); err != nil {
			tool.ClearScreen()
			log.Println("缓存目录清理失败:", err)
			fmt.Println("缓存目录清理失败，确认需清理时可手动清理或重新运行程序")
		} else {
			fmt.Println("缓存目录清理完毕")
		}
	}()

	if success := httpclient.Init(); !success {
		err := sso.HandleQRCodeLogin()
		if err != nil {
			fmt.Println("处理二维码登录失败:", err)
			return
		}
	}

	//正则对BV号进行基本检查
	BVcheak := regexp.MustCompile(`^BV[1-9A-HJ-NP-Za-km-z]{10}$`)
	var BVid string
	for {
		fmt.Printf("请输入需要下载视频的BV号：")
		if _, err := fmt.Scanln(&BVid); err != nil {
			tool.ClearScreen()
			fmt.Println("输入读取错误，请重试！")
			log.Println("读取输入错误：", err)
			continue
		}
		if BVcheak.MatchString(BVid) {
			break
		} else {
			tool.ClearScreen()
			fmt.Println("BV号格式错误，请检查格式后重试！")
		}
	}
	VideoInfoUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", BVid)

	// 获取数据
	data, err := internal.CatchData(VideoInfoUrl)
	if err != nil {
		log.Println("获取视频信息数据错误: %v\n", err)
		fmt.Println("视频信息数据获取异常，请检查网络连接或前往log文件查看详情.")
		return
	}

	// 处理数据
	Response, err := internal.ProcessResponse(data, 0)
	if err != nil {
		log.Println("处理视频详情发生错误: %v\n", err)
		fmt.Println("视频详情数据处理发生错误，请携带log文件向开发者反馈！")
		return
	}
	videoInfoResponse := Response.(*internal.VideoInfoResponse)

	DownloadURL := fmt.Sprintf("https://api.bilibili.com/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=80&fnver=0&fourk=1&qn=0", videoInfoResponse.Data.Bvid, videoInfoResponse.Data.Cid)

	data, err = internal.CatchData(DownloadURL)
	if err != nil {
		log.Println("获取下载信息数据发生错误: %v\n", err)
		fmt.Println("视频下载信息获取异常，请检查网络连接或前往log文件查看详情.")
		return
	}

	newResponse, err := internal.ProcessResponse(data, 1)
	if err != nil {
		log.Println("处理下载信息发生错误: %v\n", err)
		fmt.Println("视频下载信息处理发生错误，请携带log文件向开发者反馈！")
		return
	}

	downloadInfoResponse := newResponse.(*internal.DownloadInfoResponse)

	if err := internal.DownloadFile(downloadInfoResponse.Data.Dash.Video[0].BackupURL[0], downloadInfoResponse.Data.Dash.Audio[0].BackupURL[0], ""); err != nil {
		log.Printf("请求下载失败：%s\n", err)
		fmt.Println("请求下载失败，请检查网络连接或前往log文件查看详情.")
	}

	fmt.Println("开始视频转码：\n")
	video_processing.Transcoding(videoInfoResponse.Data.Title)
}
