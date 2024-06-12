package video_processing

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed ffmpeg/ffmpeg.exe
var embeddedFFmpeg embed.FS

func extractFFmpeg() (string, error) {
	// 读取嵌入的ffmpeg二进制文件
	data, err := embeddedFFmpeg.ReadFile("ffmpeg/ffmpeg.exe")
	if err != nil {
		return "", err
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ffmpeg")
	if err != nil {
		return "", err
	}

	// 写入ffmpeg二进制文件
	ffmpegPath := filepath.Join(tempDir, "ffmpeg.exe")
	err = os.WriteFile(ffmpegPath, data, 0644)
	if err != nil {
		return "", err
	}

	return ffmpegPath, nil
}

func Transcoding(videoName string) {
	// 提取ffmpeg
	ffmpegPath, err := extractFFmpeg()
	if err != nil {
		fmt.Println("Error extracting ffmpeg:", err)
		return
	}
	defer os.RemoveAll(filepath.Dir(ffmpegPath))

	// 定义要处理的文件目录和输出的MP4文件名
	inputDir := "./download_cache"
	outputFile := "Output.mp4"

	// 获取所有的cache文件
	caches, err := filepath.Glob(filepath.Join(inputDir, "*"))
	if err != nil {
		fmt.Println("Error finding Cache files:", err)
		return
	}

	if len(caches) == 0 {
		fmt.Println("No cache files found")
		return
	}

	// 创建文件列表
	fileList := "filelist.txt"
	fileListContent := ""
	for _, file := range caches {
		fileListContent += fmt.Sprintf("file '%s'\n", file)
	}

	err = os.WriteFile(fileList, []byte(fileListContent), 0644)
	if err != nil {
		fmt.Println("Error creating file list:", err)
		return
	}
	defer os.Remove(fileList)

	// 创建ffmpeg命令
	cmd := exec.Command(ffmpegPath, "-f", "concat", "-safe", "0", "-i", fileList, "-c", "copy", outputFile)

	// 将输出重定向到标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行ffmpeg命令
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running ffmpeg:", err)
		return
	}

	fmt.Println("\n视频文件转码成功！")
}
