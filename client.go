package amap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/enneket/amap/utils"
)

// Client 高德 API 客户端
type Client struct {
	config     *Config      // 全局配置
	httpClient *http.Client // 复用的 HTTP 客户端
}

// NewClient 创建客户端实例（校验配置合法性）
func NewClient(cfg *Config) (*Client, error) {
	if cfg.Key == "" {
		return nil, NewInvalidConfigError("API Key 不能为空")
	}
	// 初始化 HTTP 客户端（支持超时、代理）
	httpClient := &http.Client{Timeout: cfg.Timeout}
	if cfg.Proxy != "" {
		proxyURL, _ := url.Parse(cfg.Proxy)
		httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	return &Client{config: cfg, httpClient: httpClient}, nil
}

// DoRequest 通用请求方法（封装公共参数、签名、响应解析）
func (c *Client) DoRequest(path string, params map[string]string, resp interface{}) error {
	// 1. 合并公共参数（Key、签名、Timestamp 等）
	allParams := c.buildPublicParams(params)
	// 2. 签名（如果配置了 SecurityKey）
	if c.config.SecurityKey != "" {
		allParams["sig"] = utils.Sign(allParams, c.config.SecurityKey)
	}
	// 3. 构建请求 URL
	fullURL := c.config.BaseURL + path + "?" + utils.EncodeParams(allParams, true)
	// 4. 发送 HTTP 请求
	req, _ := http.NewRequest(http.MethodGet, fullURL, nil)
	req.Header.Set("User-Agent", c.config.UserAgent)
	rawResp, err := c.httpClient.Do(req)
	if err != nil {
		return NewNetworkError(err.Error())
	}
	defer rawResp.Body.Close()
	// 5. 解析响应（先解析基础响应，再解析业务响应）
	var baseResp BaseResponse
	if err := json.NewDecoder(rawResp.Body).Decode(&baseResp); err != nil {
		return NewParseError("响应解析失败: " + err.Error())
	}
	// 6. 校验 API 错误
	if baseResp.Status != "1" {
		return NewAPIError(baseResp.InfoCode, baseResp.Info)
	}
	// 7. 解析到业务响应结构体
	return json.Unmarshal(baseResp.RawJSON, resp)
}

// buildPublicParams 构建公共参数（Key、Timestamp 等）
func (c *Client) buildPublicParams(params map[string]string) map[string]string {
	publicParams := map[string]string{
		"key":       c.config.Key,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"output":    "JSON",
	}
	// 合并业务参数
	for k, v := range params {
		if v != "" { // 忽略空参数
			publicParams[k] = v
		}
	}
	return publicParams
}
