package line

// LineTrafficResponse 指定线路交通态势查询响应结果
type LineTrafficResponse struct {
	Status      string        `json:"status"`      // 状态码（1：成功；0：失败）
	Info        string        `json:"info"`        // 状态信息
	Infocode    string        `json:"infocode"`    // 状态码说明
	Trafficinfo TrafficInfo   `json:"trafficinfo"` // 交通态势信息
}

// TrafficInfo 交通态势信息
type TrafficInfo struct {
	Description string      `json:"description"` // 路况描述
	Evaluation  Evaluation  `json:"evaluation"`  // 整体路况评估
	Roads       []RoadInfo  `json:"roads"`       // 道路列表
}

// Evaluation 路况评估
type Evaluation struct {
	Expedite    int     `json:"expedite"`    // 畅通路段数
	Congested   int     `json:"congested"`   // 拥堵路段数
	Blocking    int     `json:"blocking"`    // 严重拥堵路段数
	Unknown     int     `json:"unknown"`     // 未知路段数
	Status      string  `json:"status"`      // 整体路况状态
	Description string  `json:"description"` // 整体路况描述
}

// RoadInfo 道路信息
type RoadInfo struct {
	Name        string     `json:"name"`        // 道路名称
	Status      int        `json:"status"`      // 路况状态（0：未知；1：畅通；2：缓行；3：拥堵；4：严重拥堵）
	Direction   string     `json:"direction"`   // 方向
	Lcodes      []string   `json:"lcodes"`      // 道路编码列表
	Polyline    string     `json:"polyline"`    // 道路坐标串
	Speed       float64    `json:"speed"`       // 平均车速
	Jams        []JamInfo  `json:"jams"`        // 拥堵路段列表
}

// JamInfo 拥堵路段信息
type JamInfo struct {
	Polyline    string  `json:"polyline"`    // 拥堵路段坐标串
	Status      int     `json:"status"`      // 拥堵状态（1：缓行；2：拥堵；3：严重拥堵）
	Direction   string  `json:"direction"`   // 方向
	Length      float64 `json:"length"`      // 拥堵长度
	Speed       float64 `json:"speed"`       // 拥堵路段车速
	Time        int     `json:"time"`        // 预计通过时间
	Level       int     `json:"level"`       // 拥堵等级
}
