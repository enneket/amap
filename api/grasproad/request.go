package grasproad

// GraspRoadRequest 轨迹纠偏请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/grasproad
// 用于将原始轨迹点转换为匹配道路的轨迹点

type GraspRoadRequest struct {
	SID             string `json:"sid"`              // 轨迹唯一标识（必填，用于区分不同轨迹）
	Points          string `json:"points"`           // 轨迹点列表（必填，格式："经度,纬度,时间,速度;经度,纬度,时间,速度"）
	CoordTypeInput  string `json:"coord_type_input,omitempty"`  // 输入坐标类型（可选，默认gps，可选值：gps/wgs84/gcj02）
	Extensions      string `json:"extensions,omitempty"`       // 返回结果类型（可选，base/all，默认base）
	Output          string `json:"output,omitempty"`           // 输出格式（可选，默认JSON）
	Callback        string `json:"callback,omitempty"`         // 回调函数（可选，用于JSONP跨域）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *GraspRoadRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["sid"] = req.SID       // 轨迹唯一标识为必填项，直接添加
	params["points"] = req.Points // 轨迹点列表为必填项，直接添加
	if req.CoordTypeInput != "" {
		params["coord_type_input"] = req.CoordTypeInput
	}
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Output != "" {
		params["output"] = req.Output
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	return params
}
