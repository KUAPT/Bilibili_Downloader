package detail

import (
	"Bilibili_Downloader/httpclient"
	"Bilibili_Downloader/tool"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
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

func DownloadFile(urlVideo string, urlAudio string, filepath string) error {
	var filepath1, filepath2 string

	// 使用正则表达式提取文件扩展名
	re := regexp.MustCompile(`\.([a-zA-Z0-9]+)$`)
	match1 := re.FindStringSubmatch(path.Base(urlVideo))
	match2 := re.FindStringSubmatch(path.Base(urlAudio))
	// 提取到的文件扩展名
	var videoExtension, audioExtension string
	if len(match1) > 1 {
		videoExtension = match1[1]
	}
	if len(match2) > 1 {
		audioExtension = match2[1]
	}

	if filepath == "" {
		if err := tool.CheckAndCreateCacheDir(); err != nil {
			fmt.Println("检查并创建临时下载目录失败")
		}
		// 使用提取的文件名和扩展名创建缓存路径
		filepath1 = "./download_cache/audio_cache" + audioExtension
		filepath2 = "./download_cache/video.cache" + videoExtension
	} else {
		// 检查字符串末尾是否已经有斜杠
		if !strings.HasSuffix(filepath, "/") {
			// 如果没有，则在末尾添加斜杠
			filepath += "/"
		}
		// 如果指定了文件路径，则在文件路径后添加适当的扩展名
		filepath1 = filepath + "audio_cache" + audioExtension
		filepath2 = filepath + "video.cache" + videoExtension
	}

	//client := &http.Client{}
	client := httpclient.GetClient()
	req1, err := http.NewRequest("GET", urlAudio, nil)
	if err != nil {
		return err
	}
	req2, err := http.NewRequest("GET", urlVideo, nil)
	if err != nil {
		return err
	}

	// 设置自定义请求头
	req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req1.Header.Set("Accept", "*/*")
	//req1.Header.Set("Connection", "keep-alive")
	req1.Header.Set("Referer", "https://www.bilibili.com/vedio")
	req1.Header.Set("Origin", "https://www.bilibili.com")
	// 设置自定义请求头
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0")
	req2.Header.Set("Accept", "*/*")
	//req2.Header.Set("Connection", "keep-alive")
	req2.Header.Set("Referer", "https://www.bilibili.com/vedio")
	req2.Header.Set("Origin", "https://www.bilibili.com")

	resp1, err := client.Do(req1)
	if err != nil {
		return err
	}
	resp2, err := client.Do(req2)
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	fmt.Println("发送下载请求")
	// 检查HTTP响应状态码
	if resp1.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status of audio: %s", resp1.Status)
	}
	if resp2.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status of video: %s", resp2.Status)
	}

	fmt.Println("正在下载，请耐心等待...")
	log.Println("视频下载开始")

	// 创建文件
	out1, err := os.Create(filepath1)
	if err != nil {
		return err
	}
	defer out1.Close()

	// 将HTTP响应体内容写入文件
	_, err = io.Copy(out1, resp1.Body)
	if err != nil {
		return err
	}

	// 创建文件
	out2, err := os.Create(filepath2)
	if err != nil {
		return err
	}
	defer out2.Close()

	// 将HTTP响应体内容写入文件
	_, err = io.Copy(out2, resp2.Body)
	if err != nil {
		return err
	}

	tool.ClearScreen()
	fmt.Println("下载完成！")
	log.Println("视频下载成功")
	return nil
}

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
