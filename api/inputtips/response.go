package inputtips

import (
	amapType "github.com/enneket/amap/types"
)

// InputtipsResponse 输入提示响应
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/inputtips

type InputtipsResponse struct {
	amapType.BaseResponse // 继承基础响应（Status/Info/InfoCode）
	Count  string       `json:"count"`       // 返回结果数量
	Tips   []TipItem    `json:"tips"`        // 输入提示列表
	Suggestion Suggestion `json:"suggestion,omitempty"` // 建议信息（可选）
}

// TipItem 输入提示结果项
type TipItem struct {
	ID            string   `json:"id"`             // 唯一标识
	Name          string   `json:"name"`           // 名称
	District      string   `json:"district"`       // 区域（如"朝阳区"）
	Adcode        string   `json:"adcode"`         // 行政区划编码（如"110105"）
	Location      string   `json:"location"`       // 经纬度（格式："经度,纬度"）
	Address       string   `json:"address"`        // 地址（如"望京街8号"）
	Type          string   `json:"type"`           // POI类型（如"商务住宅;楼宇;商务写字楼"）
	Typecode      string   `json:"typecode"`       // 类型编码（如"120201"）
	Weight        string   `json:"weight"`         // 权重（如"90"）
	City          string   `json:"city"`           // 城市（如"北京市"）
	Citycode      string   `json:"citycode"`       // 城市编码（如"010"）
	Districtadcode string  `json:"districtadcode"` // 区域编码（如"110105"）
	Province      string   `json:"province"`       // 省份（如"北京市"）
	BusinessArea  string   `json:"business_area"`  // 商圈（如"望京"）
	Children      []string `json:"children,omitempty"` // 子POI列表（可选）
	Groupid       string   `json:"groupid,omitempty"`  // 分组ID（可选）
	PolicyAdcode  string   `json:"policy_adcode,omitempty"` // 政策行政区划编码（可选）
	Poiweight     string   `json:"poiweight,omitempty"` // POI权重（可选）
}

// Suggestion 建议信息
type Suggestion struct {
	Keywords []string `json:"keywords,omitempty"` // 关键词建议（可选）
	Cities   []string `json:"cities,omitempty"`   // 城市建议（可选）
}
