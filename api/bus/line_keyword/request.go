package line_keyword

// LineKeywordRequest 公交路线关键字查询请求参数
type LineKeywordRequest struct {
	Keywords string `json:"keywords"` // 公交线路名称关键字（必填）
	City     string `json:"city"`     // 城市名，如："北京"（必填）
	Offset   string `json:"offset"`   // 每页记录数，默认20，最大50
	Page     string `json:"page"`     // 当前页码，默认1
}

// ToParams 将请求参数转换为map
func (req *LineKeywordRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["keywords"] = req.Keywords
	params["city"] = req.City
	params["offset"] = req.Offset
	params["page"] = req.Page
	return params
}
