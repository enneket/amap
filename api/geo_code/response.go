package geo_code

import (
	amapType "github.com/enneket/amap/types"
)

// GeoCodeResponse 地理编码响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo#t4
type GeoCodeResponse struct {
	amapType.BaseResponse               // 继承基础响应（Status/Info/InfoCode）
	Count                 string        `json:"count"`    // 匹配的地址数量
	Geocodes              []GeocodeItem `json:"geocodes"` // 地理编码结果列表
}

// GeocodeItem 地理编码结果项
type GeocodeItem struct {
	FormattedAddress string `json:"formatted_address"` // 格式化地址（省+市+区+详细地址）
	Country          string `json:"country"`           // 国家（默认中国）
	Province         string `json:"province"`          // 省份（如"北京市"）
	City             string `json:"city"`              // 城市（如"北京市"）
	Citycode         string `json:"citycode"`          // 城市编码（如"110000"）
	District         string `json:"district"`          // 区县（如"朝阳区"）
	Adcode           string `json:"adcode"`            // 行政区划编码（如"110105"）
	Street           string `json:"street"`            // 街道（如"望京街"）
	Number           string `json:"number"`            // 门牌号（如"8号"）
	Location         string `json:"location"`          // 经纬度（格式："经度,纬度"）
	Level            string `json:"level"`             // 匹配级别（如"门牌号"、"街道"、"区域"）
}
