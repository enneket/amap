package utils

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
)

// Sign 计算高德API签名 https://lbs.amap.com/faq/quota-key/key/41181
// params: 参与签名的参数（key-value），空值参数会被自动过滤
// securityKey: 高德开发者平台的安全密钥
// 返回值: 大写的MD5签名串
func Sign(params map[string]string, securityKey string) string {
	// 1. 过滤空值参数（空字符串不参与签名）
	validParams := make(map[string]string, len(params))
	for k, v := range params {
		if v != "" {
			validParams[k] = v
		}
	}

	// 2. 按参数名升序排序
	keys := make([]string, 0, len(validParams))
	for k := range validParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 3. 拼接参数为 "key=value" 格式
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(validParams[k])
	}

	// 4. 拼接安全密钥并计算MD5
	sb.WriteString(securityKey)
	hash := md5.Sum([]byte(sb.String()))

	// 5. 转换为大写十六进制字符串
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}
