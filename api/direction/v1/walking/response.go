package walking

import (
	amapType "github.com/enneket/amap/types"
)

// WalkingResponse 步行路径规划响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/direction#t4
type WalkingResponse struct {
	amapType.BaseResponse       // 继承基础响应（Status/Info/InfoCode）
	Route                 Route `json:"route"` // 路线信息
}

// Route 路线信息
type Route struct {
	Origin      string `json:"origin"`      // 起点坐标
	Destination string `json:"destination"` // 终点坐标
	Paths       []Path `json:"paths"`       // 路径列表
	Distance    string `json:"distance"`    // 总距离（米）
	Duration    string `json:"duration"`    // 总时间（秒）
	Tolls       string `json:"tolls"`       // 总费用（元，步行无费用）
}

// Path 路径信息
type Path struct {
	Distance string `json:"distance"` // 路径距离（米）
	Duration string `json:"duration"` // 路径时间（秒）
	Steps    []Step `json:"steps"`    // 步行步骤
}

// Step 步行步骤信息
type Step struct {
	Instruction     string `json:"instruction"`      // 步行指示
	Road            any    `json:"road"`             // 道路名称
	Distance        string `json:"distance"`         // 距离（米）
	Orientation     any    `json:"orientation"`      // 方向
	Duration        string `json:"duration"`         // 时间（秒）
	Polyline        string `json:"polyline"`         // 坐标集合
	Action          any    `json:"action"`           // 主要动作
	AssistantAction any    `json:"assistant_action"` // 辅助动作
	WalkType        string `json:"walk_type"`        // 这段路是否存在特殊的方式
}
