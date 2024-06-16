package detail

import (
	"Bilibili_Downloader/httpclient"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func CatchData(Url string) ([]byte, error) {

	//client := &http.Client{}
	client := httpclient.GetClient()
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}

	// 设置自定义请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.Header.Set("Origin", "https://www.bilibili.com/")
	//req.Header.Set("Cookie", Cookie)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取 URL 失败，状态码: %d", resp.StatusCode)
	}
	fmt.Printf("响应成功，状态码：%v\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}

func DownloadFile(url string, filepath string) error {
	if filepath == "" {
		if err := CheckAndCreateCacheDir(); err != nil {
			fmt.Println("检查并创建临时下载目录失败")
		}
		filepath = "./download_cache/video_cache.mp4"
	}

	//client := &http.Client{}
	client := httpclient.GetClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// 设置自定义请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.Header.Set("Origin", "https://www.bilibili.com/")
	//req.Header.Set("Cookie", Cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("正在下载，请耐心等待...")
	log.Println("视频下载开始")

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将HTTP响应体内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("下载完成！")
	log.Println("视频下载成功")
	return nil
}
