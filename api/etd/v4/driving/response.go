package driving

import (
	amapType "github.com/enneket/amap/types"
)

// ETDDrivingResponseV4 未来驾车路径规划v4响应
type ETDDrivingResponseV4 struct {
	amapType.BaseResponse    // 继承基础响应（Status/Info/InfoCode）
	Route                    RouteV4 `json:"route"` // 路线信息
}

// RouteV4 路线信息
type RouteV4 struct {
	Origin      string    `json:"origin"`       // 起点坐标
	Destination string    `json:"destination"`  // 终点坐标
	Paths       []PathV4  `json:"paths"`        // 路径列表
	Distance    string    `json:"distance"`     // 总距离（米）
	Duration    string    `json:"duration"`     // 总时间（秒）
	Tolls       string    `json:"tolls"`        // 总费用（元）
	ETD         string    `json:"etd"`          // 预计出发时间
	ETA         string    `json:"eta"`          // 预计到达时间
}

// PathV4 路径信息
type PathV4 struct {
	Distance     string   `json:"distance"`      // 路径距离（米）
	Duration     string   `json:"duration"`      // 路径时间（秒）
	Steps        []StepV4 `json:"steps"`         // 导航路段
	Polyline     string   `json:"polyline"`      // 路径坐标集合
	Tolls        string   `json:"tolls"`         // 费用（元）
	TollDistance string   `json:"toll_distance"` // 收费路段距离（米）
	TrafficLight string   `json:"traffic_light"` // 红绿灯数量
}

// StepV4 导航路段信息
type StepV4 struct {
	Instruction     string   `json:"instruction"`      // 驾驶指示
	Orientation     string   `json:"orientation"`      // 方向
	Road            string   `json:"road"`             // 道路名称
	Distance        string   `json:"distance"`         // 距离（米）
	Duration        string   `json:"duration"`         // 时间（秒）
	Polyline        string   `json:"polyline"`         // 坐标集合
	Action          string   `json:"action"`           // 主要动作
	AssistantAction string   `json:"assistant_action"` // 辅助动作
	Tolls           string   `json:"tolls"`            // 费用（元）
	TollRoad        string   `json:"toll_road"`        // 收费道路
	TrafficLight    string   `json:"traffic_light"`    // 红绿灯数量
}
