package detail

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProcessResponse 处理 JSON 响应并返回 Response 结构体
func ProcessResponse(data []byte, flag int) (interface{}, error) {
	var err error
	if flag == 0 {
		var response VideoInfoResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("JSON 解组失败: %w", err)
		}
		fmt.Println("数据解组正常！")
		return &response, nil
	} else if flag == 1 {
		var response DownloadInfoResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("JSON 解组失败: %w", err)
		}
		fmt.Println("数据解组正常！")
		return &response, nil
	}
	return nil, fmt.Errorf("不支持的 flag 值: %d", flag)
	//fmt.Printf("解组后的数据：%v", response)
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
		err = os.Mkdir(cacheDirPath, 0755)
		if err != nil {
			return fmt.Errorf("临时下载目录创建失败: %v", err)
		}
		fmt.Println("成功创建目录:", cacheDirPath)
	} else {
		// 目录存在，报错
		return fmt.Errorf("目录 '%s' 已经存在", cacheDirPath)
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
	return nil
}

// renameFile 修改文件名（不包含扩展名）
func RenameFile(filepath, newName string) (string, error) {
	// 分离文件名和扩展名
	dir := filepath.Dir(filepath)
	ext := filepath.Ext(filepath)
	oldName := strings.TrimSuffix(filepath.Base(filepath), ext)

	// 生成新的文件路径
	newFilepath := filepath.Join(dir, newName+ext)

	// 重命名文件
	err := os.Rename(filepath, newFilepath)
	if err != nil {
		return "", err
	}
	return newFilepath, nil
}
