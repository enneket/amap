package position

import (
	amapType "github.com/enneket/amap/types"
)

// HardwarePositionResponse 硬件定位响应
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/hardware-location
// 返回基于硬件设备信息计算的地理位置信息

type HardwarePositionResponse struct {
	amapType.BaseResponse           // 继承基础响应（Status/Info/InfoCode）
	DeviceID              string    `json:"deviceid"`            // 设备唯一标识
	Latitude              float64   `json:"latitude"`            // 纬度（WGS84坐标系）
	Longitude             float64   `json:"longitude"`           // 经度（WGS84坐标系）
	Accuracy              float64   `json:"accuracy"`            // 定位精度（米）
	Speed                 float64   `json:"speed,omitempty"`     // 速度（km/h，可选）
	Direction             float64   `json:"direction,omitempty"` // 方向（度，可选）
	Altitude              float64   `json:"altitude,omitempty"`  // 海拔高度（米，可选）
	Floor                 int       `json:"floor,omitempty"`     // 室内楼层（可选）
	Timestamp             string    `json:"timestamp"`           // 定位时间戳
	LocationType          string    `json:"location_type"`       // 定位类型（gps/wifi/basestation/hybrid）
	Address               string    `json:"address,omitempty"`   // 详细地址（可选）
	POI                   []POIInfo `json:"poi,omitempty"`       // 周边POI信息（可选）
	AdInfo                *AdInfo   `json:"ad_info,omitempty"`   // 行政区划信息（可选）
}

// POIInfo 兴趣点信息
// 包含定位点周边的POI信息

type POIInfo struct {
	Name      string  `json:"name"`      // POI名称
	Distance  int     `json:"distance"`  // 距离定位点的距离（米）
	Latitude  float64 `json:"latitude"`  // POI纬度
	Longitude float64 `json:"longitude"` // POI经度
	Type      string  `json:"type"`      // POI类型
}

// AdInfo 行政区划信息
// 包含定位点所在的行政区划信息

type AdInfo struct {
	Province     string `json:"province"`     // 省份名称
	City         string `json:"city"`         // 城市名称
	District     string `json:"district"`     // 区县名称
	Adcode       string `json:"adcode"`       // 行政区划编码
	CityCode     string `json:"citycode"`     // 城市编码
	ProvinceCode string `json:"provincecode"` // 省份编码
}
