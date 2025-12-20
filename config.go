package amap

import "time"

// Config 高德 API 全局配置
type Config struct {
	Key         string        // 高德 API Key（必填）
	SecurityKey string        // 安全密钥（可选，用于签名）
	Timeout     time.Duration // 请求超时（默认 5s）
	Proxy       string        // HTTP 代理地址（可选）
	UserAgent   string        // 请求 UA（默认 amap-go/1.0）
}

// NewConfig 创建默认配置（只需传入必填的 Key）
func NewConfig(key string) *Config {
	return &Config{
		Key:       key,
		Timeout:   5 * time.Second,
		UserAgent: "amap-go/1.0",
	}
}
