package station_id

// StationIDRequest 公交站ID查询请求参数
type StationIDRequest struct {
	ID   string `json:"id"`   // 公交站点ID（必填）
	City string `json:"city"` // 城市名，如："北京"（必填）
}

// ToParams 将请求参数转换为map
func (req *StationIDRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["id"] = req.ID
	params["city"] = req.City
	return params
}
