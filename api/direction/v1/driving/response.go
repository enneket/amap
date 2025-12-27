package driving

import (
	amapType "github.com/enneket/amap/types"
)

// DrivingResponse 驾车路径规划响应
type DrivingResponse struct {
	amapType.BaseResponse       // 继承基础响应（Status/Info/InfoCode）
	Route                 Route `json:"route"` // 路线信息
}

// Route 路线信息
type Route struct {
	Origin      string `json:"origin"`      // 起点坐标
	Destination string `json:"destination"` // 终点坐标
	TaxiCost    string `json:"taxi_cost"`   // 打车费用（元）
	Paths       []Path `json:"paths"`       // 驾车换乘方案
}

// Path 路径信息
type Path struct {
	Distance      string `json:"distance"`       // 路径距离（米）
	Duration      string `json:"duration"`       // 预计行驶时间
	Strategy      string `json:"strategy"`       // 导航策略
	Tolls         string `json:"tolls"`          // 费用（元）
	Restrictions  string `json:"restrictions"`   // 限行结果
	TrafficLights string `json:"traffic_lights"` // 红绿灯个数
	TollDistance  string `json:"toll_distance"`  // 收费路段距离（米）
	Steps         []Step `json:"steps"`          // 导航路段
}

// Step 导航路段信息
type Step struct {
	Instruction     string `json:"instruction"`      // 行驶指示
	Orientation     string `json:"orientation"`      // 方向
	Road            string `json:"road"`             // 道路名称
	Distance        string `json:"distance"`         // 此路段距离
	Tolls           string `json:"tolls"`            // 此段收费
	TollDistance    string `json:"toll_distance"`    // 收费路段距离
	TollRoad        string `json:"toll_road"`        // 主要收费道路
	Polyline        string `json:"polyline"`         // 此路段坐标点串
	Action          string `json:"action"`           // 导航主要动作
	AssistantAction string `json:"assistant_action"` // 导航辅助动作
	Tmcs            Tmcs   `json:"tmcs"`             // 驾车导航详细信息

}

type Tmcs struct {
	Distance string `json:"distance"` // 距此段路的长度（米）
	Status   string `json:"status"`   // 此段路的交通情况
	Polyline string `json:"polyline"` // 此段路的轨迹
	Cities   []City `json:"cities"`   // 路线途经行政区划
}

type City struct {
	Name      string     `json:"name"`      // 城市名称
	CityCode  string     `json:"citycode"`  // 城市编码
	AdCode    string     `json:"adcode"`    // 区域编码
	Districts []District `json:"districts"` // 区县信息
}

type District struct {
	Name   string `json:"name"`   // 区县名称
	AdCode string `json:"adcode"` // 区域编码
}
