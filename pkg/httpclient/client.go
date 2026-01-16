package httpclient

import (
	"Bilibili_Downloader/pkg/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)

// 定义一个全局的 http.Client 变量
var client *http.Client
var once sync.Once
var initialized bool

// Init 初始化函数，创建并配置一个带有 cookiejar 的 http.Client
func Init() bool {
	success := true
	once.Do(func() {
		// 加载之前保存的 cookies
		cookies := config.LoadCookies()
		// 创建一个 cookie jar
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Println("创建cookie jar失败:", err)
			success = false
			return
		}

		if cookies != nil {
			// 设置 cookies
			biliURL, _ := url.Parse("https://api.bilibili.com/")
			jar.SetCookies(biliURL, cookies)
		} else {
			log.Println("未找到保存的cookies")
			success = false
		}

		// 创建HTTP客户端，设置超时
		client = &http.Client{
			Jar:     jar,
			Timeout: DefaultTimeout,
		}
		initialized = true
	})
	return success
}

// ChangeClient 修改全局HTTP客户端
func ChangeClient(newClient *http.Client) {
	client = newClient
	initialized = true
}

// GetClient 获取全局的 http.Client 实例
func GetClient() *http.Client {
	if !initialized {
		Init()
	}
	if client == nil {
		// 若初始化失败，返回一个默认客户端
		jar, _ := cookiejar.New(nil)
		client = &http.Client{
			Jar:     jar,
			Timeout: DefaultTimeout,
		}
		log.Println("使用默认HTTP客户端")
		fmt.Println("警告：未加载Cookie，部分功能可能受限")
	}
	return client
}

// Do 执行HTTP请求，支持 context
func Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return GetClient().Do(req.WithContext(ctx))
}

// Get 执行GET请求，支持 context
func Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return GetClient().Do(req)
}
