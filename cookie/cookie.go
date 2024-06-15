package cookie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// 存储Cookie
func StoreCookies(cookies []*http.Cookie) {
	// 创建一个临时文件来存储 cookies
	file, err := os.CreateTemp("", "cookies*.json")
	if err != nil {
		fmt.Println("创建临时文件失败:", err)
		return
	}
	defer file.Close()

	// 将 cookies 转换为 JSON 格式
	cookiesJSON, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		fmt.Println("转换 cookies 到 JSON 失败:", err)
		return
	}

	// 将 JSON 写入文件
	if err := ioutil.WriteFile(file.Name(), cookiesJSON, 0644); err != nil {
		fmt.Println("写入 cookies 文件失败:", err)
		return
	}

	fmt.Println("Cookies 已保存到:", file.Name())
}

// 加载之前保存的 cookies
func LoadCookies() []*http.Cookie {
	// 读取之前保存的 cookies 文件
	content, err := ioutil.ReadFile("path/to/cookies.json")
	if err != nil {
		fmt.Println("读取 cookies 文件失败:", err)
		return nil
	}

	var cookies []*http.Cookie
	err = json.Unmarshal(content, &cookies)
	if err != nil {
		fmt.Println("解析 cookies 失败:", err)
		return nil
	}

	return cookies
}
