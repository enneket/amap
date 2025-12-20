package re_geo_code

import (
	"strconv"
)

// ReGeocodeRequest 逆地理编码请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo#t5
type ReGeocodeRequest struct {
	Location   string `json:"location"`   // 必填：经纬度（格式："经度,纬度"，如"116.481028,39.921983"）
	Poitype    string `json:"poitype"`    // 可选：POI类型过滤（如"餐饮|酒店"，仅extensions=all时生效）
	Radius     int    `json:"radius"`     // 可选：搜索半径（单位：米，默认1000，最大3000）
	Extensions string `json:"extensions"` // 可选：返回信息类型（默认"base"基础信息，"all"详细信息）
	RoadLevel  int    `json:"roadlevel"`  // 可选：是否返回道路等级（默认false，仅extensions=all时生效）
	HomeOrCorp string `json:"homeorcorp"` // 可选：是否优化 POI 返回顺序
	Callback   string `json:"callback"`   // 可选：JSONP回调函数名
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
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.RoadLevel > 0 {
		params["roadlevel"] = strconv.Itoa(req.RoadLevel)
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Poitype != "" {
		params["poitype"] = req.Poitype
	}
	return params
}
