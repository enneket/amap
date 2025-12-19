package electric

import (
	amapType "github.com/enneket/amap/types"
)

// ElectricResponseV2 电动车路线规划v2响应
type ElectricResponseV2 struct {
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
	Tolls       string    `json:"tolls"`        // 总费用（元）
}

// PathV2 路径信息
type PathV2 struct {
	Distance     string       `json:"distance"`      // 路径距离（米）
	Duration     string       `json:"duration"`      // 路径时间（秒）
	Steps        []StepV2     `json:"steps"`         // 导航路段
	Polyline     string       `json:"polyline"`      // 路径坐标集合
	Tolls        string       `json:"tolls"`         // 费用（元）
	TollDistance string       `json:"toll_distance"` // 收费路段距离（米）
	ChargeInfo   ChargeInfoV2 `json:"charge_info"`   // 充电信息
}

// ChargeInfoV2 充电信息
type ChargeInfoV2 struct {
	BatteryUsage   string        `json:"battery_usage"`    // 电量消耗（kWh）
	ChargeStations []ChargeStation `json:"charge_stations"` // 充电站点列表
	TotalChargeFee string        `json:"total_charge_fee"` // 总充电费用（元）
}

// ChargeStation 充电站点信息
type ChargeStation struct {
	Id           string  `json:"id"`            // 站点ID
	Name         string  `json:"name"`          // 站点名称
	Location     string  `json:"location"`      // 站点坐标
	Distance     string  `json:"distance"`      // 距离（米）
	ChargeFee    string  `json:"charge_fee"`    // 充电费用（元）
	ChargeTime   string  `json:"charge_time"`   // 充电时间（分钟）
	BatteryAfter string  `json:"battery_after"` // 充电后电量（%）
}

// StepV2 导航路段信息
type StepV2 struct {
	Instruction     string  `json:"instruction"`      // 驾驶指示
	Orientation     string  `json:"orientation"`      // 方向
	Road            string  `json:"road"`             // 道路名称
	Distance        string  `json:"distance"`         // 距离（米）
	Duration        string  `json:"duration"`         // 时间（秒）
	Polyline        string  `json:"polyline"`         // 坐标集合
	Action          string  `json:"action"`           // 主要动作
	AssistantAction string  `json:"assistant_action"` // 辅助动作
	ChargeStation   *ChargeStation `json:"charge_station,omitempty"` // 充电站点（如果有）
}
