package position

import (
	amapType "github.com/enneket/amap/types"
)

// HardwarePositionRequest 硬件定位请求参数（v5版本）
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/hardware-location
// 支持通过硬件设备信息（如GPS、基站、WiFi等）获取地理位置信息，v5版本增强了定位精度和多源数据融合能力

type HardwarePositionRequest struct {
	Key           string                 `json:"key,omitempty"`                  // 开发者密钥（可选，核心客户端已自动填充，可覆盖）
	DeviceID      string                 `json:"deviceid,omitempty"`            // 设备唯一标识（可选）
	GPS           string                 `json:"gps,omitempty"`                 // GPS原始数据（可选，格式：纬度,经度,速度,方向,时间戳,精度）
	WiFi          string                 `json:"wifi,omitempty"`                // WiFi原始数据（可选，格式：mac1,signal1,ssid1,channel1|mac2,signal2,ssid2,channel2）
	BaseStation   string                 `json:"basestation,omitempty"`         // 基站原始数据（可选，格式：mcc,mnc,lac,cellid,signal,ctype）
	Bluetooth     string                 `json:"bluetooth,omitempty"`           // 蓝牙原始数据（可选，格式：mac1,signal1,type1|mac2,signal2,type2）
	Barometer     string                 `json:"barometer,omitempty"`           // 气压计数据（可选，单位：hPa）
	Accelerometer string                 `json:"accelerometer,omitempty"`       // 加速度计数据（可选，格式：x,y,z）
	Gyroscope     string                 `json:"gyroscope,omitempty"`           // 陀螺仪数据（可选，格式：x,y,z）
	Magnetometer  string                 `json:"magnetometer,omitempty"`        // 磁力计数据（可选，格式：x,y,z）
	Orientation   string                 `json:"orientation,omitempty"`         // 设备方向数据（可选，格式：pitch,roll,yaw）
	Pressure      string                 `json:"pressure,omitempty"`            // 压力传感器数据（可选，单位：Pa）
	Light         string                 `json:"light,omitempty"`               // 光线传感器数据（可选，单位：lux）
	Temperature   string                 `json:"temperature,omitempty"`         // 温度传感器数据（可选，单位：℃）
	Humidity      string                 `json:"humidity,omitempty"`            // 湿度传感器数据（可选，单位：%）
	SensorTime    string                 `json:"sensortime,omitempty"`          // 传感器数据采集时间戳（可选）
	PositionMode  string                 `json:"positionmode,omitempty"`        // 定位模式（可选，1：高精度模式，2：低功耗模式，3：仅设备模式）
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
	if req.Orientation != "" {
		params["orientation"] = req.Orientation
	}
	if req.Pressure != "" {
		params["pressure"] = req.Pressure
	}
	if req.Light != "" {
		params["light"] = req.Light
	}
	if req.Temperature != "" {
		params["temperature"] = req.Temperature
	}
	if req.Humidity != "" {
		params["humidity"] = req.Humidity
	}
	if req.SensorTime != "" {
		params["sensortime"] = req.SensorTime
	}
	if req.PositionMode != "" {
		params["positionmode"] = req.PositionMode
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
