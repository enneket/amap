package line

import (
	amapType "github.com/enneket/amap/types"
)

// LineTrafficRequest 指定线路交通态势查询请求参数
type LineTrafficRequest struct {
	Path            string                  `json:"path"`                    // 路径坐标串（必填，格式：lng1,lat1;lng2,lat2;...）
	Level           string                  `json:"level,omitempty"`         // 路况等级（可选，默认all）
	CoordinateType  amapType.CoordinateType `json:"coordinate_type,omitempty"` // 输入/输出坐标系（可选，默认gcj02）
	Output          amapType.OutputType     `json:"output,omitempty"`          // 输出格式（可选，默认JSON）
	Language        amapType.LanguageType   `json:"language,omitempty"`        // 语言（可选，默认中文）
	Timestamp       string                  `json:"timestamp,omitempty"`       // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *LineTrafficRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["path"] = req.Path // 路径坐标串为必填项
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
