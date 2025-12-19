package driving

import (
	amapType "github.com/enneket/amap/types"
)

// DrivingRequest 驾车路径规划请求参数
type DrivingRequest struct {
	Origin          string                  `json:"origin"`                   // 起点坐标（必填，格式：经度,纬度）
	Destination     string                  `json:"destination"`              // 终点坐标（必填，格式：经度,纬度）
	Strategy        string                  `json:"strategy,omitempty"`       // 驾车策略（可选，默认0=最快捷）
	Waypoints       string                  `json:"waypoints,omitempty"`      // 途经点（可选，格式：lng1,lat1|lng2,lat2）
	DepartureTime   string                  `json:"departure_time,omitempty"` // 出发时间（可选，格式：YYYY-MM-DD HH:mm）
	VehicleType     string                  `json:"vehicle_type,omitempty"`   // 车辆类型（可选，默认0=小型车）
	PlateNumber     string                  `json:"plate_number,omitempty"`   // 车牌号（可选，用于规避限行）
	CoordinateType  amapType.CoordinateType `json:"coordinate_type,omitempty"` // 输入/输出坐标系（可选，默认gcj02）
	Output          amapType.OutputType     `json:"output,omitempty"`          // 输出格式（可选，默认JSON）
	Language        amapType.LanguageType   `json:"language,omitempty"`        // 语言（可选，默认中文）
	Timestamp       string                  `json:"timestamp,omitempty"`       // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *DrivingRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["origin"] = req.Origin     // 起点坐标为必填项
	params["destination"] = req.Destination // 终点坐标为必填项
	if req.Strategy != "" {
		params["strategy"] = req.Strategy
	}
	if req.Waypoints != "" {
		params["waypoints"] = req.Waypoints
	}
	if req.DepartureTime != "" {
		params["departure_time"] = req.DepartureTime
	}
	if req.VehicleType != "" {
		params["vehicle_type"] = req.VehicleType
	}
	if req.PlateNumber != "" {
		params["plate_number"] = req.PlateNumber
	}
	if req.CoordinateType != "" {
		params["coordinate_type"] = string(req.CoordinateType)
	}
	if req.Output != "" {
		params["output"] = string(req.Output)
	}
	if req.Language != "" {
		params["language"] = string(req.Language)
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}
	return params
}
