package re_geo_code

import (
	"strconv"

	amapType "github.com/enneket/amap/types"
)

// ReGeocodeRequest 逆地理编码请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo/#regeo
type ReGeocodeRequest struct {
	Location       string                  `json:"location"`        // 必填：经纬度（格式："经度,纬度"，如"116.481028,39.921983"）
	Radius         int                     `json:"radius"`          // 可选：搜索半径（单位：米，默认1000，最大3000）
	CoordinateType amapType.CoordinateType `json:"coordinate_type"` // 可选：输入坐标系（默认gcj02，支持wgs84/bd09ll）
	Extensions     string                  `json:"extensions"`      // 可选：返回信息类型（默认"base"基础信息，"all"详细信息）
	Output         amapType.OutputType     `json:"output"`          // 可选：响应格式（默认JSON，支持XML）
	Language       amapType.LanguageType   `json:"language"`        // 可选：响应语言（默认zh_cn，支持en）
	Callback       string                  `json:"callback"`        // 可选：JSONP回调函数名
	Timestamp      string                  `json:"timestamp"`       // 可选：时间戳（核心客户端已自动填充，可覆盖）
	Poitype        string                  `json:"poitype"`         // 可选：POI类型过滤（如"餐饮|酒店"，仅extensions=all时生效）
	HouseNumber    string                  `json:"housenumber"`     // 可选：是否返回门牌号（"true"/"false"，默认false，仅extensions=all时生效）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *ReGeocodeRequest) ToParams() map[string]string {
	params := make(map[string]string)
	if req.Location != "" {
		params["location"] = req.Location
	}
	if req.Radius > 0 {
		params["radius"] = strconv.Itoa(req.Radius)
	}
	if req.CoordinateType != "" {
		params["coordinate_type"] = string(req.CoordinateType)
	}
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
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
	if req.Poitype != "" {
		params["poitype"] = req.Poitype
	}
	if req.HouseNumber != "" {
		params["housenumber"] = req.HouseNumber
	}
	return params
}
