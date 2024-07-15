package main

import (
	"Bilibili_Downloader/internal"
	"Bilibili_Downloader/internal/sso"
	"Bilibili_Downloader/internal/toolkit"
	"Bilibili_Downloader/internal/video_processing"
	"Bilibili_Downloader/pkg/httpclient"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// 创建一个新的读取器
	reader := bufio.NewReader(os.Stdin)

	//初始化日志文件
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
		// 读取一个字符
		_, _ = reader.ReadString('\n')
		log.Println("程序执行完毕，正常退出")
	}()

	//提交清理计划
	defer func() {
		if err := toolkit.RemoveCacheDir(); err != nil {
			toolkit.ClearScreen()
			log.Println("缓存目录清理失败:", err)
			fmt.Println("缓存目录清理失败，确认需清理时可手动清理或重新运行程序")
		} else {
			fmt.Println("缓存目录清理完毕")
		}
	}()

	//初始化客户端
	if success := httpclient.Init(); !success {
		err := sso.HandleQRCodeLogin()
		if err != nil {
			fmt.Println("处理二维码登录失败:", err)
			return
		}
	}

	//主逻辑循环
	for {
		//获取用户BV号输入并检查
		BVid := toolkit.CatchAndCheckBVid()
		VideoInfoUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", BVid)

		//获取视频信息
		data, err := internal.CatchData(VideoInfoUrl)
		if err != nil {
			log.Println("获取视频信息数据错误: %v\n", err)
			fmt.Println("视频信息数据获取异常，请检查网络连接或前往log文件查看详情.")
			return
		}

		//视频信息反序列化
		Response, err := internal.ProcessResponse(data, 0)
		if err != nil {
			log.Println("处理视频详情发生错误: %v\n", err)
			fmt.Println("视频详情数据处理发生错误，请携带log文件向开发者反馈！")
			return
		}
		videoInfoResponse, ok := Response.(*toolkit.VideoInfoResponse)
		if !ok {
			log.Println("视频详情数据类型断言失败")
			fmt.Println("程序运行发生异常，请携带log日志文件联系开发者！")
			break
		}

		//打印视频详细信息并进行确认
		toolkit.ClearScreen()
		internal.ConfirmVideoExplanation(videoInfoResponse)

		Default := -2
		var actions map[string]map[string]int64
		if videoInfoResponse.Data.Ugc_season.Sections != nil || len(videoInfoResponse.Data.Pages) > 1 {
			var part int64
			//确认是否使用多分P选择/连续下载功能
			fmt.Printf("检测到视频含有分P，是否使用多分P选择/连续下载功能？(y/n):")
			if toolkit.YesOrNo() {
				Default = 0
				if len(videoInfoResponse.Data.Pages) > 1 {
					internal.PrintDiversityInformationPart2(videoInfoResponse)
					part = 2
				} else {
					internal.PrintDiversityInformationPart1(videoInfoResponse)
					part = 1
				}
				if actions = toolkit.GetMaps(videoInfoResponse, part); actions == nil {
					return
				}
			} else {
				actions = map[string]map[string]int64{
					videoInfoResponse.Data.Title: {
						videoInfoResponse.Data.Bvid: videoInfoResponse.Data.Cid,
					},
				}
			}
			//////缺少在分P列表中对当前视频的突出显示//////
		} else {
			actions = map[string]map[string]int64{
				videoInfoResponse.Data.Title: {
					videoInfoResponse.Data.Bvid: videoInfoResponse.Data.Cid,
				},
			}
		}

		for title, video := range actions {
			for bvid, cid := range video {
				//请求视频取流地址
				DownloadURL := fmt.Sprintf("https://api.bilibili.com/x/player/wbi/playurl?bvid=%s&cid=%d&fnval=4048", bvid, cid)
				data, err = internal.CatchData(DownloadURL)
				if err != nil {
					log.Println("获取下载信息数据发生错误: %v\n", err)
					fmt.Println("视频下载信息获取异常，请检查网络连接或前往log文件查看详情.")
					return
				}

				//反序列化视频流信息
				newResponse, err := internal.ProcessResponse(data, 1)
				if err != nil {
					log.Println("处理下载信息发生错误: %v\n", err)
					fmt.Println("视频下载信息处理发生错误，请携带log文件向开发者反馈！")
					return
				}
				downloadInfoResponse, ok := newResponse.(*toolkit.DownloadInfoResponse)
				if !ok {
					log.Println("视频下载数据类型断言失败")
					fmt.Println("程序运行发生异常，请携带log日志文件联系开发者！")
					break
				}

				//处理用户选择
				if len(downloadInfoResponse.Data.AcceptDescription) == 0 || downloadInfoResponse.Data.AcceptDescription[0] == `试看` {
					log.Println("当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址.")
					fmt.Println("\n当前账号可能没有观看（下载）该视频的权限，无法获取视频下载地址.")
				} else {
					if Default == 0 {
						fmt.Printf("当前为多分P下载模式，是否默认下载可获取的最高分辨率?(y/n)；")
						if toolkit.YesOrNo() {
							Default = 1
						} else {
							Default = -1
						}
					}
					videoIndex, _, resolutionDescription := toolkit.ObtainUserResolutionSelection(int64(Default), title, downloadInfoResponse)

					//请求视频下载
					if err := internal.DownloadFile(downloadInfoResponse.Data.Dash.Video[videoIndex].BackupURL[0], downloadInfoResponse.Data.Dash.Audio[0].BackupURL[0], ""); err != nil {
						log.Printf("请求下载失败：%s\n", err)
						fmt.Println("请求下载失败，请检查网络连接或前往log文件查看详情.")
					}

					//视频音频混流转码
					fmt.Println("开始视频转码：\n")
					video_processing.Transcoding(title, resolutionDescription)
					time.Sleep(500 * time.Millisecond)
				}
			}
		}

		fmt.Printf("是否继续下载其他视频？(y/n):")
		if isContinue := toolkit.YesOrNo(); isContinue {
			toolkit.ClearScreen()
			continue
		} else {
			break
		}
	}
}
