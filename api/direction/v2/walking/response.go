package walking

import (
	amapType "github.com/enneket/amap/types"
)

// WalkingResponseV2 步行路线规划v2响应
type WalkingResponseV2 struct {
	amapType.BaseResponse    // 继承基础响应（Status/Info/InfoCode）
	Route                    RouteV2 `json:"route"` // 路线信息
}

// RouteV2 路线信息
type RouteV2 struct {
	Origin      string    `json:"origin"`       // 起点坐标
	Destination string    `json:"destination"`  // 终点坐标
	Paths       []PathV2  `json:"paths"`        // 路径列表
	Distance    string    `json:"distance"`     // 总距离（米）
	Duration    string    `json:"duration"`     // 总时间（秒）
}

// PathV2 路径信息
type PathV2 struct {
	Distance    string    `json:"distance"`     // 路径距离（米）
	Duration    string    `json:"duration"`     // 路径时间（秒）
	Steps       []StepV2  `json:"steps"`        // 导航路段
	Polyline    string    `json:"polyline"`     // 路径坐标集合
}

// StepV2 导航路段信息
type StepV2 struct {
	Instruction string   `json:"instruction"`   // 步行指示
	Orientation string   `json:"orientation"`   // 方向
	Road        string   `json:"road"`          // 道路名称
	Distance    string   `json:"distance"`      // 距离（米）
	Duration    string   `json:"duration"`      // 时间（秒）
	Polyline    string   `json:"polyline"`      // 坐标集合
	Action      string   `json:"action"`        // 主要动作
	AssistantAction string `json:"assistant_action"` // 辅助动作
	WalkType    string   `json:"walk_type"`     // 步行类型
}
