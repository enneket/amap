package weatherinfo

import (
	amapType "github.com/enneket/amap/types"
)

// WeatherinfoRequest 天气信息请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/weatherinfo

type WeatherinfoRequest struct {
	City       string              `json:"city"`                 // 城市ID或名称（必填）
	Extensions string              `json:"extensions,omitempty"` // 返回结果扩展（可选，base或all，默认base）
	Output     amapType.OutputType `json:"output,omitempty"`     // 输出格式（可选，默认JSON）
	Callback   string              `json:"callback,omitempty"`   // 回调函数（可选，用于JSONP跨域）
	Timestamp  string              `json:"timestamp,omitempty"`  // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *WeatherinfoRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["city"] = req.City // 城市为必填项，直接添加

	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Output != "" {
		params["output"] = string(req.Output)
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}

	return params
}
