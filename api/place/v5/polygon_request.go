package v5

// PolygonSearchRequest POI搜索2.0多边形搜索请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/newpoisearch
// 基于多边形边界的搜索，用于查询指定多边形区域内的POI

type PolygonSearchRequest struct {
	Keyword     string `json:"keyword"`     // 搜索关键词（必填）
	Polygon     string `json:"polygon"`     // 多边形范围（格式："经度1,纬度1;经度2,纬度2;..."，必填）
	Types       string `json:"types,omitempty"`   // POI类型（可选，多个类型用|分隔）
	Sortrule    string `json:"sortrule,omitempty"`   // 排序规则（可选，0：综合排序，1：距离排序）
	Offset      string `json:"offset,omitempty"`    // 每页条数（可选，1-50，默认20）
	Page        string `json:"page,omitempty"`      // 页码（可选，默认1）
	Extensions  string `json:"extensions,omitempty"`  // 返回结果类型（可选，base/all，默认base）
	Filter      string `json:"filter,omitempty"`     // 过滤条件（可选，如"price:100-200"）
	Origin      string `json:"origin,omitempty"`     // 起点坐标（可选，用于距离排序）
	Language    string `json:"language,omitempty"`   // 语言（可选，默认中文）
}

// ToParams 将多边形搜索请求参数转换为map[string]string格式
func (req *PolygonSearchRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["keyword"] = req.Keyword // 关键词为必填项，直接添加
	params["polygon"] = req.Polygon // 多边形范围为必填项，直接添加
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
	if req.Filter != "" {
		params["filter"] = req.Filter
	}
	if req.Origin != "" {
		params["origin"] = req.Origin
	}
	if req.Language != "" {
		params["language"] = req.Language
	}
	return params
}
