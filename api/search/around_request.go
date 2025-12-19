package search

// AroundSearchRequest 周边搜索请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/search
// 基于中心点和半径的搜索，用于查询指定区域内的POI

type AroundSearchRequest struct {
	Keyword         string `json:"keyword"`          // 搜索关键词（必填）
	Location        string `json:"location"`        // 中心点坐标（经度,纬度，必填）
	Radius          string `json:"radius"`          // 搜索半径（单位：米，可选，默认3000）
	Types           string `json:"types,omitempty"`  // POI类型（可选，多个类型用|分隔）
	Sortrule        string `json:"sortrule,omitempty"` // 排序规则（可选，0：综合排序，1：距离排序）
	Offset          string `json:"offset,omitempty"` // 每页条数（可选，1-50，默认20）
	Page            string `json:"page,omitempty"`   // 页码（可选，默认1）
	Extensions      string `json:"extensions,omitempty"` // 返回结果类型（可选，base/all，默认base）
	Output          string `json:"output,omitempty"`  // 输出格式（可选，默认JSON）
	Callback        string `json:"callback,omitempty"` // 回调函数（可选，用于JSONP跨域）
	Sig             string `json:"sig,omitempty"`     // 签名（可选，需结合安全密钥使用）
	Filter          string `json:"filter,omitempty"`  // 过滤条件（可选，如"price:100-200"）
	Origin          string `json:"origin,omitempty"`  // 起点坐标（可选，用于距离排序）
	Time            string `json:"time,omitempty"`     // 时间戳（可选，用于签名）
	Language        string `json:"language,omitempty"` // 语言（可选，默认中文）
}

// ToParams 将周边搜索请求参数转换为map[string]string格式
func (req *AroundSearchRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["keyword"] = req.Keyword // 关键词为必填项，直接添加
	params["location"] = req.Location // 中心点坐标为必填项，直接添加
	if req.Radius != "" {
		params["radius"] = req.Radius
	}
	if req.Types != "" {
		params["types"] = req.Types
	}
	if req.Sortrule != "" {
		params["sortrule"] = req.Sortrule
	}
	if req.Offset != "" {
		params["offset"] = req.Offset
	}
	if req.Page != "" {
		params["page"] = req.Page
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
	if req.Filter != "" {
		params["filter"] = req.Filter
	}
	if req.Origin != "" {
		params["origin"] = req.Origin
	}
	if req.Time != "" {
		params["time"] = req.Time
	}
	if req.Language != "" {
		params["language"] = req.Language
	}
	return params
}
