package utils

import (
	"net/url"
	"sort"
	"strings"
)

// EncodeParams 将参数map转换为URL编码的查询字符串
// params: 待编码的参数（key-value）
// sorted: 是否按key升序排列（建议true，与签名逻辑一致）
// 返回值: 编码后的查询字符串（如 "address=北京市&city=北京"）
func EncodeParams(params map[string]string, sorted bool) string {
	values := url.Values{}

	// 1. 过滤空值并添加到url.Values
	for k, v := range params {
		if v != "" {
			values.Add(k, v)
		}
	}

	// 2. 如果需要排序，手动拼接（url.Values.Encode() 不保证顺序）
	if sorted {
		keys := make([]string, 0, len(values))
		for k := range values {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var encodedParts []string
		for _, k := range keys {
			// 对key和value分别编码
			encodedKey := url.QueryEscape(k)
			encodedValue := url.QueryEscape(values.Get(k))
			encodedParts = append(encodedParts, encodedKey+"="+encodedValue)
		}
		return strings.Join(encodedParts, "&")
	}

	// 3. 不排序则直接使用内置方法
	return values.Encode()
}

// EncodeParamsDefault 默认编码（按key升序）
func EncodeParamsDefault(params map[string]string) string {
	return EncodeParams(params, true)
}
