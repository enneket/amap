package district

import (
	amapType "github.com/enneket/amap/types"
)

// DistrictResponse 行政区查询响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/district
// 返回匹配的行政区列表，支持嵌套返回子级行政区

type DistrictResponse struct {
	amapType.BaseResponse           // 继承基础响应（Status/Info/InfoCode）
	Count       string            `json:"count"`       // 匹配的行政区数量
	Districts   []DistrictItem    `json:"districts"`   // 行政区列表
	Suggestion  *SuggestionItem   `json:"suggestion,omitempty"` // 建议词列表（可选）
}

// DistrictItem 行政区信息
// 支持嵌套返回子级行政区

type DistrictItem struct {
	Name        string          `json:"name"`        // 行政区名称（如"北京市"）
	Level       string          `json:"level"`       // 行政区级别（country/province/city/district/street/town）
	Adcode      string          `json:"adcode"`      // 行政区划编码
	Citycode    string          `json:"citycode"`    // 城市编码（仅城市和区县级别有值）
	Center      string          `json:"center"`      // 行政区中心点坐标（经度,纬度）
	Polyline    string          `json:"polyline"`    // 行政区边界坐标（仅extensions=all时返回）
	Districts   []DistrictItem  `json:"districts"`   // 子级行政区列表（根据subdistrict参数决定返回级别）
	ParentCity  []string        `json:"parent_city"`  // 父级城市（仅区县级别有值）
	BusinessAreas []string      `json:"business_areas,omitempty"` // 商圈信息（仅城市级别有值）
	Regions     []RegionItem    `json:"regions,omitempty"` // 区域列表（扩展信息）
}

// SuggestionItem 建议词列表
// 当关键字匹配不准确时返回建议词

type SuggestionItem struct {
	Keywords    []string        `json:"keywords"`    // 建议关键字列表
	Cities      []string        `json:"cities"`      // 建议城市列表
}

// RegionItem 区域信息
// 扩展字段，包含区域的详细信息

type RegionItem struct {
	Center      string          `json:"center"`      // 区域中心点坐标
	Name        string          `json:"name"`        // 区域名称
	Level       string          `json:"level"`       // 区域级别
	Polyline    string          `json:"polyline"`    // 区域边界坐标
}
