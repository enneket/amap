package line_id

// StationInfo 站点信息
type StationInfo struct {
	ID       string `json:"id"`       // 站点ID
	Name     string `json:"name"`     // 站点名称
	Location string `json:"location"` // 站点坐标
}

// LineIDResponse 公交路线ID查询响应结果
type LineIDResponse struct {
	Status    string        `json:"status"`    // 返回结果状态值，0表示失败，1表示成功
	Info      string        `json:"info"`      // 返回状态说明
	InfoCode  string        `json:"infocode"`  // 返回状态码
	LineID    string        `json:"lineid"`    // 线路ID
	Name      string        `json:"name"`      // 线路名称
	Type      string        `json:"type"`      // 线路类型
	FirstTime string        `json:"start_time"` // 首班车时间
	LastTime  string        `json:"end_time"` // 末班车时间
	Distance  string        `json:"distance"`  // 线路距离（米）
	Polyline  string        `json:"polyline"`  // 线路坐标集合
	Stations  []StationInfo `json:"stations"`  // 站点列表
}
