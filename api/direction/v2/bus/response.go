package bus

import (
	amapType "github.com/enneket/amap/types"
)

// BusResponseV2 公交路线规划v2响应
type BusResponseV2 struct {
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
	Transits    []Transit `json:"transits"`     // 公交换乘方案
}

// Transit 公交换乘方案
type Transit struct {
	Distance     string      `json:"distance"`      // 换乘距离（米）
	Duration     string      `json:"duration"`      // 换乘时间（秒）
	WalkingDistance string    `json:"walking_distance"` // 步行距离（米）
	BusLines     []BusLine   `json:"buslines"`      // 公交路线
	Steps        []StepV2    `json:"steps"`         // 换乘步骤
}

// BusLine 公交路线信息
type BusLine struct {
	Name         string      `json:"name"`          // 公交路线名称
	BusLineType  string      `json:"busline_type"`  // 公交类型
	DepartureBusStation BusStation `json:"departure_busstation"` // 出发站点
	ArrivalBusStation   BusStation `json:"arrival_busstation"`   // 到达站点
	ViaBusStations []BusStation `json:"via_busstations"` // 途经站点
}

// BusStation 公交站点信息
type BusStation struct {
	Id           string      `json:"id"`            // 站点ID
	Name         string      `json:"name"`          // 站点名称
	Location     string      `json:"location"`      // 站点坐标
}

// StepV2 导航路段信息
type StepV2 struct {
	Instruction string      `json:"instruction"`   // 指示
	Orientation string      `json:"orientation"`   // 方向
	Road        string      `json:"road"`          // 道路名称
	Distance    string      `json:"distance"`      // 距离（米）
	Duration    string      `json:"duration"`      // 时间（秒）
	Polyline    string      `json:"polyline"`      // 坐标集合
	Action      string      `json:"action"`        // 主要动作
	AssistantAction string  `json:"assistant_action"` // 辅助动作
	WalkType    string      `json:"walk_type"`     // 步行类型
	BusLine     *BusLine    `json:"busline,omitempty"` // 公交路线信息（公交路段）
}
