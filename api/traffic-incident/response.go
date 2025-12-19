package traffic_incident

import (
	amapType "github.com/enneket/amap/types"
)

// TrafficIncidentResponse 交通事件查询响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/traffic-incident
// 返回指定区域内的交通事件列表，包含事件详情

type TrafficIncidentResponse struct {
	amapType.BaseResponse           // 继承基础响应（Status/Info/InfoCode）
	Trafficincidents []TrafficIncidentItem `json:"trafficincidents"` // 交通事件列表
}

// TrafficIncidentItem 交通事件详情
// 包含事件的基本信息、位置、类型、级别、影响范围等

type TrafficIncidentItem struct {
	Id               string `json:"id"`               // 事件ID
	Location         string `json:"location"`         // 事件位置（经度,纬度）
	Type             string `json:"type"`             // 事件类型（1-12，对应不同类型的交通事件）
	TypeDes          string `json:"type_des"`         // 事件类型描述
	Level            string `json:"level"`            // 事件级别（1-4，对应轻微、一般、严重、非常严重）
	LevelDes         string `json:"level_des"`        // 事件级别描述
	Description      string `json:"description"`      // 事件描述
	Polyline         string `json:"polyline"`         // 事件影响区域坐标串
	Road             string `json:"road"`             // 事件所在道路
	StartTime        string `json:"start_time"`       // 事件开始时间（时间戳）
	EndTime          string `json:"end_time"`         // 事件结束时间（时间戳）
	Direction        string `json:"direction"`        // 事件影响方向
	Status           string `json:"status"`           // 事件状态（0-结束，1-进行中）
	ImpactLevel      string `json:"impact_level"`      // 影响程度（0-无影响，1-轻微，2-中等，3-严重）
	AffectRoadLength string `json:"affect_road_length"` // 影响道路长度（米）
	Jams             []JamItem `json:"jams"`          // 拥堵信息（仅extensions=all时返回）
	FirstReportTime  string `json:"first_report_time"` // 首次上报时间（时间戳）
	LastReportTime   string `json:"last_report_time"`  // 最新上报时间（时间戳）
}

// JamItem 拥堵信息
// 包含拥堵的详细信息

type JamItem struct {
	Location    string `json:"location"`    // 拥堵位置（经度,纬度）
	Direction   string `json:"direction"`   // 拥堵方向
	Length      string `json:"length"`      // 拥堵长度（米）
	Level       string `json:"level"`       // 拥堵级别（1-4）
	Status      string `json:"status"`      // 拥堵状态
	Speed       string `json:"speed"`       // 拥堵路段平均速度（km/h）
	Polyline    string `json:"polyline"`    // 拥堵路段坐标串
	StartTime   string `json:"start_time"`   // 拥堵开始时间（时间戳）
	EndTime     string `json:"end_time"`     // 拥堵结束时间（时间戳）
}
