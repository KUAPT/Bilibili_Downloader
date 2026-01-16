package toolkit

import (
	"fmt"
	"os"
	"strings"
)

// CheckAndCreateDir 检查并创建指定目录
func CheckAndCreateDir(dir string) error {
	configDir := dir

	// 检查目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			return fmt.Errorf("无法创建目录 %s: %v", configDir, err)
		}
		fmt.Println("目录已创建:", configDir)
	} else if err != nil {
		// 其他错误
		return fmt.Errorf("检查目录 %s 时出错: %v", configDir, err)
	} else {
		// 目录已存在
		fmt.Println("目录已存在:", configDir)
	}
	return nil
}

// CheckAndCleanFileName 检查文件名是否包含不允许的字符，并进行清理
func CheckAndCleanFileName(fileName string) string {
	disallowedChars := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	// 检查文件名中的每个字符
	for _, char := range disallowedChars {
		if strings.Contains(fileName, char) {
			// 替换不允许的字符为下划线
			fileName = strings.ReplaceAll(fileName, char, "_")
		}
	}
	return fileName
}
