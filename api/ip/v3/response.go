package ip

import (
	amapType "github.com/enneket/amap/types"
)

// IPConfigResponse IP定位响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/ipconfig
// 返回IP地址对应的地理位置信息

type IPConfigResponse struct {
	amapType.BaseResponse           // 继承基础响应（Status/Info/InfoCode）
	IP         string        `json:"ip"`         // 查询的IP地址
	Province   string        `json:"province"`   // 省份名称
	City       string        `json:"city"`       // 城市名称
	District   string        `json:"district"`   // 区县名称
	Adcode     string        `json:"adcode"`     // 行政区划编码
	Center     string        `json:"center"`     // 城市中心点坐标（经度,纬度）
	ISP        string        `json:"isp"`        // 互联网服务提供商
	Country    string        `json:"country"`    // 国家名称
	Location   *LocationInfo `json:"location,omitempty"` // 详细位置信息（扩展）
}

// LocationInfo 详细位置信息
// 包含更详细的位置描述和坐标信息

type LocationInfo struct {
	Lat        string        `json:"lat"`        // 纬度
	Lon        string        `json:"lon"`        // 经度
	Address    string        `json:"address"`    // 详细地址描述
	CityCode   string        `json:"city_code"`   // 城市编码
	ProvinceCode string      `json:"province_code"` // 省份编码
	DistrictCode string      `json:"district_code"` // 区县编码
	ISPInfo    *ISPInfo      `json:"isp_info,omitempty"` // ISP详细信息
}

// ISPInfo 互联网服务提供商详细信息
// 包含ISP的名称、类型等信息

type ISPInfo struct {
	Name       string        `json:"name"`       // ISP名称
	Type       string        `json:"type"`       // ISP类型（如电信、联通、移动等）
	MCC        string        `json:"mcc"`        // 移动国家码
	MNC        string        `json:"mnc"`        // 移动网络码
}
