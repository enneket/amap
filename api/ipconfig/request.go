package ipconfig

import (
	amapType "github.com/enneket/amap/types"
)

// IPConfigRequest IP定位请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/ipconfig
// 支持通过IP地址查询地理位置信息，可指定返回结果类型

type IPConfigRequest struct {
	IP         string                  `json:"ip"`                  // IP地址（必填，IPv4/IPv6格式）
	Output     amapType.OutputType     `json:"output,omitempty"`     // 输出格式（可选，默认JSON）
	Language   amapType.LanguageType   `json:"language,omitempty"`   // 语言（可选，默认中文）
	Callback   string                  `json:"callback,omitempty"`   // 回调函数（可选，用于JSONP跨域）
	Timestamp  string                  `json:"timestamp,omitempty"`  // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *IPConfigRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["ip"] = req.IP // IP地址为必填项
	if req.Output != "" {
		params["output"] = string(req.Output)
	}
	if req.Language != "" {
		params["language"] = string(req.Language)
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}
	return params
}
