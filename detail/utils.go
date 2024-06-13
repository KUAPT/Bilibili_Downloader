package detail

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ProcessResponse 处理 JSON 响应并返回 Response 结构体
func ProcessResponse(data []byte, flag int) (interface{}, error) {
	var err error
	if flag == 0 {
		var response VideoInfoResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("视频信息解组失败: %w", err)
		}
		fmt.Println("视频信息数据解组正常！")
		log.Println("视频信息数据解组正常")
		return &response, nil
	} else if flag == 1 {
		var response DownloadInfoResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("下载信息解组失败: %w", err)
		}
		fmt.Println("下载信息数据解组正常！")
		log.Println("下载信息数据解组正常")
		return &response, nil
	}
	return nil, fmt.Errorf("不支持的 flag 值: %d", flag)
}

func CheckAndCreateCacheDir() error {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	cacheDirPath := filepath.Join(currentDir, "download_cache")

	// 检查目录是否存在
	if _, err := os.Stat(cacheDirPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err = os.MkdirAll(cacheDirPath, 0755)
		if err != nil {
			return fmt.Errorf("临时下载目录创建失败: %v", err)
		}
		fmt.Println("成功创建缓存目录:", cacheDirPath)
		log.Println("建立缓存目录正常:", cacheDirPath)
	} else if err != nil {
		// 如果 os.Stat 返回了错误，但不是 os.IsNotExist
		return fmt.Errorf("检查目录时发生错误: %v", err)
	} else {
		// 目录存在
		fmt.Println("缓存目录已经存在，继续使用:", cacheDirPath)
		log.Println("缓存目录已经存在，继续使用:", cacheDirPath)
	}

	return nil
}

func RemoveCacheDir() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %v", err)
	}

	cacheDirPath := filepath.Join(currentDir, "download_cache")

	if err := os.RemoveAll(cacheDirPath); err != nil {
		return fmt.Errorf("移除临时cache目录失败: %v", err)
	}

	fmt.Println("成功移除cache目录:", cacheDirPath)
	log.Println("成功移除cache目录:", cacheDirPath)
	return nil
}

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		log.Println("无法清屏，不支持的平台！")
	}
}
