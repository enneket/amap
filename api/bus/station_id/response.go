package station_id

// StationLine 公交线路信息
type StationLine struct {
	LineID    string        `json:"lineid"` // 线路ID
	LineName  string        `json:"name"`   // 线路名称
	FirstTime string        `json:"start_time"` // 首班车时间
	LastTime  string        `json:"end_time"` // 末班车时间
	Distance  string        `json:"distance"` // 线路距离（米）
	Stations  []StationInfo `json:"stations"` // 站点列表
}

// StationInfo 站点信息
type StationInfo struct {
	ID       string `json:"id"`       // 站点ID
	Name     string `json:"name"`     // 站点名称
	Location string `json:"location"` // 站点坐标
}

// StationIDResponse 公交站ID查询响应结果
type StationIDResponse struct {
	Status    string        `json:"status"`    // 返回结果状态值，0表示失败，1表示成功
	Info      string        `json:"info"`      // 返回状态说明
	InfoCode  string        `json:"infocode"`  // 返回状态码
	StationID string        `json:"stationid"` // 公交站点ID
	Name      string        `json:"name"`      // 公交站点名称
	Location  string        `json:"location"`  // 公交站点坐标
	Lines     []StationLine `json:"lines"`     // 经过该站点的公交线路列表
}
