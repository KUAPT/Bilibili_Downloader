package sso

// QRCodeResponse 定义二维码生成接口的响应结构
type QRCodeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		URL       string `json:"url"`
		QRCodeKey string `json:"qrcode_key"`
	} `json:"data"`
}

// PollingResponse 定义扫码状态接口的响应结构
type PollingResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		URL       string `json:"url"`
		Code      int    `json:"code"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}
