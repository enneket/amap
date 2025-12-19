package search

// AOISearchRequest AOI边界查询请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/search
// 查询AOI（兴趣面）的边界信息，用于获取指定区域的边界坐标

type AOISearchRequest struct {
	ID             string `json:"id,omitempty"`             // AOI ID（可选，与location二选一）
	Location       string `json:"location,omitempty"`       // 中心点坐标（经度,纬度，可选，与id二选一）
	Types          string `json:"types,omitempty"`          // AOI类型（可选，多个类型用|分隔）
	Extensions     string `json:"extensions,omitempty"`     // 返回结果类型（可选，base/all，默认base）
	Output         string `json:"output,omitempty"`         // 输出格式（可选，默认JSON）
	Callback       string `json:"callback,omitempty"`       // 回调函数（可选，用于JSONP跨域）
	Sig            string `json:"sig,omitempty"`            // 签名（可选，需结合安全密钥使用）
	Time           string `json:"time,omitempty"`           // 时间戳（可选，用于签名）
	Language       string `json:"language,omitempty"`       // 语言（可选，默认中文）
}

// ToParams 将AOI边界查询请求参数转换为map[string]string格式
func (req *AOISearchRequest) ToParams() map[string]string {
	params := make(map[string]string)
	if req.ID != "" {
		params["id"] = req.ID // AOI ID（与location二选一）
	}
	if req.Location != "" {
		params["location"] = req.Location // 中心点坐标（与id二选一）
	}
	if req.Types != "" {
		params["types"] = req.Types
	}
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Output != "" {
		params["output"] = req.Output
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Sig != "" {
		params["sig"] = req.Sig
	}
	if req.Time != "" {
		params["time"] = req.Time
	}
	if req.Language != "" {
		params["language"] = req.Language
	}
	return params
}
