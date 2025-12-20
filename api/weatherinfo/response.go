package weatherinfo

import (
	amapType "github.com/enneket/amap/types"
)

// WeatherinfoResponse 天气信息响应
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/weatherinfo

type WeatherinfoResponse struct {
	amapType.BaseResponse // 继承基础响应（Status/Info/InfoCode）
	Weatherinfo WeatherInfo `json:"weatherinfo"` // 天气信息
	Forecasts   []Forecast   `json:"forecasts,omitempty"` // 预报信息（extensions=all时返回）
	Suggestion  Suggestion   `json:"suggestion,omitempty"` // 生活指数建议（extensions=all时返回）
}

// WeatherInfo 实时天气信息
type WeatherInfo struct {
	City        string `json:"city"`         // 城市名称（如"北京市"）
	CityID      string `json:"cityid"`       // 城市ID（如"101010100"）
	Temp        string `json:"temp"`         // 实时温度（如"22"）
	WD          string `json:"WD"`           // 风向（如"东南风"）
	WS          string `json:"WS"`           // 风力（如"1级"）
	SD          string `json:"SD"`           // 湿度（如"40%"）
	AP          string `json:"AP"`           // 气压（如"1013hPa"）
	NJD         string `json:"njd"`          // 能见度（如"10km"）
	WSE         string `json:"WSE"`          // 风力等级（如"1"）
	Time        string `json:"time"`         // 数据发布时间（如"10:30"）
	IsRadar     string `json:"isRadar"`      // 是否有雷达图（如"1"表示有）
	Radar       string `json:"Radar"`        // 雷达图URL（如"JC_RADAR_AZ9010_JB"）
	Weather     string `json:"weather"`      // 天气状况（如"晴"）
	Temperature string `json:"temperature"`  // 温度范围（如"10~22℃"）
	Winddirection string `json:"winddirection"` // 风向（如"东南"）
	Windpower   string `json:"windpower"`    // 风力（如"1-2级"）
	Humidity    string `json:"humidity"`     // 湿度（如"40%"）
}

// Forecast 预报信息
type Forecast struct {
	City          string `json:"city"`          // 城市名称
	Adcode        string `json:"adcode"`        // 行政区划编码
	Province      string `json:"province"`      // 省份名称
	Reporttime    string `json:"reporttime"`    // 预报发布时间
	Castype       string `json:"castype"`       // 预报类型（如"1"表示24小时预报）
	Forecast      []DailyForecast `json:"forecast"` // 每日预报列表
}

// DailyForecast 每日预报
type DailyForecast struct {
	Date          string `json:"date"`          // 日期（如"2023-05-20"）
	Week          string `json:"week"`          // 星期几（如"六"）
	Dayweather    string `json:"dayweather"`    // 白天天气状况（如"晴"）
	Nightweather  string `json:"nightweather"`  // 夜间天气状况（如"晴"）
	Daytemp       string `json:"daytemp"`       // 白天温度（如"22"）
	Nighttemp     string `json:"nighttemp"`     // 夜间温度（如"10"）
	Daywind       string `json:"daywind"`       // 白天风向（如"东南风"）
	Nightwind     string `json:"nightwind"`     // 夜间风向（如"东南风"）
	Daypower      string `json:"daypower"`      // 白天风力（如"1级"）
	Nightpower    string `json:"nightpower"`    // 夜间风力（如"1级"）
	Daytemp_float float64 `json:"daytemp_float"` // 白天温度（浮点数）
	Nighttemp_float float64 `json:"nighttemp_float"` // 夜间温度（浮点数）
}

// Suggestion 生活指数建议
type Suggestion struct {
	Comf     SuggestionItem `json:"comf"`     // 舒适度指数
	Cw       SuggestionItem `json:"cw"`       // 洗车指数
	Drsg     SuggestionItem `json:"drsg"`     // 穿衣指数
	Flu      SuggestionItem `json:"flu"`      // 感冒指数
	Sport    SuggestionItem `json:"sport"`    // 运动指数
	Trav     SuggestionItem `json:"trav"`     // 旅游指数
	Uv       SuggestionItem `json:"uv"`       // 紫外线指数
	Air      SuggestionItem `json:"air"`      // 空气污染扩散条件指数
	Ac       SuggestionItem `json:"ac"`       // 空调开启指数
	Ag       SuggestionItem `json:"ag"`       // 过敏指数
	Gl       SuggestionItem `json:"gl"`       // 太阳镜指数
	Mu        SuggestionItem `json:"mu"`       // 化妆指数
	Drying    SuggestionItem `json:"drying"`   // 晾晒指数
	Dressing  SuggestionItem `json:"dressing"` // 穿衣指数（详细）
	Fishing   SuggestionItem `json:"fishing"`  // 钓鱼指数
	Spf       SuggestionItem `json:"spf"`      // 防晒指数
}

// SuggestionItem 生活指数建议项
type SuggestionItem struct {
	Brf  string `json:"brf"`  // 简短描述（如"舒适"）
	Txt  string `json:"txt"`  // 详细描述（如"白天温度适宜，风力不大，相信您在这样的天气条件下，应会感到比较清爽和舒适。"）
	Type string `json:"type"` // 指数类型（如"comf"）
}
