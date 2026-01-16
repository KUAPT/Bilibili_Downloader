package internal

import (
	"Bilibili_Downloader/internal/paths"
	"fmt"
	"log"
	"os"
)

func InitLog() *os.File {
	logPath := paths.Get().LogPath
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("程序运行，开始日志记录")
	return logFile
}

func CacheClean() {
	cacheDir := paths.Get().CacheDir
	if err := os.RemoveAll(cacheDir); err != nil {
		log.Println("缓存目录清理失败:", err)
		fmt.Println("缓存目录清理失败，确认需清理时可手动清理或重新运行程序")
	} else {
		fmt.Println("缓存目录清理完毕")
	}
}
