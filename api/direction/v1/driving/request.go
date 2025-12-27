package driving

// DrivingRequest 驾车路径规划请求参数
type DrivingRequest struct {
	Origin          string `json:"origin"`                    // 出发点
	Destination     string `json:"destination"`               // 目的地
	OriginID        string `json:"originid,omitempty"`        // 出发点 poiid
	DestinationID   string `json:"destinationid,omitempty"`   // 目的地 poiid
	DestinationType string `json:"destinationtype,omitempty"` // 终点的 poi 类别
	Strategy        string `json:"strategy,omitempty"`        // 驾车选择策略
	Waypoints       string `json:"waypoints,omitempty"`       // 途经点
	AvoidPolygons   string `json:"avoidpolygons,omitempty"`   // 避让区域
	Province        string `json:"province,omitempty"`        // 用汉字填入车牌省份缩写，用于判断是否限行
	Number          string `json:"number,omitempty"`          // 填入除省份及标点之外，车牌的字母和数字（需大写）。用于判断限行相关。
	CarType         string `json:"cartype,omitempty"`         // 车辆类型
	Ferry           string `json:"ferry,omitempty"`           // 在路径规划中，是否使用轮渡
	RoadAggregation string `json:"roadaggregation,omitempty"` // 是否返回路径聚合信息
	NoSteps         string `json:"nosteps,omitempty"`         // 是否返回步骤信息
	Callback        string `json:"callback,omitempty"`        // 回调函数名，用于 JSONP 回调
}

// ToParams 将请求参数转换为map[string]string格式
func (req *DrivingRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["origin"] = req.Origin           // 起点坐标为必填项
	params["destination"] = req.Destination // 终点坐标为必填项
	if req.OriginID != "" {
		params["originid"] = req.OriginID
	}
	if req.DestinationID != "" {
		params["destinationid"] = req.DestinationID
	}
	if req.DestinationType != "" {
		params["destinationtype"] = req.DestinationType
	}
	if req.Strategy != "" {
		params["strategy"] = req.Strategy
	}
	if req.Waypoints != "" {
		params["waypoints"] = req.Waypoints
	}
	if req.AvoidPolygons != "" {
		params["avoidpolygons"] = req.AvoidPolygons
	}
	if req.Province != "" {
		params["province"] = req.Province
	}
	if req.Number != "" {
		params["number"] = req.Number
	}
	if req.CarType != "" {
		params["cartype"] = req.CarType
	}
	if req.Ferry != "" {
		params["ferry"] = req.Ferry
	}
	if req.RoadAggregation != "" {
		params["roadaggregation"] = req.RoadAggregation
	}
	if req.NoSteps != "" {
		params["nosteps"] = req.NoSteps
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	return params
}
