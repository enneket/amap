package bus

// BusRequest 公交路线查询请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api/direction#t5
type BusRequest struct {
	Origin      string `json:"origin"`               // 起点坐标（必填）
	Destination string `json:"destination"`          // 终点坐标（必填）
	City        string `json:"city,omitempty"`       // 城市（可选，默认根据坐标判断）
	CityD       string `json:"cityd,omitempty"`      // 城市（可选，默认根据坐标判断）
	Extensions  string `json:"extensions,omitempty"` // 返回结果类型（可选，base/all，默认base）
	Strategy    string `json:"strategy,omitempty"`   // 策略（可选，默认0）
	NightFlag   string `json:"nightflag,omitempty"`  // 是否计算夜班车（可选，默认0）
	Date        string `json:"date,omitempty"`       // 日期（可选，默认当前日期）
	Time        string `json:"time,omitempty"`       // 时间（可选，默认当前时间）
	Callback    string `json:"callback,omitempty"`   // 回调函数名（可选，JSONP格式）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *BusRequest) ToParams() map[string]string {
	params := make(map[string]string)
	params["origin"] = req.Origin           // 起点坐标为必填项
	params["destination"] = req.Destination // 终点坐标为必填项
	if req.City != "" {
		params["city"] = req.City
	}
	if req.CityD != "" {
		params["cityd"] = req.CityD
	}
	if req.Extensions != "" {
		params["extensions"] = req.Extensions
	}
	if req.Strategy != "" {
		params["strategy"] = req.Strategy
	}
	if req.NightFlag != "" {
		params["nightflag"] = req.NightFlag
	}
	if req.Date != "" {
		params["date"] = req.Date
	}
	if req.Time != "" {
		params["time"] = req.Time
	}
	if req.Callback != "" {
		params["callback"] = req.Callback
	}
	return params
}
