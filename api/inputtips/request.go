package inputtips

import (
	amapType "github.com/enneket/amap/types"
)

// InputtipsRequest 输入提示请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/inputtips

type InputtipsRequest struct {
	Keywords   string                `json:"keywords"`             // 关键字（必填）
	City       string                `json:"city,omitempty"`       // 城市（可选）
	Type       string                `json:"type,omitempty"`       // POI类型（可选）
	Location   string                `json:"location,omitempty"`   // 经纬度（可选）
	Radius     string                `json:"radius,omitempty"`     // 搜索半径（可选）
	Offset     string                `json:"offset,omitempty"`     // 返回结果数量（可选，默认10）
	Page       string                `json:"page,omitempty"`       // 当前页码（可选，默认1）
	Extensions string                `json:"extensions,omitempty"` // 返回结果扩展（可选，默认base）
	Callback   string                `json:"callback,omitempty"`   // 回调函数（可选，用于JSONP跨域）
	Datatype   string                `json:"datatype,omitempty"`   // 返回数据类型（可选，默认all）
	Citylimit  string                `json:"citylimit,omitempty"`  // 是否限制城市（可选，默认false）
	Children   string                `json:"children,omitempty"`   // 是否返回子POI（可选，默认false）
	Language   amapType.LanguageType `json:"language,omitempty"`   // 语言（可选，默认中文）
	Timestamp  string                `json:"timestamp,omitempty"`  // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *InputtipsRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["keywords"] = req.Keywords // 关键字为必填项，直接添加

	if req.City != "" {
		params["city"] = req.City
	}
	if req.Type != "" {
		params["type"] = req.Type
	}
	if req.Location != "" {
		params["location"] = req.Location
	}
	if req.Radius != "" {
		params["radius"] = req.Radius
	}
	if req.Offset != "" {
		params["offset"] = req.Offset
	}
	if req.Page != "" {
		params["page"] = req.Page
	}
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	if req.Datatype != "" {
		params["datatype"] = req.Datatype
	}
	if req.Citylimit != "" {
		params["citylimit"] = req.Citylimit
	}
	if req.Children != "" {
		params["children"] = req.Children
	}
	if req.Language != "" {
		params["language"] = string(req.Language)
	}
	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}

	return params
}
