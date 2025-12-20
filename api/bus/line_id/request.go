package line_id

// LineIDRequest 公交路线ID查询请求参数
type LineIDRequest struct {
	ID   string `json:"id"`   // 公交线路ID（必填）
	City string `json:"city"` // 城市名，如："北京"（必填）
}

// ToParams 将请求参数转换为map
func (req *LineIDRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["id"] = req.ID
	params["city"] = req.City
	return params
}
