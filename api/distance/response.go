package distance

import (
	amapType "github.com/enneket/amap/types"
)

// DistanceResponse 距离测量响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/direction/#distance

type DistanceResponse struct {
	amapType.BaseResponse                  // 继承基础响应（Status/Info/InfoCode）
	Results               []DistanceResult `json:"results"` // 距离测量结果列表
}

// DistanceResult 距离测量结果项
type DistanceResult struct {
	OriginId      string `json:"origin_id"` // 起点ID（当传入originid时返回）
	DestinationId string `json:"dest_id"`   // 终点ID（当传入destid时返回）
	Distance      string `json:"distance"`  // 距离（单位：米）
	Duration      string `json:"duration"`  // 预计时间（单位：秒，仅当type为2/3/4时返回）
	Info          string `json:"info"`      // 状态信息（如"OK"表示成功）
	Status        string `json:"status"`    // 状态码（如"1"表示成功）
}
