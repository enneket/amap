package line_keyword

// Suggestion 搜索建议
type Suggestion struct {
	Keywords []string `json:"keywords"` // 关键词建议
	Cities   []string `json:"cities"`   // 城市建议
}

// LineDetail 公交线路详细信息
type LineDetail struct {
	LineID      string `json:"lineid"`      // 线路ID
	Name        string `json:"name"`        // 线路名称
	Type        string `json:"type"`        // 线路类型
	FirstTime   string `json:"start_time"`   // 首班车时间
	LastTime    string `json:"end_time"`     // 末班车时间
	Distance    string `json:"distance"`    // 线路距离（米）
	FromStation string `json:"from_stop"`   // 起点站
	ToStation   string `json:"to_stop"`     // 终点站
}

// LineKeywordResponse 公交路线关键字查询响应结果
type LineKeywordResponse struct {
	Status     string        `json:"status"`     // 返回结果状态值，0表示失败，1表示成功
	Info       string        `json:"info"`       // 返回状态说明
	InfoCode   string        `json:"infocode"`   // 返回状态码
	Count      string        `json:"count"`      // 匹配的公交线路总数
	Suggestion Suggestion    `json:"suggestion"` // 搜索建议
	Lines      []LineDetail  `json:"lines"`      // 公交线路列表
}
