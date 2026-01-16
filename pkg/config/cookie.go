package config

import (
	"Bilibili_Downloader/internal/paths"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func StoreCookies(cookies []*http.Cookie) {
	p := paths.Get()
	if err := os.MkdirAll(p.ConfigDir, 0755); err != nil {
		log.Println("配置目录创建失败：", err)
	}

	cookiePath := p.CookieFilePath()
	file, err := os.Create(cookiePath)
	if err != nil {
		fmt.Println("创建Cookie文件失败:", err)
		log.Println("创建Cookie文件失败:", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("Close Cookie file失败：", err)
		}
	}()

	cookiesJSON, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		fmt.Println("转换 cookies 到 JSON 失败:", err)
		log.Println("转换 cookies 到 JSON 失败:", err)
		return
	}

	if err := os.WriteFile(file.Name(), cookiesJSON, 0644); err != nil {
		fmt.Println("写入 cookies 文件失败:", err)
		log.Println("写入 cookies 文件失败:", err)
		return
	}

	fmt.Println("Cookies 已保存到:", file.Name())
	log.Println("Cookies 已保存到:", file.Name())
}

func LoadCookies() []*http.Cookie {
	cookiePath := paths.Get().CookieFilePath()
	content, err := os.ReadFile(cookiePath)
	if err != nil {
		fmt.Println("未成功加载已保存的配置文件:", err)
		log.Println("未成功加载已保存的配置文件:", err)
		return nil
	}

	var cookies []*http.Cookie
	err = json.Unmarshal(content, &cookies)
	if err != nil {
		fmt.Println("解析 cookies 失败:", err)
		log.Println("解析 cookies 失败:", err)
		return nil
	}
	fmt.Println("Cookie加载成功，前十个字符为：", cookies[0].Value[:10])
	log.Println("Cookie加载成功")

	return cookies
}
