package bicycling

import (
	amapType "github.com/enneket/amap/types"
)

// BicyclingResponse 骑行路径规划响应
type BicyclingResponse struct {
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
}

// Path 路径信息
type Path struct {
	Distance    string  `json:"distance"`     // 路径距离（米）
	Duration    string  `json:"duration"`     // 路径时间（秒）
	Steps       []Step  `json:"steps"`        // 骑行步骤
	Polyline    string  `json:"polyline"`     // 路径坐标集合
}

// Step 骑行步骤信息
type Step struct {
	Instruction string `json:"instruction"`   // 骑行指示
	Orientation string `json:"orientation"`   // 方向
	Road        string `json:"road"`          // 道路名称
	Distance    string `json:"distance"`      // 距离（米）
	Duration    string `json:"duration"`      // 时间（秒）
	Polyline    string `json:"polyline"`      // 坐标集合
	Action      string `json:"action"`        // 主要动作
	AssistantAction string `json:"assistant_action"` // 辅助动作
}
