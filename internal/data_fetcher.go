package internal

import (
	"Bilibili_Downloader/pkg/httpclient"
	"Bilibili_Downloader/pkg/toolkit"
	"Bilibili_Downloader/pkg/toolkit/data_struct"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CatchData(url string) ([]byte, error) {
	client := httpclient.GetClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	toolkit.SetBilibiliHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Close resp.Body失败:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取 URL 失败，状态码: %d", resp.StatusCode)
	}
	log.Printf("请求响应成功，状态码：%v\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}

func ProcessResponse(data []byte, flag int) (interface{}, error) {
	var err error
	if flag == 0 {
		var response data_struct.VideoInfoResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return nil, fmt.Errorf("视频信息解组失败: %w", err)
		}
		fmt.Println("视频信息数据解组正常！")
		log.Println("视频信息数据解组正常")
		return &response, nil
	}
	if flag == 1 {
		var response data_struct.DownloadInfoResponse
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
