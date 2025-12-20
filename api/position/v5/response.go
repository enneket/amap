package position

import (
	amapType "github.com/enneket/amap/types"
)

// HardwarePositionResponse 硬件定位响应（v5版本）
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/hardware-location
// 返回基于硬件设备信息计算的地理位置信息，v5版本增强了定位精度和多源数据融合能力

type HardwarePositionResponse struct {
	amapType.BaseResponse               // 继承基础响应（Status/Info/InfoCode）
	DeviceID              string        `json:"deviceid"`              // 设备唯一标识
	Latitude              float64       `json:"latitude"`              // 纬度（WGS84坐标系）
	Longitude             float64       `json:"longitude"`             // 经度（WGS84坐标系）
	Accuracy              float64       `json:"accuracy"`              // 定位精度（米）
	Speed                 float64       `json:"speed,omitempty"`       // 速度（km/h，可选）
	Direction             float64       `json:"direction,omitempty"`   // 方向（度，可选）
	Altitude              float64       `json:"altitude,omitempty"`    // 海拔高度（米，可选）
	Floor                 int           `json:"floor,omitempty"`       // 室内楼层（可选）
	Timestamp             string        `json:"timestamp"`             // 定位时间戳
	LocationType          string        `json:"location_type"`         // 定位类型（gps/wifi/basestation/hybrid）
	Address               string        `json:"address,omitempty"`     // 详细地址（可选）
	POI                   []POIInfo     `json:"poi,omitempty"`         // 周边POI信息（可选）
	AdInfo                *AdInfo       `json:"ad_info,omitempty"`     // 行政区划信息（可选）
	SensorInfo            *SensorInfo   `json:"sensor_info,omitempty"` // 传感器数据信息（可选，v5新增）
	Confidence            float64       `json:"confidence,omitempty"`  // 定位置信度（0-100，可选，v5新增）
	Indoor                bool          `json:"indoor,omitempty"`      // 是否为室内定位（可选，v5新增）
	MapMatch              *MapMatchInfo `json:"map_match,omitempty"`   // 地图匹配信息（可选，v5新增）
	TraceID               string        `json:"trace_id"`              // 定位请求唯一标识（v5新增）
}

// POIInfo 兴趣点信息
// 包含定位点周边的POI信息，v5版本增强了POI数据的丰富度

type POIInfo struct {
	Name      string  `json:"name"`              // POI名称
	Distance  int     `json:"distance"`          // 距离定位点的距离（米）
	Latitude  float64 `json:"latitude"`          // POI纬度
	Longitude float64 `json:"longitude"`         // POI经度
	Type      string  `json:"type"`              // POI类型
	Address   string  `json:"address,omitempty"` // POI详细地址（可选，v5新增）
	Phone     string  `json:"phone,omitempty"`   // POI电话（可选，v5新增）
}

// AdInfo 行政区划信息
// 包含定位点所在的行政区划信息

type AdInfo struct {
	Province     string `json:"province"`           // 省份名称
	City         string `json:"city"`               // 城市名称
	District     string `json:"district"`           // 区县名称
	Adcode       string `json:"adcode"`             // 行政区划编码
	CityCode     string `json:"citycode"`           // 城市编码
	ProvinceCode string `json:"provincecode"`       // 省份编码
	Township     string `json:"township,omitempty"` // 乡镇名称（可选，v5新增）
	Village      string `json:"village,omitempty"`  // 村庄名称（可选，v5新增）
}

// SensorInfo 传感器数据信息
// 包含定位计算使用的传感器数据信息（v5新增）

type SensorInfo struct {
	GPSValid         bool     `json:"gps_valid,omitempty"`         // GPS数据是否有效（可选）
	WiFiValid        bool     `json:"wifi_valid,omitempty"`        // WiFi数据是否有效（可选）
	BaseStationValid bool     `json:"basestation_valid,omitempty"` // 基站数据是否有效（可选）
	BluetoothValid   bool     `json:"bluetooth_valid,omitempty"`   // 蓝牙数据是否有效（可选）
	UsedSensorTypes  []string `json:"used_sensor_types,omitempty"` // 参与定位计算的传感器类型列表（可选）
}

// MapMatchInfo 地图匹配信息
// 包含定位点与地图道路的匹配信息（v5新增）

type MapMatchInfo struct {
	RoadName  string  `json:"road_name,omitempty"` // 匹配的道路名称（可选）
	RoadType  string  `json:"road_type,omitempty"` // 道路类型（可选）
	Lane      string  `json:"lane,omitempty"`      // 车道信息（可选）
	Offset    float64 `json:"offset,omitempty"`    // 偏离道路中心线的距离（米，可选）
	Direction float64 `json:"direction,omitempty"` // 道路方向（度，可选）
}
