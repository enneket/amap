package distance

import (
	"strconv"

	amapType "github.com/enneket/amap/types"
)

// DistanceRequest 距离测量请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/direction/#distance

type DistanceRequest struct {
	Origins        string                  `json:"origins"`                   // 起点坐标列表（必填），格式："经度1,纬度1|经度2,纬度2|..."，最多支持100个起点
	Destination    string                  `json:"destination"`               // 终点坐标列表（必填），格式："经度1,纬度1|经度2,纬度2|..."，最多支持100个终点
	Type           int                     `json:"type,omitempty"`            // 计算方式（可选，默认1）：1-直线距离，2-驾车导航距离，3-公交规划距离，4-步行规划距离
	CoordinateType amapType.CoordinateType `json:"coordinate_type,omitempty"` // 输入/输出坐标系（可选，默认gcj02，支持wgs84/bd09ll）
	Output         amapType.OutputType     `json:"output,omitempty"`          // 输出格式（可选，默认JSON）
	Language       amapType.LanguageType   `json:"language,omitempty"`        // 语言（可选，默认中文）
	Callback       string                  `json:"callback,omitempty"`        // 回调函数（可选，用于JSONP跨域）
	Timestamp      string                  `json:"timestamp,omitempty"`       // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *DistanceRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["origins"] = req.Origins         // 起点为必填项
	params["destination"] = req.Destination // 终点为必填项

	if req.Type != 0 {
		params["type"] = strconv.Itoa(req.Type)
	}

	if req.CoordinateType != "" {
		params["coordinate_type"] = string(req.CoordinateType)
	}

	if req.Output != "" {
		params["output"] = string(req.Output)
	}

	if req.Language != "" {
		params["language"] = string(req.Language)
	}

	if req.Callback != "" {
		params["callback"] = req.Callback
	}

	if req.Timestamp != "" {
		params["timestamp"] = req.Timestamp
	}

	return params
}
