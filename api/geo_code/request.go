package geo_code

import (
	amapType "github.com/enneket/amap/types"
)

// GeocodeRequest 地理编码请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo/#geo
type GeocodeRequest struct {
	Address        string                  `json:"address"`                   // 待解析的地址（必填）
	City           string                  `json:"city,omitempty"`            // 城市（可选）
	CoordinateType amapType.CoordinateType `json:"coordinate_type,omitempty"` // 输出坐标系（可选，默认gcj02，支持wgs84/bd09ll）
	Output         amapType.OutputType     `json:"output,omitempty"`          // 输出格式（可选，默认JSON）
	Language       amapType.LanguageType   `json:"language,omitempty"`        // 语言（可选，默认中文）
	Callback       string                  `json:"callback,omitempty"`        // 回调函数（可选，用于JSONP跨域）
	Timestamp      string                  `json:"timestamp,omitempty"`       // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *GeocodeRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["address"] = req.Address // 地址为必填项，直接添加
	if req.City != "" {
		params["city"] = req.City
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
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}
	return params
}
