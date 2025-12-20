package traffic_incident

import (
	amapType "github.com/enneket/amap/types"
)

// TrafficIncidentRequest 交通事件查询请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/traffic-incident
// 支持查询指定区域内的交通事件，可按事件级别和类型筛选
// 矩形区域参数格式：左下经纬度,右上经纬度（如：116.351147,39.904989,116.480317,39.976564）

type TrafficIncidentRequest struct {
	Level       string                  `json:"level"`                 // 事件级别（必填，1-4，默认1：所有级别）
	Type        string                  `json:"type"`                  // 事件类型（必填，1-12，多个用|分隔，默认所有类型）
	Rectangle   string                  `json:"rectangle"`             // 查询区域（必填，左下经纬度,右上经纬度）
	Extensions  string                  `json:"extensions,omitempty"`  // 返回结果类型（可选，base/all，默认base）
	Output      amapType.OutputType     `json:"output,omitempty"`      // 输出格式（可选，默认JSON）
	Callback    string                  `json:"callback,omitempty"`    // 回调函数（可选，用于JSONP跨域）
	Language    amapType.LanguageType   `json:"language,omitempty"`    // 语言（可选，默认中文）
	Timestamp   string                  `json:"timestamp,omitempty"`   // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *TrafficIncidentRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["level"] = req.Level     // 事件级别为必填项
	params["type"] = req.Type       // 事件类型为必填项
	params["rectangle"] = req.Rectangle // 查询区域为必填项
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Output != "" {
		params["output"] = string(req.Output)
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Language != "" {
		params["language"] = string(req.Language)
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}
	return params
}
