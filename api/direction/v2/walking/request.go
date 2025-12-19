package walking

import (
	amapType "github.com/enneket/amap/types"
)

// WalkingRequestV2 步行路线规划v2请求参数
type WalkingRequestV2 struct {
	Origin          string                  `json:"origin"`                   // 起点坐标（必填，格式：经度,纬度）
	Destination     string                  `json:"destination"`              // 终点坐标（必填，格式：经度,纬度）
	Strategy        string                  `json:"strategy,omitempty"`       // 步行策略（可选，默认0=最快捷）
	CoordinateType  amapType.CoordinateType `json:"coordinate_type,omitempty"` // 输入/输出坐标系（可选，默认gcj02）
	Output          amapType.OutputType     `json:"output,omitempty"`          // 输出格式（可选，默认JSON）
	Language        amapType.LanguageType   `json:"language,omitempty"`        // 语言（可选，默认中文）
	Timestamp       string                  `json:"timestamp,omitempty"`       // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *WalkingRequestV2) ToParams() map[string]string {
	params := make(map[string]string)
	params["origin"] = req.Origin     // 起点坐标为必填项
	params["destination"] = req.Destination // 终点坐标为必填项
	if req.Strategy != "" {
		params["strategy"] = req.Strategy
	}
	if req.CoordinateType != "" {
		params["coordinate_type"] = string(req.CoordinateType)
	}
	if req.Output != "" {
		params["output"] = string(req.Output)
	}
	if req.Language != "" {
		params["language"] = string(req.Language)
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}
	return params
}
