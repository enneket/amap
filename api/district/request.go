package district

import (
	amapType "github.com/enneket/amap/types"
)

// DistrictRequest 行政区查询请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/district
// 支持通过关键字搜索行政区，可指定返回子级行政区的级别
// 关键字支持：行政区名称、adcode、citycode等

type DistrictRequest struct {
	Keywords     string                  `json:"keywords"`               // 关键字（必填，如行政区名称、adcode、citycode等）
	Subdistrict  string                  `json:"subdistrict,omitempty"`  // 子级行政区级别（可选，0-3，默认1：返回下一级）
	Page         string                  `json:"page,omitempty"`         // 分页（可选，默认1）
	Offset       string                  `json:"offset,omitempty"`       // 每页条数（可选，默认20）
	Extensions   string                  `json:"extensions,omitempty"`   // 扩展信息（可选，base/all，默认base）
	Filter       string                  `json:"filter,omitempty"`       // 筛选条件（可选，如"citycode:110000"）
	Output       amapType.OutputType     `json:"output,omitempty"`       // 输出格式（可选，默认JSON）
	Callback     string                  `json:"callback,omitempty"`     // 回调函数（可选，用于JSONP跨域）
	Language     amapType.LanguageType   `json:"language,omitempty"`     // 语言（可选，默认中文）
	Timestamp    string                  `json:"timestamp,omitempty"`    // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *DistrictRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["keywords"] = req.Keywords // 关键字为必填项，直接添加
	if req.Subdistrict != "" {
		params["subdistrict"] = req.Subdistrict
	}
	if req.Page != "" {
		params["page"] = req.Page
	}
	if req.Offset != "" {
		params["offset"] = req.Offset
	}
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Filter != "" {
		params["filter"] = req.Filter
	}
	if req.Output != "" {
		params["output"] = string(req.Output)
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Language != "" {
		params["language"] = string(req.Language)
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}
	return params
}
