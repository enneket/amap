package polygon

import (
	"fmt"
)

// PolygonSearchRequest 多边形搜索请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/search
// 基于多边形边界的搜索，用于查询指定多边形区域内的POI
type PolygonSearchRequest struct {
	Keyword    string `json:"keyword,omitempty"`    // 搜索关键词（可选）
	Types      string `json:"types,omitempty"`      // POI类型（可选，多个类型用|分隔）
	City       string `json:"city,omitempty"`       // 搜索城市（可选，默认全国）
	Offset     int    `json:"offset,omitempty"`     // 每页条数（可选，1-50，默认20）
	Page       int    `json:"page,omitempty"`       // 页码（可选，默认1）
	Extensions string `json:"extensions,omitempty"` // 返回结果类型（可选，base/all，默认base）
	Filter     string `json:"filter,omitempty"`     // 过滤条件（可选，如"price:100-200"）
	Language   string `json:"language,omitempty"`   // 语言（可选，默认中文）
	Polygon    string `json:"polygon"`              // 多边形范围（必填，格式："经度1,纬度1;经度2,纬度2;..."）
}

// ToParams 将多边形搜索请求参数转换为map[string]string格式
func (req *PolygonSearchRequest) ToParams() map[string]string {
	params := make(map[string]string)
	
	if req.Keyword != "" {
		params["keyword"] = req.Keyword
	}
	if req.Types != "" {
		params["types"] = req.Types
	}
	if req.City != "" {
		params["city"] = req.City
	}
	if req.Offset != 0 {
		params["offset"] = fmt.Sprintf("%d", req.Offset)
	}
	if req.Page != 0 {
		params["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Filter != "" {
		params["filter"] = req.Filter
	}
	if req.Language != "" {
		params["language"] = req.Language
	}
	params["polygon"] = req.Polygon // 多边形范围为必填项

	return params
}
