package convert

// ConvertRequest 坐标转换请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/convert
// 支持将其他坐标系的坐标转换为高德坐标系（GCJ02）
// 支持批量转换，一次最多转换40对坐标

type ConvertRequest struct {
	Locations string `json:"locations"`  // 待转换的坐标列表（必填，格式："经度,纬度;经度,纬度"，最多40对）
	CoordSys  string `json:"coordsys"`   // 原坐标系（必填，可选值：gps, mapbar, baidu, autonavi）
	Output    string `json:"output,omitempty"`  // 输出格式（可选，默认JSON）
	Callback  string `json:"callback,omitempty"` // 回调函数（可选，用于JSONP跨域）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *ConvertRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["locations"] = req.Locations // 坐标列表为必填项，直接添加
	params["coordsys"] = req.CoordSys   // 原坐标系为必填项，直接添加
	if req.Output != "" {
		params["output"] = req.Output
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	return params
}
