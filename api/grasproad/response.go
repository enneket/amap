package grasproad

import (
	amapType "github.com/enneket/amap/types"
)

// GraspRoadResponse 轨迹纠偏响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/grasproad
// 返回轨迹纠偏后的结果，包括道路匹配信息、纠偏后的轨迹点等

type GraspRoadResponse struct {
	amapType.BaseResponse // 继承基础响应（Status/Info/InfoCode）
	SID      string     `json:"sid"`      // 轨迹唯一标识
	Paths    []PathItem `json:"paths"`    // 纠偏后的轨迹路径列表
}

// PathItem 纠偏后的轨迹路径
// 包含轨迹的基本信息和轨迹点列表

type PathItem struct {
	Points   []PointItem `json:"points"`   // 纠偏后的轨迹点列表
	Distance int         `json:"distance,omitempty"` // 轨迹总距离（单位：米，仅extensions=all时返回）
	Time     int         `json:"time,omitempty"`     // 轨迹总时间（单位：秒，仅extensions=all时返回）
	Steps    []StepItem  `json:"steps,omitempty"`    // 轨迹分段信息（仅extensions=all时返回）
}

// PointItem 纠偏后的轨迹点
// 包含经纬度、时间、速度、方向等信息

type PointItem struct {
	Location string  `json:"location"` // 纠偏后的坐标（经度,纬度）
	Time     int64   `json:"time"`     // 时间戳（秒）
	Speed    float64 `json:"speed"`    // 速度（单位：km/h）
	Direction int    `json:"direction,omitempty"` // 方向（单位：度，仅extensions=all时返回）
	RoadID   string  `json:"road_id,omitempty"`   // 道路ID（仅extensions=all时返回）
	RoadName string  `json:"road_name,omitempty"` // 道路名称（仅extensions=all时返回）
	POIID    string  `json:"poi_id,omitempty"`    // POI ID（仅extensions=all时返回）
	POIName  string  `json:"poi_name,omitempty"`  // POI名称（仅extensions=all时返回）
	MatchType int    `json:"match_type,omitempty"` // 匹配类型（0：不匹配，1：匹配，仅extensions=all时返回）
	Status   int     `json:"status,omitempty"`   // 轨迹点状态（0：正常，1：异常，仅extensions=all时返回）
}

// StepItem 轨迹分段信息
// 包含道路信息和分段轨迹点列表

type StepItem struct {
	StartIndex int         `json:"start_index"` // 起始点索引
	EndIndex   int         `json:"end_index"`   // 结束点索引
	Road       RoadItem    `json:"road,omitempty"` // 道路信息（仅extensions=all时返回）
	Points     []PointItem `json:"points,omitempty"` // 分段轨迹点列表（仅extensions=all时返回）
}

// RoadItem 道路信息
// 包含道路名称、道路ID、道路类型等

type RoadItem struct {
	ID      string `json:"id,omitempty"`      // 道路ID
	Name    string `json:"name,omitempty"`    // 道路名称
	Type    int    `json:"type,omitempty"`    // 道路类型（0：未知，1：高速公路，2：城市快速路，3：国道，4：省道，5：县道，6：乡道，7：村道，8：其他）
	Level   int    `json:"level,omitempty"`   // 道路等级（0：未知，1：高速，2：快速，3：主干道，4：次干道，5：支路）
	Width   float64 `json:"width,omitempty"`  // 道路宽度（单位：米）
	Lanes   int    `json:"lanes,omitempty"`   // 车道数
	MaxSpeed int   `json:"max_speed,omitempty"` // 最高限速（单位：km/h）
}
