package v5

// SearchRequest POI搜索2.0通用请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/newpoisearch
// 支持关键词搜索、周边搜索、多边形搜索等多种搜索方式

type SearchRequest struct {
	Keyword     string `json:"keyword"`     // 搜索关键词（必填）
	Types       string `json:"types,omitempty"`   // POI类型（可选，多个类型用|分隔）
	City        string `json:"city,omitempty"`    // 搜索城市（可选，默认全国）
	Citylimit   string `json:"citylimit,omitempty"`  // 仅在指定城市内搜索（可选，0/1，默认0）
	Offset      string `json:"offset,omitempty"`    // 每页条数（可选，1-50，默认20）
	Page        string `json:"page,omitempty"`      // 页码（可选，默认1）
	Extensions  string `json:"extensions,omitempty"`  // 返回结果类型（可选，base/all，默认base）
	Filter      string `json:"filter,omitempty"`     // 过滤条件（可选，如"price:100-200"）
	Origin      string `json:"origin,omitempty"`     // 起点坐标（可选，用于距离排序）
	Sortrule    string `json:"sortrule,omitempty"`   // 排序规则（可选，0：综合排序，1：距离排序）
	Radius      string `json:"radius,omitempty"`     // 搜索半径（可选，单位：米，周边搜索时使用）
	Rectangle   string `json:"rectangle,omitempty"`  // 矩形范围（可选，格式："左下角经度,左下角纬度,右上角经度,右上角纬度"）
	Polygon     string `json:"polygon,omitempty"`    // 多边形范围（可选，格式："经度1,纬度1;经度2,纬度2;..."）
	Adcode      string `json:"adcode,omitempty"`     // 行政区划编码筛选（可选）
	Building    string `json:"building,omitempty"`   // 建筑物筛选（可选）
	Location    string `json:"location,omitempty"`   // 中心点坐标（可选，周边搜索时使用）
	Language    string `json:"language,omitempty"`   // 语言（可选，默认中文）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *SearchRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["keyword"] = req.Keyword // 关键词为必填项，直接添加
	if req.Types != "" {
		params["types"] = req.Types
	}
	if req.City != "" {
		params["city"] = req.City
	}
	if req.Citylimit != "" {
		params["citylimit"] = req.Citylimit
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
	if req.Sortrule != "" {
		params["sortrule"] = req.Sortrule
	}
	if req.Radius != "" {
		params["radius"] = req.Radius
	}
	if req.Rectangle != "" {
		params["rectangle"] = req.Rectangle
	}
	if req.Polygon != "" {
		params["polygon"] = req.Polygon
	}
	if req.Adcode != "" {
		params["adcode"] = req.Adcode
	}
	if req.Building != "" {
		params["building"] = req.Building
	}
	if req.Location != "" {
		params["location"] = req.Location
	}
	if req.Language != "" {
		params["language"] = req.Language
	}
	return params
}
