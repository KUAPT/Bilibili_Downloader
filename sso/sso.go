package sso

import (
	"Bilibili_Downloader/cookie"
	"Bilibili_Downloader/httpclient"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// RequestQRCode 请求并获取二维码信息
func RequestQRCode(client *http.Client) (string, string, error) {
	resp, err := client.Get("https://passport.bilibili.com/x/passport-login/web/qrcode/generate")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var qrResp QRCodeResponse
	err = json.Unmarshal(body, &qrResp)
	if err != nil {
		return "", "", err
	}

	if qrResp.Code != 0 {
		return "", "", fmt.Errorf("申请二维码失败: %s", qrResp.Message)
	}

	return qrResp.Data.QRCodeKey, qrResp.Data.URL, nil
}

// PollQRCodeStatus 轮询二维码状态
func PollQRCodeStatus(client *http.Client, token string) (int, []*http.Cookie, error) {
	qrStatusURL := fmt.Sprintf("https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=%s", url.QueryEscape(token))
	resp, err := client.Get(qrStatusURL)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	var pollResp PollingResponse
	err = json.Unmarshal(body, &pollResp)
	if err != nil {
		return 0, nil, err
	}

	if pollResp.Code != 0 {
		return pollResp.Data.Code, nil, fmt.Errorf("轮询扫码状态失败: %s", pollResp.Message)
	}

	return pollResp.Data.Code, resp.Cookies(), nil
}

// HandleQRCodeLogin 处理二维码登录流程
func HandleQRCodeLogin() error {
	// 创建 cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("创建 cookie jar 失败: %v", err)
	}
	client := &http.Client{Jar: jar}
	httpclient.ChangeClinet(client)

	token, qrURL, err := RequestQRCode(client)
	if err != nil {
		return err
	}

	err = DisplayQRCodeInTerminal(qrURL)
	if err != nil {
		return err
	}

	for {
		status, cookies, err := PollQRCodeStatus(client, token)
		if err != nil {
			return err
		}

		switch status {
		case 86101: // 未扫码
			time.Sleep(2 * time.Second)
			continue
		case 86038: // 二维码超时或失效
			return fmt.Errorf("二维码失效或超时")
		case 86090: // 已扫描未确认
			fmt.Println("二维码已扫描，等待确认")
		case 0: // 登录成功
			fmt.Println("扫描成功")
			cookie.StoreCookies(cookies)
			return nil
		default:
			return fmt.Errorf("未知状态码: %d", status)
		}
	}
}
