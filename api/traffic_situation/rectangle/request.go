package rectangle

import (
	amapType "github.com/enneket/amap/types"
)

// RectangleTrafficRequest 矩形区域内交通态势查询请求参数
type RectangleTrafficRequest struct {
	Rectangle       string                  `json:"rectangle"`                // 矩形区域（必填，格式：左下角lng,左下角lat;右上角lng,右上角lat）
	Level           string                  `json:"level,omitempty"`         // 路况等级（可选，默认all）
	CoordinateType  amapType.CoordinateType `json:"coordinate_type,omitempty"` // 输入/输出坐标系（可选，默认gcj02）
	Output          amapType.OutputType     `json:"output,omitempty"`          // 输出格式（可选，默认JSON）
	Language        amapType.LanguageType   `json:"language,omitempty"`        // 语言（可选，默认中文）
	Timestamp       string                  `json:"timestamp,omitempty"`       // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *RectangleTrafficRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["rectangle"] = req.Rectangle // 矩形区域为必填项
	if req.Level != "" {
		params["level"] = req.Level
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
