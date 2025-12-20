package walking

// WalkingRequest 步行路径规划请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/direction#t4
type WalkingRequest struct {
	Origin        string `json:"origin"`                   // 起点坐标（必填，格式：经度,纬度）
	Destination   string `json:"destination"`              // 终点坐标（必填，格式：经度,纬度）
	OriginID      string `json:"origin_id,omitempty"`      // 起点 POI ID（可选）
	DestinationID string `json:"destination_id,omitempty"` // 目的地 POI ID（可选）
	Callback      string `json:"callback,omitempty"`       // 回调函数名（可选，JSONP格式）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *WalkingRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["origin"] = req.Origin           // 起点坐标为必填项
	params["destination"] = req.Destination // 终点坐标为必填项
	if req.OriginID != "" {
		params["origin_id"] = string(req.OriginID)
	}
	if req.DestinationID != "" {
		params["destination_id"] = string(req.DestinationID)
	}
	if req.Callback != "" {
		params["callback"] = string(req.Callback)
	}
	return params
}
