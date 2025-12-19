package geo_code

import (
	appType "github.com/enneket/amap/types"
)

// GeocodeRequest 地理编码请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo/#geo
type GeocodeRequest struct {
	Address        string                 `json:"address"`         // 必填：详细地址（如"北京市朝阳区望京SOHO"）
	City           string                 `json:"city"`            // 可选：城市名称（缩小搜索范围，如"北京"）
	CoordinateType appType.CoordinateType `json:"coordinate_type"` // 可选：输出坐标系（默认gcj02，支持wgs84/bd09ll）
	Output         appType.OutputType     `json:"output"`          // 可选：响应格式（默认JSON，支持XML）
	Language       appType.LanguageType   `json:"language"`        // 可选：响应语言（默认zh_cn，支持en）
	Callback       string                 `json:"callback"`        // 可选：JSONP回调函数名（前端跨域使用）
	Timestamp      string                 `json:"timestamp"`       // 可选：时间戳（核心客户端已自动填充，可覆盖）
}

// ToParams 转换为请求参数map（过滤空值）
func (r *GeocodeRequest) ToParams() map[string]string {
	params := make(map[string]string)
	if r.Address != "" {
		params["address"] = r.Address
	}
	if r.City != "" {
		params["city"] = r.City
	}
	if r.CoordinateType != "" {
		params["coordinate_type"] = string(r.CoordinateType)
	}
	if r.Output != "" {
		params["output"] = string(r.Output)
	}
	if r.Language != "" {
		params["language"] = string(r.Language)
	}
	if r.Callback != "" {
		params["callback"] = r.Callback
	}
	if r.Timestamp != "" {
		params["timestamp"] = r.Timestamp
	}
	return params
}
