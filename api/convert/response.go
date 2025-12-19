package convert

import (
	amapType "github.com/enneket/amap/types"
)

// ConvertResponse 坐标转换响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/convert
// 返回坐标转换结果，支持批量转换

// ConvertResponse 坐标转换响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/convert
// 返回坐标转换结果，支持批量转换

type ConvertResponse struct {
	amapType.BaseResponse // 继承基础响应（Status/Info/InfoCode）
	Locations string      `json:"locations"` // 转换后的坐标列表，格式："经度,纬度;经度,纬度"
}
