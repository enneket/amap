package aoi

import (
	"fmt"
)

// AOISearchRequest AOI边界查询请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/search
// 用于获取指定AOI（兴趣面）的边界坐标
type AOISearchRequest struct {
	ID         string `json:"id,omitempty"`         // AOI ID（可选，与location二选一）
	Location   string `json:"location,omitempty"`   // 中心点坐标（可选，与id二选一）
	Keyword    string `json:"keyword,omitempty"`    // 搜索关键词（可选）
	Offset     int    `json:"offset,omitempty"`     // 每页条数（可选，1-50，默认20）
	Page       int    `json:"page,omitempty"`       // 页码（可选，默认1）
	Extensions string `json:"extensions,omitempty"` // 返回结果类型（可选，base/all，默认base）
	Language   string `json:"language,omitempty"`   // 语言（可选，默认中文）
}

// ToParams 将AOI查询请求参数转换为map[string]string格式
func (req *AOISearchRequest) ToParams() map[string]string {
	params := make(map[string]string)

	if req.ID != "" {
		params["id"] = req.ID
	}
	if req.Location != "" {
		params["location"] = req.Location
	}
	if req.Keyword != "" {
		params["keyword"] = req.Keyword
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
	if req.Language != "" {
		params["language"] = req.Language
	}

	return params
}
