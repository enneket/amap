package bus

import (
	amapType "github.com/enneket/amap/types"
)

// BusResponse 公交路线查询响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/direction#t5
type BusResponse struct {
	amapType.BaseResponse       // 继承基础响应（Status/Info/InfoCode）
	Route                 Route `json:"route"` // 路线信息
}

// Route 路线信息
type Route struct {
	Origin      string    `json:"origin"`      // 起点坐标
	Destination string    `json:"destination"` // 终点坐标
	Distance    string    `json:"distance"`    // 总距离（米）
	TaxiCost    string    `json:"taxi_cost"`   // 打车费（元）
	Transits    []Transit `json:"transits"`    // 换乘次数
}

type Transit struct {
	Cost          string      `json:"cost"`           // 换乘费用（元）
	Duration      string      `json:"duration"`       // 换乘时间（秒）
	BusLineName   string      `json:"busline_name"`   // 公交线路名称
	BusLineID     string      `json:"busline_id"`     // 公交线路ID
	DepartureStop BusStopInfo `json:"departure_stop"` // 上车站点
	ArrivalStop   BusStopInfo `json:"arrival_stop"`   // 下车站点
}

// Path 路径信息
type Path struct {
	Distance     string `json:"distance"`      // 路径总距离（米）
	Duration     string `json:"duration"`      // 路径总时间（秒）
	Steps        []Step `json:"steps"`         // 导航路段列表
	Polyline     string `json:"polyline"`      // 路径坐标集合
	Transits     int    `json:"transits"`      // 换乘次数
	Cost         int    `json:"cost"`          // 票价（元）
	WalkDistance string `json:"walk_distance"` // 步行距离（米）
}

// Step 导航路段信息
type Step struct {
	Instruction     string       `json:"instruction"`             // 路段指示
	Orientation     string       `json:"orientation"`             // 方向
	Road            string       `json:"road"`                    // 道路名称
	Distance        string       `json:"distance"`                // 距离（米）
	Duration        string       `json:"duration"`                // 时间（秒）
	Polyline        string       `json:"polyline"`                // 坐标集合
	Type            int          `json:"type"`                    // 路段类型（0:步行, 1:公交）
	Action          string       `json:"action"`                  // 主要动作
	AssistantAction string       `json:"assistant_action"`        // 辅助动作
	BusLineInfo     *BusLineInfo `json:"bus_line_info,omitempty"` // 公交线信息（仅当type=1时返回）
	WalkDetail      *WalkDetail  `json:"walk_detail,omitempty"`   // 步行详细信息（仅当type=0时返回）
}

// BusLineInfo 公交线信息
type BusLineInfo struct {
	BusLineName   string        `json:"busline_name"`   // 公交线路名称
	BusLineID     string        `json:"busline_id"`     // 公交线路ID
	BusLineType   int           `json:"busline_type"`   // 公交类型（0:未知, 1:公交, 2:地铁）
	DepartureStop BusStopInfo   `json:"departure_stop"` // 上车站点
	ArrivalStop   BusStopInfo   `json:"arrival_stop"`   // 下车站点
	PassStopList  []BusStopInfo `json:"pass_stop_list"` // 途经站点列表
	BusNumber     string        `json:"busnumber"`      // 公交车辆数
}

// BusStopInfo 公交站点信息
type BusStopInfo struct {
	Name     string `json:"name"`     // 站点名称
	Location string `json:"location"` // 站点坐标
}

// WalkDetail 步行详细信息
type WalkDetail struct {
	Instruction string `json:"instruction"` // 步行指示
	Orientation string `json:"orientation"` // 方向
	Road        string `json:"road"`        // 道路名称
	Distance    string `json:"distance"`    // 距离（米）
	Duration    string `json:"duration"`    // 时间（秒）
	Polyline    string `json:"polyline"`    // 坐标集合
}
