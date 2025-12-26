# 高德地图 API Go 客户端

> **AI 生成提醒**: 本项目目前绝大部分由 AI 辅助生成的代码，包括但不限于：
> - API 服务的实现代码
> - 单元测试用例
> - 部分示例代码
> - 文档内容
> 
> **使用建议**:
> - 在生产环境使用前，请充分测试相关功能
> - 建议定期审查和更新 AI 生成的代码
> - 如发现问题，欢迎提交 Issue 或 Pull Request

这是一个高德地图 API 的 Go 语言客户端库，提供了简洁易用的接口来访问高德地图的各种服务。

## 功能特性

- 支持多种高德地图 API 服务
- 简洁易用的接口设计
- 完整的错误处理
- 支持超时配置
- 支持代理配置
- 支持 API 签名

## 安装

```bash
go get github.com/enneket/amap
```

## 快速开始

### 初始化客户端

```go
import "github.com/enneket/amap"

func main() {
    // 创建配置
    config := &amap.Config{
        Key:    "your_amap_api_key",
        SecurityKey: "your_amap_security_key", // 可选，用于生成签名
        Timeout: 30 * time.Second,
    }
    
    // 初始化客户端
    client, err := amap.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用客户端调用 API
    // ...
}
```

## 支持的 API

### 地理编码与逆地理编码
- `GeoCode`: 地理编码
- `ReGeocode`: 逆地理编码

### 路径规划
- `Walking`: 步行路径规划 (v1)
- `Driving`: 驾车路径规划 (v1)
- `Bicycling`: 骑行路径规划 (v1)
- `WalkingV2`: 步行路径规划 (v2)
- `DrivingV2`: 驾车路径规划 (v2)
- `BicyclingV2`: 骑行路径规划 (v2)
- `BusV2`: 公交路线规划 (v2)
- `ElectricV2`: 电动车路线规划 (v2)
- `ETDDrivingV4`: 未来驾车路径规划 (v4)

### 距离测量
- `Distance`: 距离测量

### 行政区划查询
- `District`: 行政区查询

### 交通态势
- `TrafficIncident`: 交通事件查询
- `LineTrafficStatus`: 指定线路交通态势查询
- `CircleTrafficStatus`: 圆形区域内交通态势查询
- `RectangleTrafficStatus`: 矩形区域内交通态势查询

### IP 定位
- `IPConfig`: IP 定位 (v3)
- `IPV5Config`: IP 定位 (v5)

### 坐标转换
- `Convert`: 坐标转换

### 轨迹纠偏
- `GraspRoad`: 轨迹纠偏

### POI 搜索
- `PlaceV3ID`: POI ID 查询 (v3)
- `PlaceV3Text`: POI 文本搜索 (v3)
- `PlaceV3Around`: POI 周边搜索 (v3)
- `PlaceV3Polygon`: POI 多边形搜索 (v3)
- `PlaceV3AOI`: POI AOI 查询 (v3)
- `PlaceV5ID`: POI ID 查询 (v5)
- `PlaceV5Text`: POI 文本搜索 (v5)
- `PlaceV5Around`: POI 周边搜索 (v5)
- `PlaceV5Polygon`: POI 多边形搜索 (v5)
- `PlaceV5AOI`: POI AOI 查询 (v5)

### 输入提示
- `Inputtips`: 输入提示

### 天气信息
- `Weatherinfo`: 天气信息

### 硬件定位
- `HardwarePosition`: 硬件定位 (v1)
- `HardwarePositionV5`: 硬件定位 (v5)

### 公交查询
- `BusStationID`: 公交站 ID 查询
- `BusStationKeyword`: 公交站关键字查询
- `BusLineID`: 公交路线 ID 查询
- `BusLineKeyword`: 公交路线关键字查询

## 使用示例

### 地理编码

```go
req := &geoCode.GeocodeRequest{
    Address: "北京市朝阳区阜通东大街6号",
    City:    "北京",
}

resp, err := client.GeoCode(req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("经纬度:", resp.Geocodes[0].Location)
```

### 驾车路径规划

```go
req := &drivingV2.DrivingRequestV2{
    Origin:      "116.481028,39.989643",
    Destination: "116.514203,39.905409",
}

resp, err := client.DrivingV2(req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("距离:", resp.Route.Paths[0].Distance)
fmt.Println("耗时:", resp.Route.Paths[0].Duration)
```

### 天气信息

```go
req := &weatherinfo.WeatherinfoRequest{
    City:  "北京",
    Extensions: "base",
}

resp, err := client.Weatherinfo(req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("天气:", resp.Lives[0].Weather)
fmt.Println("温度:", resp.Lives[0].Temperature)
```

## 配置选项

| 配置项 | 类型 | 说明 | 是否必填 |
|-------|------|------|----------|
| Key | string | 高德 API Key | 是 |
| SecurityKey | string | 高德 API 安全密钥，用于生成签名 | 否 |
| Timeout | time.Duration | 请求超时时间 | 否，默认 30 秒 |
| Proxy | string | 代理地址 | 否 |
| BaseURL | string | API 基础 URL | 否，默认 https://restapi.amap.com/v3 |
| UserAgent | string | HTTP 请求 User-Agent | 否，默认 amap-go-client/1.0 |

## 错误处理

所有 API 调用都会返回标准的 Go 错误，错误类型包括：

- `InvalidConfigError`: 配置错误
- `NetworkError`: 网络错误
- `APIError`: API 返回的错误

您可以使用类型断言来处理特定类型的错误：

```go
resp, err := client.GeoCode(req)
if err != nil {
    if apiErr, ok := err.(*amapErr.APIError); ok {
        fmt.Printf("API 错误: %s (代码: %s)\n", apiErr.Message, apiErr.Code)
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 注意事项

1. 请确保您已经在高德开放平台申请了相应的 API Key
2. 不同的 API 服务可能需要不同的权限，请在高德开放平台控制台中开启相应的服务
3. 部分 API 服务可能会产生费用，请参考高德开放平台的定价策略
4. 请合理使用 API，避免频繁调用导致的限流
5. **AI 生成代码**: 本项目包含 AI 辅助生成的代码，使用前请注意：
   - 建议在生产环境部署前进行充分的测试验证
   - 定期检查和更新 AI 生成的代码以确保质量和安全性
   - 对于关键业务逻辑，建议进行人工代码审查

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
