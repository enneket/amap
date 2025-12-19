package driving

import (
	amapType "github.com/enneket/amap/types"
)

// DrivingResponse 驾车路径规划响应
type DrivingResponse struct {
	amapType.BaseResponse    // 继承基础响应（Status/Info/InfoCode）
	Route                    Route `json:"route"` // 路线信息
}

// Route 路线信息
type Route struct {
	Origin      string   `json:"origin"`       // 起点坐标
	Destination string   `json:"destination"`  // 终点坐标
	Paths       []Path   `json:"paths"`        // 路径列表
	Distance    string   `json:"distance"`     // 总距离（米）
	Duration    string   `json:"duration"`     // 总时间（秒）
	Tolls       string   `json:"tolls"`        // 总费用（元）
}

// Path 路径信息
type Path struct {
	Distance    string  `json:"distance"`     // 路径距离（米）
	Duration    string  `json:"duration"`     // 路径时间（秒）
	Steps       []Step  `json:"steps"`        // 导航路段
	Polyline    string  `json:"polyline"`     // 路径坐标集合
	Tolls       string  `json:"tolls"`        // 费用（元）
	TollDistance string  `json:"toll_distance"` // 收费路段距离（米）
}

// Step 导航路段信息
type Step struct {
	Instruction string   `json:"instruction"`   // 驾驶指示
	Orientation string   `json:"orientation"`   // 方向
	Road        string   `json:"road"`          // 道路名称
	Distance    string   `json:"distance"`      // 距离（米）
	Duration    string   `json:"duration"`      // 时间（秒）
	Polyline    string   `json:"polyline"`      // 坐标集合
	Action      string   `json:"action"`        // 主要动作
	AssistantAction string `json:"assistant_action"` // 辅助动作
	Tolls       string   `json:"tolls"`        // 费用（元）
	TollRoad    string   `json:"toll_road"`     // 收费道路
}
