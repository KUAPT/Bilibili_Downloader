package httpclient

import (
	"Bilibili_Downloader/pkg/cookie"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

// 定义一个全局的 http.Client 变量
var client *http.Client
var once sync.Once

// 初始化函数，创建并配置一个带有 cookiejar 的 http.Client
func Init() bool {
	success := true
	once.Do(func() {
		// 加载之前保存的 cookies
		cookies := cookie.LoadCookies()
		if cookies != nil {
			// 创建一个 cookie jar 并设置 cookies
			jar, _ := cookiejar.New(nil)
			URL, _ := url.Parse("https://api.bilibili.com/")
			jar.SetCookies(URL, cookies)

			// 使用带有 cookies 的 http.Client
			client = &http.Client{Jar: jar}

			// 使用 client 进行操作
		} else {
			success = false
		}
	})
	return success
}

func ChangeClinet(newClient *http.Client) {
	client = newClient
}

// 获取全局的 http.Client 实例
func GetClient() *http.Client {
	return client
}
