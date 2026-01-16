package main

import (
	"Bilibili_Downloader/internal"
	"Bilibili_Downloader/internal/app"
	"Bilibili_Downloader/internal/paths"
	"Bilibili_Downloader/internal/ui"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	flagBVID     = flag.String("bvid", "", "指定要下载的 BV 号")
	flagNonInter = flag.Bool("non-interactive", false, "非交互模式运行")
	flagNoUpdate = flag.Bool("no-update", false, "禁用更新检查")
	flagOutput   = flag.String("output", "", "指定输出路径")
	flagVersion  = flag.Bool("version", false, "显示版本号并退出")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println("Bilibili_Downloader v1.4.2")
		os.Exit(0)
	}

	if err := paths.Init(); err != nil {
		fmt.Println("初始化路径失败:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logFile := internal.InitLog()
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Println("log文件close失败:", err)
		}
	}()

	defer internal.CacheClean()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		s := <-signalChan
		log.Printf("收到信号：%v，开始退出程序...\n", s)
		cancel()
		_ = os.Stdin.Close()
	}()

	var userInterface ui.UI
	nonInteractive := *flagNonInter

	if nonInteractive {
		if *flagBVID == "" {
			fmt.Println("错误: --non-interactive 模式必须提供 --bvid 参数")
			flag.Usage()
			os.Exit(1)
		}
		userInterface = ui.NewNonInteractiveUI()
	} else {
		userInterface = ui.NewStdIOUI()
	}

	if !app.InitClient() {
		if nonInteractive {
			fmt.Println("错误: Cookie 未配置或已过期，非交互模式无法登录")
			os.Exit(1)
		}
		if err := app.HandleLogin(); err != nil {
			fmt.Println("处理二维码登录失败:", err)
			return
		}
	}

	cfg := app.Config{
		BVID:       *flagBVID,
		NoUpdate:   *flagNoUpdate || nonInteractive,
		NonInter:   nonInteractive,
		OutputPath: *flagOutput,
	}

	application := app.New(userInterface, cfg)

	defer func() {
		if ctx.Err() == nil && userInterface.IsInteractive() {
			userInterface.WaitForExit()
		}
		log.Println("程序执行完毕，正常退出")
	}()

	if err := application.Run(ctx); err != nil {
		if err != context.Canceled {
			log.Printf("程序运行错误: %v\n", err)
			userInterface.ShowError(fmt.Sprintf("程序运行错误: %v", err))
		}
	}
}
