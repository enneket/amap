package station_keyword

// Suggestion 搜索建议
type Suggestion struct {
	Keywords []string `json:"keywords"` // 关键词建议
	Cities   []string `json:"cities"`   // 城市建议
}

// StationDetail 公交站点详细信息
type StationDetail struct {
	ID       string `json:"id"`       // 站点ID
	Name     string `json:"name"`     // 站点名称
	Location string `json:"location"` // 站点坐标
	CityID   string `json:"cityid"`   // 城市ID
	CityName string `json:"cityname"` // 城市名称
	Address  string `json:"address"`  // 站点地址
}

// StationKeywordResponse 公交站关键字查询响应结果
type StationKeywordResponse struct {
	Status     string          `json:"status"`     // 返回结果状态值，0表示失败，1表示成功
	Info       string          `json:"info"`       // 返回状态说明
	InfoCode   string          `json:"infocode"`   // 返回状态码
	Count      string          `json:"count"`      // 匹配的公交站点总数
	Suggestion Suggestion      `json:"suggestion"` // 搜索建议
	Stations   []StationDetail `json:"stations"`   // 公交站点列表
}
