package search

// IDSearchRequest ID查询请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/search
// 根据POI ID查询详细信息，用于获取指定POI的完整数据

type IDRequest struct {
	ID             string `json:"id"`             // POI ID（必填）
	Extensions     string `json:"extensions,omitempty"` // 返回结果类型（可选，base/all，默认base）
	Output         string `json:"output,omitempty"`  // 输出格式（可选，默认JSON）
	Callback       string `json:"callback,omitempty"` // 回调函数（可选，用于JSONP跨域）
	Sig            string `json:"sig,omitempty"`     // 签名（可选，需结合安全密钥使用）
	Time           string `json:"time,omitempty"`     // 时间戳（可选，用于签名）
	Language       string `json:"language,omitempty"` // 语言（可选，默认中文）
}

// ToParams 将ID查询请求参数转换为map[string]string格式
func (req *IDRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["id"] = req.ID // POI ID为必填项，直接添加
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
