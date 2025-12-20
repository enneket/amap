package position

import (
	amapType "github.com/enneket/amap/types"
)

// HardwarePositionRequest 硬件定位请求参数
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/hardware-location
// 支持通过硬件设备信息（如GPS、基站、WiFi等）获取地理位置信息

type HardwarePositionRequest struct {
	Key           string                 `json:"key,omitempty"`                  // 开发者密钥（可选，核心客户端已自动填充，可覆盖）
	DeviceID      string                 `json:"deviceid,omitempty"`            // 设备唯一标识（可选）
	GPS           string                 `json:"gps,omitempty"`                 // GPS原始数据（可选，格式：纬度,经度,速度,方向,时间戳）
	WiFi          string                 `json:"wifi,omitempty"`                // WiFi原始数据（可选，格式：mac1,signal1,ssid1|mac2,signal2,ssid2）
	BaseStation   string                 `json:"basestation,omitempty"`         // 基站原始数据（可选，格式：mcc,mnc,lac,cellid,signal）
	Bluetooth     string                 `json:"bluetooth,omitempty"`           // 蓝牙原始数据（可选，格式：mac1,signal1|mac2,signal2）
	Barometer     string                 `json:"barometer,omitempty"`           // 气压计数据（可选，单位：hPa）
	Accelerometer string                 `json:"accelerometer,omitempty"`       // 加速度计数据（可选，格式：x,y,z）
	Gyroscope     string                 `json:"gyroscope,omitempty"`           // 陀螺仪数据（可选，格式：x,y,z）
	Magnetometer  string                 `json:"magnetometer,omitempty"`        // 磁力计数据（可选，格式：x,y,z）
	Output        amapType.OutputType    `json:"output,omitempty"`              // 输出格式（可选，默认JSON）
	Language      amapType.LanguageType  `json:"language,omitempty"`            // 语言（可选，默认中文）
	Callback      string                 `json:"callback,omitempty"`            // 回调函数（可选，用于JSONP跨域）
	Timestamp     string                 `json:"timestamp,omitempty"`           // 时间戳（可选，核心客户端已自动填充，可覆盖）
}

// ToParams 将请求参数转换为map[string]string格式
func (req *HardwarePositionRequest) ToParams() map[string]string {
	params := make(map[string]string)
	if req.Key != "" {
		params["key"] = req.Key
	}
	if req.DeviceID != "" {
		params["deviceid"] = req.DeviceID
	}
	if req.GPS != "" {
		params["gps"] = req.GPS
	}
	if req.WiFi != "" {
		params["wifi"] = req.WiFi
	}
	if req.BaseStation != "" {
		params["basestation"] = req.BaseStation
	}
	if req.Bluetooth != "" {
		params["bluetooth"] = req.Bluetooth
	}
	if req.Barometer != "" {
		params["barometer"] = req.Barometer
	}
	if req.Accelerometer != "" {
		params["accelerometer"] = req.Accelerometer
	}
	if req.Gyroscope != "" {
		params["gyroscope"] = req.Gyroscope
	}
	if req.Magnetometer != "" {
		params["magnetometer"] = req.Magnetometer
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
