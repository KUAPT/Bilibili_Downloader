package main

import (
	"aBlibliTop/detail"
	"bufio"
	"fmt"
	"log"
	"os"
)

//-ldflags="-s -w"

func main() {
	defer detail.RemoveCacheDir()

	fmt.Println("请输入你需要下载视频的BV号：")
	var VideoId string
	if _, err := fmt.Scanln(&VideoId); err != nil {
		fmt.Println("BV号输入错误！")
	}
	VideoInfoUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", VideoId)

	cookie, flag := detail.LoadConfig()
	if flag != 0 {
		fmt.Println("配置加载错误，请检查config目录下json文件是否存在或格式是否正确！")
		return
	}
	fmt.Println("配置文件加载成功！")
	if len(cookie) >= 10 {
		fmt.Printf("当前使用的Cookie为（前十个字符）：%v\n", cookie[:10])
	} else {
		fmt.Println("cookie已成功加载，但可能存在问题，请检查填写是否正常（确认无误可忽略此警告）")
	}

	// 获取数据
	data, err := detail.CatchData(VideoInfoUrl, cookie)
	if err != nil {
		log.Fatalf("获取视频信息数据错误: %v\n", err)
		return
	}

	// 打印原始数据以进行调试
	//fmt.Printf("原始数据: %s\n", data)

	// 处理数据
	Response, err := detail.ProcessResponse(data, 0)
	if err != nil {
		log.Fatalf("处理视频信息响应错误: %v\n", err)
		return
	}
	videoInfoResponse := Response.(*detail.VideoInfoResponse)
	//fmt.Printf("视频信息解组后的数据：%v\n", videoInfoResponse)

	DownloadURL := fmt.Sprintf("https://api.bilibili.com/x/player/wbi/playurl?bvid=%s&cid=%d", videoInfoResponse.Data.Bvid, videoInfoResponse.Data.Cid)

	data, err = detail.CatchData(DownloadURL, cookie)
	if err != nil {
		log.Fatalf("获取下载信息数据错误: %v\n", err)
		return
	}
	//fmt.Printf("原始数据: %s\n\n\n\n\n", data)

	newResponse, err := detail.ProcessResponse(data, 1)
	if err != nil {
		log.Fatalf("下载信息处理响应错误: %v\n", err)
		return
	}

	downloadInfoResponse := newResponse.(*detail.DownloadInfoResponse)
	//fmt.Printf("下载信息解组后的数据：%v\n", downloadInfoResponse)

	if err := detail.DownloadFile(downloadInfoResponse.Data.Durl[0].URL, "", cookie); err != nil {
		fmt.Printf("请求下载失败：%s\n", err)
	}

	detail.Transcoding(videoInfoResponse.Data.Title)

	fmt.Println("程序执行完毕，请按Enter键退出...")
	// 创建一个新的读取器
	reader := bufio.NewReader(os.Stdin)
	// 读取一个字符
	_, _ = reader.ReadString('\n')
}
