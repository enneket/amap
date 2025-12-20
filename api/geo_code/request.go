package geo_code

// GeocodeRequest 地理编码请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo#t4
type GeocodeRequest struct {
	Address  string `json:"address"`            // 待解析的地址（必填）
	City     string `json:"city,omitempty"`     // 城市（可选）
	Callback string `json:"callback,omitempty"` // 回调函数（可选，用于JSONP跨域）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *GeocodeRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["address"] = req.Address // 地址为必填项，直接添加
	if req.City != "" {
		params["city"] = req.City
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	return params
}
