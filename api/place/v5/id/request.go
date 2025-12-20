package id

// IDRequest POI搜索2.0 ID查询请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/newpoisearch
// 根据POI ID查询详细信息

type IDRequest struct {
	ID         string `json:"id"`                   // POI ID（必填）
	Extensions string `json:"extensions,omitempty"` // 返回结果类型（可选，base/all，默认base）
	Language   string `json:"language,omitempty"`   // 语言（可选，默认中文）
}

// ToParams 将ID查询请求参数转换为map[string]string格式
func (req *IDRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["id"] = req.ID // ID为必填项，直接添加
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Language != "" {
		params["language"] = req.Language
	}
	return params
}
