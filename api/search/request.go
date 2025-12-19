package search

// SearchRequest 高级搜索请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/search
// 支持关键词搜索、周边搜索、多边形搜索等多种搜索方式

type SearchRequest struct {
	Keyword         string `json:"keyword"`          // 搜索关键词（必填）
	Types           string `json:"types,omitempty"`  // POI类型（可选，多个类型用|分隔）
	City            string `json:"city,omitempty"`   // 搜索城市（可选，默认全国）
	Citylimit       string `json:"citylimit,omitempty"` // 仅在指定城市内搜索（可选，0/1，默认0）
	Children        string `json:"children,omitempty"`  // 返回子POI（可选，0/1，默认0）
	Offset          string `json:"offset,omitempty"` // 每页条数（可选，1-50，默认20）
	Page            string `json:"page,omitempty"`   // 页码（可选，默认1）
	Extensions      string `json:"extensions,omitempty"` // 返回结果类型（可选，base/all，默认base）
	Output          string `json:"output,omitempty"`  // 输出格式（可选，默认JSON）
	Callback        string `json:"callback,omitempty"` // 回调函数（可选，用于JSONP跨域）
	Sig             string `json:"sig,omitempty"`     // 签名（可选，需结合安全密钥使用）
	Filter          string `json:"filter,omitempty"`  // 过滤条件（可选，如"price:100-200"）
	Origin          string `json:"origin,omitempty"`  // 起点坐标（可选，用于距离排序）
	Sortrule        string `json:"sortrule,omitempty"` // 排序规则（可选，0：综合排序，1：距离排序）
	Radius          string `json:"radius,omitempty"`  // 搜索半径（可选，单位：米，周边搜索时使用）
	Rectangle       string `json:"rectangle,omitempty"` // 矩形范围（可选，格式："左下角经度,左下角纬度,右上角经度,右上角纬度"）
	Polygon         string `json:"polygon,omitempty"`  // 多边形范围（可选，格式："经度1,纬度1;经度2,纬度2;..."）
	Datatype        string `json:"datatype,omitempty"` // 返回数据类型（可选，如"poi,bus,subway"）
	Province        string `json:"province,omitempty"` // 省份筛选（可选）
	Citycode        string `json:"citycode,omitempty"` // 城市编码筛选（可选）
	Adcode          string `json:"adcode,omitempty"`   // 行政区划编码筛选（可选）
	Industrycode    string `json:"industrycode,omitempty"` // 行业编码筛选（可选）
	Township        string `json:"township,omitempty"` // 乡镇筛选（可选）
	Street          string `json:"street,omitempty"`   // 街道筛选（可选）
	Building        string `json:"building,omitempty"` // 建筑物筛选（可选）
	Floor           string `json:"floor,omitempty"`    // 楼层筛选（可选）
	Location        string `json:"location,omitempty"` // 中心点坐标（可选，周边搜索时使用）
	Style           string `json:"style,omitempty"`    // 返回结果样式（可选，0：标准样式，1：详细样式）
	Time            string `json:"time,omitempty"`     // 时间戳（可选，用于签名）
	Language        string `json:"language,omitempty"` // 语言（可选，默认中文）
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
	if req.Children != "" {
		params["children"] = req.Children
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
	if req.Datatype != "" {
		params["datatype"] = req.Datatype
	}
	if req.Province != "" {
		params["province"] = req.Province
	}
	if req.Citycode != "" {
		params["citycode"] = req.Citycode
	}
	if req.Adcode != "" {
		params["adcode"] = req.Adcode
	}
	if req.Industrycode != "" {
		params["industrycode"] = req.Industrycode
	}
	if req.Township != "" {
		params["township"] = req.Township
	}
	if req.Street != "" {
		params["street"] = req.Street
	}
	if req.Building != "" {
		params["building"] = req.Building
	}
	if req.Floor != "" {
		params["floor"] = req.Floor
	}
	if req.Location != "" {
		params["location"] = req.Location
	}
	if req.Style != "" {
		params["style"] = req.Style
	}
	if req.Time != "" {
		params["time"] = req.Time
	}
	if req.Language != "" {
		params["language"] = req.Language
	}
	return params
}
