package amap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	busLineID "github.com/enneket/amap/api/bus/line_id"
	busLineKeyword "github.com/enneket/amap/api/bus/line_keyword"
	busStationID "github.com/enneket/amap/api/bus/station_id"
	busStationKeyword "github.com/enneket/amap/api/bus/station_keyword"
	convert "github.com/enneket/amap/api/convert"
	bicyclingV1 "github.com/enneket/amap/api/direction/v1/bicycling"
	drivingV1 "github.com/enneket/amap/api/direction/v1/driving"
	walkingV1 "github.com/enneket/amap/api/direction/v1/walking"
	bicyclingV2 "github.com/enneket/amap/api/direction/v2/bicycling"
	busV2 "github.com/enneket/amap/api/direction/v2/bus"
	drivingV2 "github.com/enneket/amap/api/direction/v2/driving"
	electricV2 "github.com/enneket/amap/api/direction/v2/electric"
	walkingV2 "github.com/enneket/amap/api/direction/v2/walking"
	distance "github.com/enneket/amap/api/distance"
	district "github.com/enneket/amap/api/district"
	geoCode "github.com/enneket/amap/api/geo_code"
	grasproad "github.com/enneket/amap/api/grasproad"
	"github.com/enneket/amap/api/inputtips"
	ipV3 "github.com/enneket/amap/api/ip/v3"
	ipV5 "github.com/enneket/amap/api/ip/v5"
	placev3aoi "github.com/enneket/amap/api/place/v3/aoi"
	placev3around "github.com/enneket/amap/api/place/v3/around"
	placev3id "github.com/enneket/amap/api/place/v3/id"
	placev3polygon "github.com/enneket/amap/api/place/v3/polygon"
	placev3text "github.com/enneket/amap/api/place/v3/text"
	placev5aoi "github.com/enneket/amap/api/place/v5/aoi"
	placev5around "github.com/enneket/amap/api/place/v5/around"
	placev5id "github.com/enneket/amap/api/place/v5/id"
	placev5polygon "github.com/enneket/amap/api/place/v5/polygon"
	placev5text "github.com/enneket/amap/api/place/v5/text"
	positionV1 "github.com/enneket/amap/api/position/v1"
	positionV5 "github.com/enneket/amap/api/position/v5"
	reGeoCode "github.com/enneket/amap/api/re_geo_code"
	trafficIncident "github.com/enneket/amap/api/traffic-incident"
	circle "github.com/enneket/amap/api/traffic-situation/circle"
	line "github.com/enneket/amap/api/traffic-situation/line"
	rectangle "github.com/enneket/amap/api/traffic-situation/rectangle"
	"github.com/enneket/amap/api/weatherinfo"
	amapErr "github.com/enneket/amap/errors"
	amapType "github.com/enneket/amap/types"
	"github.com/enneket/amap/utils"
)

// Client 高德 API 客户端
type Client struct {
	config     *Config      // 全局配置
	httpClient *http.Client // 复用的 HTTP 客户端
}

// NewClient 创建客户端实例（校验配置合法性）
func NewClient(cfg *Config) (*Client, error) {
	if cfg.Key == "" {
		return nil, amapErr.NewInvalidConfigError("API Key 不能为空")
	}
	// 初始化 HTTP 客户端（支持超时、代理）
	httpClient := &http.Client{Timeout: cfg.Timeout}
	if cfg.Proxy != "" {
		proxyURL, _ := url.Parse(cfg.Proxy)
		httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	return &Client{config: cfg, httpClient: httpClient}, nil
}

// DoRequest 通用请求方法（封装公共参数、签名、响应解析）
func (c *Client) DoRequest(path string, params map[string]string, resp interface{}) error {
	// 1. 合并公共参数（Key、签名、Timestamp 等）
	allParams := c.buildPublicParams(params)
	// 2. 签名（如果配置了 SecurityKey）
	if c.config.SecurityKey != "" {
		allParams["sig"] = utils.Sign(allParams, c.config.SecurityKey)
	}
	// 3. 构建请求 URL
	fullURL := c.config.BaseURL + path + "?" + utils.EncodeParams(allParams, true)
	// 4. 发送 HTTP 请求
	req, _ := http.NewRequest(http.MethodGet, fullURL, nil)
	req.Header.Set("User-Agent", c.config.UserAgent)
	rawResp, err := c.httpClient.Do(req)
	if err != nil {
		return amapErr.NewNetworkError(err.Error())
	}
	defer rawResp.Body.Close()
	// 5. 解析响应（先解析基础响应，再解析业务响应）
	baseResp, _, err := amapType.ReadBaseResponse(rawResp.Body)
	if err != nil {
		return err
	}
	// 6. 校验 API 错误
	if baseResp.Status != "1" {
		return amapErr.NewAPIError(baseResp.InfoCode, baseResp.Info)
	}
	// 7. 解析到业务响应结构体
	return json.Unmarshal(baseResp.RawJSON, resp)
}

// buildPublicParams 构建公共参数（Key、Timestamp 等）
func (c *Client) buildPublicParams(params map[string]string) map[string]string {
	publicParams := map[string]string{
		"key":       c.config.Key,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"output":    "JSON",
	}
	// 合并业务参数
	for k, v := range params {
		if v != "" { // 忽略空参数
			publicParams[k] = v
		}
	}
	return publicParams
}

// GeoCode 地理编码API调用方法（为amap.Client绑定方法）
func (c *Client) GeoCode(req *geoCode.GeocodeRequest) (*geoCode.GeoCodeResponse, error) {
	// 校验必填参数
	if req.Address == "" {
		return nil, amapErr.NewInvalidConfigError("地理编码：address参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp geoCode.GeoCodeResponse
	if err := c.DoRequest("geocode/geo", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ReGeocode 逆地理编码API调用方法
func (c *Client) ReGeocode(req *reGeoCode.ReGeocodeRequest) (*reGeoCode.ReGeocodeResponse, error) {
	// 校验必填参数
	if req.Location == "" {
		return nil, amapErr.NewInvalidConfigError("逆地理编码：location参数不能为空")
	}
	// 简单校验经纬度格式（经度,纬度）
	if !strings.Contains(req.Location, ",") {
		return nil, amapErr.NewInvalidConfigError("逆地理编码：location格式错误，应为\"经度,纬度\"")
	}

	// 处理默认值
	if req.Radius <= 0 {
		req.Radius = 1000 // 默认搜索半径1000米
	}
	if req.Extensions == "" {
		req.Extensions = "base" // 默认返回基础信息
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp reGeoCode.ReGeocodeResponse
	if err := c.DoRequest("geocode/regeo", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Walking 步行路径规划API调用方法（v1）
func (c *Client) Walking(req *walkingV1.WalkingRequest) (*walkingV1.WalkingResponse, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("步行路径规划：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("步行路径规划：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("步行路径规划：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp walkingV1.WalkingResponse
	if err := c.DoRequest("direction/walking", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Driving 驾车路径规划API调用方法（v1）
func (c *Client) Driving(req *drivingV1.DrivingRequest) (*drivingV1.DrivingResponse, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("驾车路径规划：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("驾车路径规划：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("驾车路径规划：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp drivingV1.DrivingResponse
	if err := c.DoRequest("direction/driving", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Bicycling 骑行路径规划API调用方法（v1）
func (c *Client) Bicycling(req *bicyclingV1.BicyclingRequest) (*bicyclingV1.BicyclingResponse, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("骑行路径规划：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("骑行路径规划：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("骑行路径规划：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp bicyclingV1.BicyclingResponse
	if err := c.DoRequest("direction/bicycling", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// WalkingV2 步行路径规划API调用方法（v2）
func (c *Client) WalkingV2(req *walkingV2.WalkingRequestV2) (*walkingV2.WalkingResponseV2, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("步行路径规划v2：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("步行路径规划v2：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("步行路径规划v2：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp walkingV2.WalkingResponseV2
	if err := c.DoRequest("direction/v2/walking", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DrivingV2 驾车路径规划API调用方法（v2）
func (c *Client) DrivingV2(req *drivingV2.DrivingRequestV2) (*drivingV2.DrivingResponseV2, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("驾车路径规划v2：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("驾车路径规划v2：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("驾车路径规划v2：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp drivingV2.DrivingResponseV2
	if err := c.DoRequest("direction/v2/driving", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BicyclingV2 骑行路径规划API调用方法（v2）
func (c *Client) BicyclingV2(req *bicyclingV2.BicyclingRequestV2) (*bicyclingV2.BicyclingResponseV2, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("骑行路径规划v2：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("骑行路径规划v2：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("骑行路径规划v2：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp bicyclingV2.BicyclingResponseV2
	if err := c.DoRequest("direction/v2/bicycling", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BusV2 公交路线规划API调用方法（v2）
func (c *Client) BusV2(req *busV2.BusRequestV2) (*busV2.BusResponseV2, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("公交路径规划v2：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("公交路径规划v2：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("公交路径规划v2：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp busV2.BusResponseV2
	if err := c.DoRequest("direction/v2/transit/integrated", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ElectricV2 电动车路线规划API调用方法（v2）
func (c *Client) ElectricV2(req *electricV2.ElectricRequestV2) (*electricV2.ElectricResponseV2, error) {
	// 校验必填参数
	if req.Origin == "" {
		return nil, amapErr.NewInvalidConfigError("电动车路径规划v2：origin参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("电动车路径规划v2：destination参数不能为空")
	}
	// 简单校验经纬度格式
	if !strings.Contains(req.Origin, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("电动车路径规划v2：坐标格式错误，应为\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp electricV2.ElectricResponseV2
	if err := c.DoRequest("direction/v2/electric", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Distance 距离测量API调用方法
func (c *Client) Distance(req *distance.DistanceRequest) (*distance.DistanceResponse, error) {
	// 校验必填参数
	if req.Origins == "" {
		return nil, amapErr.NewInvalidConfigError("距离测量：origins参数不能为空")
	}
	if req.Destination == "" {
		return nil, amapErr.NewInvalidConfigError("距离测量：destination参数不能为空")
	}

	// 简单校验经纬度格式（至少包含一个有效的经纬度对）
	if !strings.Contains(req.Origins, ",") || !strings.Contains(req.Destination, ",") {
		return nil, amapErr.NewInvalidConfigError("距离测量：坐标格式错误，应为\"经度,纬度|经度,纬度\"或单个\"经度,纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp distance.DistanceResponse
	if err := c.DoRequest("direction/distance", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// District 行政区查询API调用方法
// 支持通过关键字搜索行政区，可指定返回子级行政区的级别
func (c *Client) District(req *district.DistrictRequest) (*district.DistrictResponse, error) {
	// 校验必填参数
	if req.Keywords == "" {
		return nil, amapErr.NewInvalidConfigError("行政区查询：keywords参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp district.DistrictResponse
	if err := c.DoRequest("config/district", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// TrafficIncident 交通事件查询API调用方法
// 支持查询指定区域内的交通事件，可按事件级别和类型筛选
func (c *Client) TrafficIncident(req *trafficIncident.TrafficIncidentRequest) (*trafficIncident.TrafficIncidentResponse, error) {
	// 校验必填参数
	if req.Level == "" {
		return nil, amapErr.NewInvalidConfigError("交通事件查询：level参数不能为空")
	}
	if req.Type == "" {
		return nil, amapErr.NewInvalidConfigError("交通事件查询：type参数不能为空")
	}
	if req.Rectangle == "" {
		return nil, amapErr.NewInvalidConfigError("交通事件查询：rectangle参数不能为空")
	}

	// 简单校验矩形区域格式（必须包含3个逗号，如："116.351147,39.904989,116.480317,39.976564"）
	commaCount := 0
	for _, char := range req.Rectangle {
		if char == ',' {
			commaCount++
		}
	}
	if commaCount != 3 {
		return nil, amapErr.NewInvalidConfigError("交通事件查询：rectangle格式错误，应为\"左下经度,左下纬度,右上经度,右上纬度\"")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp trafficIncident.TrafficIncidentResponse
	if err := c.DoRequest("traffic/status", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// IPConfig IP定位API调用方法
// 支持通过IP地址查询地理位置信息，返回省份、城市、区县、ISP等信息
func (c *Client) IPConfig(req *ipV3.IPConfigRequest) (*ipV3.IPConfigResponse, error) {
	// 校验必填参数
	if req.IP == "" {
		return nil, amapErr.NewInvalidConfigError("IP定位：ip参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp ipV3.IPConfigResponse
	if err := c.DoRequest("ip", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// IPV5Config IP定位API调用方法（v5）
// 支持通过IP地址查询地理位置信息，返回省份、城市、区县、ISP等信息
func (c *Client) IPV5Config(req *ipV5.IPConfigRequest) (*ipV5.IPConfigResponse, error) {
	// 校验必填参数
	if req.IP == "" {
		return nil, amapErr.NewInvalidConfigError("IP定位v5：ip参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp ipV5.IPConfigResponse
	if err := c.DoRequest("v5/ip", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Convert 坐标转换API调用方法
// 支持将其他坐标系的坐标转换为高德坐标系（GCJ02）
// 支持批量转换，一次最多转换40对坐标
func (c *Client) Convert(req *convert.ConvertRequest) (*convert.ConvertResponse, error) {
	// 校验必填参数
	if req.Locations == "" {
		return nil, amapErr.NewInvalidConfigError("坐标转换：locations参数不能为空")
	}
	if req.CoordSys == "" {
		return nil, amapErr.NewInvalidConfigError("坐标转换：coordsys参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp convert.ConvertResponse
	if err := c.DoRequest("convert", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GraspRoad 轨迹纠偏API调用方法
// 用于将原始轨迹点转换为匹配道路的轨迹点，支持批量处理
func (c *Client) GraspRoad(req *grasproad.GraspRoadRequest) (*grasproad.GraspRoadResponse, error) {
	// 校验必填参数
	if req.SID == "" {
		return nil, amapErr.NewInvalidConfigError("轨迹纠偏：sid参数不能为空")
	}
	if req.Points == "" {
		return nil, amapErr.NewInvalidConfigError("轨迹纠偏：points参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp grasproad.GraspRoadResponse
	if err := c.DoRequest("grasproad", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV3ID POI ID查询API调用方法（v3）
// 根据POI ID查询详细信息
func (c *Client) PlaceV3ID(req *placev3id.IDRequest) (*placev3id.IDResponse, error) {
	// 校验必填参数
	if req.ID == "" {
		return nil, amapErr.NewInvalidConfigError("POI ID查询：id参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev3id.IDResponse
	if err := c.DoRequest("place/detail", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV3Text POI文本搜索API调用方法（v3）
// 基于关键词的搜索，支持矩形范围搜索
func (c *Client) PlaceV3Text(req *placev3text.TextSearchRequest) (*placev3text.TextSearchResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev3text.TextSearchResponse
	if err := c.DoRequest("place/text", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV3Around POI周边搜索API调用方法（v3）
// 基于中心点和半径的搜索，用于查询指定区域内的POI
func (c *Client) PlaceV3Around(req *placev3around.AroundSearchRequest) (*placev3around.AroundSearchResponse, error) {
	// 校验必填参数
	if req.Location == "" {
		return nil, amapErr.NewInvalidConfigError("POI周边搜索：location参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev3around.AroundSearchResponse
	if err := c.DoRequest("place/around", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV3Polygon POI多边形搜索API调用方法（v3）
// 基于多边形边界的搜索，用于查询指定多边形区域内的POI
func (c *Client) PlaceV3Polygon(req *placev3polygon.PolygonSearchRequest) (*placev3polygon.PolygonSearchResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev3polygon.PolygonSearchResponse
	if err := c.DoRequest("place/polygon", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV3AOI POI AOI查询API调用方法（v3）
// 用于查询指定AOI区域内的POI
func (c *Client) PlaceV3AOI(req *placev3aoi.AOISearchRequest) (*placev3aoi.AOISearchResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev3aoi.AOISearchResponse
	if err := c.DoRequest("place/aoi", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV5ID POI ID查询API调用方法（v5）
// 根据POI ID查询详细信息
func (c *Client) PlaceV5ID(req *placev5id.IDRequest) (*placev5id.IDResponse, error) {
	// 校验必填参数
	if req.ID == "" {
		return nil, amapErr.NewInvalidConfigError("POI ID查询v5：id参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev5id.IDResponse
	if err := c.DoRequest("v5/place/detail", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV5Text POI文本搜索API调用方法（v5）
// 基于关键词的搜索，用于查询指定区域内的POI
func (c *Client) PlaceV5Text(req *placev5text.TextSearchRequest) (*placev5text.TextSearchResponse, error) {
	// 校验必填参数
	if req.Keyword == "" {
		return nil, amapErr.NewInvalidConfigError("POI文本搜索v5：keyword参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev5text.TextSearchResponse
	if err := c.DoRequest("v5/place/text", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV5Around POI周边搜索API调用方法（v5）
// 基于中心点和半径的搜索，用于查询指定区域内的POI
func (c *Client) PlaceV5Around(req *placev5around.AroundSearchRequest) (*placev5around.AroundSearchResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev5around.AroundSearchResponse
	if err := c.DoRequest("v5/place/around", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV5Polygon POI多边形搜索API调用方法（v5）
// 基于多边形边界的搜索，用于查询指定多边形区域内的POI
func (c *Client) PlaceV5Polygon(req *placev5polygon.PolygonSearchRequest) (*placev5polygon.PolygonSearchResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev5polygon.PolygonSearchResponse
	if err := c.DoRequest("v5/place/polygon", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PlaceV5AOI POI AOI查询API调用方法（v5）
// 用于查询指定AOI区域内的POI
func (c *Client) PlaceV5AOI(req *placev5aoi.AOISearchRequest) (*placev5aoi.AOISearchResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp placev5aoi.AOISearchResponse
	if err := c.DoRequest("v5/place/aoi", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Inputtips 输入提示API调用方法
// 支持根据关键字获取输入提示，可指定城市、类型等过滤条件
func (c *Client) Inputtips(req *inputtips.InputtipsRequest) (*inputtips.InputtipsResponse, error) {
	// 校验必填参数
	if req.Keywords == "" {
		return nil, amapErr.NewInvalidConfigError("输入提示：keywords参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp inputtips.InputtipsResponse
	if err := c.DoRequest("assistant/inputtips", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Weatherinfo 天气信息API调用方法
// 支持获取实时天气、预报信息和生活指数建议
func (c *Client) Weatherinfo(req *weatherinfo.WeatherinfoRequest) (*weatherinfo.WeatherinfoResponse, error) {
	// 校验必填参数
	if req.City == "" {
		return nil, amapErr.NewInvalidConfigError("天气信息：city参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp weatherinfo.WeatherinfoResponse
	if err := c.DoRequest("weather/weatherInfo", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// HardwarePosition 硬件定位API调用方法（v1）
// 支持通过硬件设备信息（如GPS、基站、WiFi等）获取地理位置信息
func (c *Client) HardwarePosition(req *positionV1.HardwarePositionRequest) (*positionV1.HardwarePositionResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp positionV1.HardwarePositionResponse
	if err := c.DoRequest("position/v1/hardware", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// HardwarePositionV5 硬件定位API调用方法（v5）
// 支持通过硬件设备信息（如GPS、基站、WiFi等）获取地理位置信息，v5版本增强了定位精度和多源数据融合能力
func (c *Client) HardwarePositionV5(req *positionV5.HardwarePositionRequest) (*positionV5.HardwarePositionResponse, error) {
	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp positionV5.HardwarePositionResponse
	if err := c.DoRequest("v5/position/hardware", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// LineTrafficStatus 指定线路交通态势查询API调用方法
// 支持查询指定线路的交通态势信息
func (c *Client) LineTrafficStatus(req *line.LineTrafficRequest) (*line.LineTrafficResponse, error) {
	// 校验必填参数
	if req.Path == "" {
		return nil, amapErr.NewInvalidConfigError("指定线路交通态势查询：path参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp line.LineTrafficResponse
	if err := c.DoRequest("v3/traffic/status/road", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CircleTrafficStatus 圆形区域内交通态势查询API调用方法
// 支持查询指定圆形区域内的交通态势信息
func (c *Client) CircleTrafficStatus(req *circle.CircleTrafficRequest) (*circle.CircleTrafficResponse, error) {
	// 校验必填参数
	if req.Center == "" {
		return nil, amapErr.NewInvalidConfigError("圆形区域交通态势查询：center参数不能为空")
	}
	if req.Radius == "" {
		return nil, amapErr.NewInvalidConfigError("圆形区域交通态势查询：radius参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp circle.CircleTrafficResponse
	if err := c.DoRequest("v3/traffic/status/circle", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RectangleTrafficStatus 矩形区域内交通态势查询API调用方法
// 支持查询指定矩形区域内的交通态势信息
func (c *Client) RectangleTrafficStatus(req *rectangle.RectangleTrafficRequest) (*rectangle.RectangleTrafficResponse, error) {
	// 校验必填参数
	if req.Rectangle == "" {
		return nil, amapErr.NewInvalidConfigError("矩形区域交通态势查询：rectangle参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp rectangle.RectangleTrafficResponse
	if err := c.DoRequest("v3/traffic/status/rectangle", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BusStationID 公交站ID查询API调用方法
// 根据公交站点ID查询经过该站点的所有公交线路详细信息
func (c *Client) BusStationID(req *busStationID.StationIDRequest) (*busStationID.StationIDResponse, error) {
	// 校验必填参数
	if req.ID == "" {
		return nil, amapErr.NewInvalidConfigError("公交站ID查询：id参数不能为空")
	}
	if req.City == "" {
		return nil, amapErr.NewInvalidConfigError("公交站ID查询：city参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp busStationID.StationIDResponse
	if err := c.DoRequest("bus/linename", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BusStationKeyword 公交站关键字查询API调用方法
// 根据公交站点名称关键字查询公交站点及经过该站点的公交线路信息
func (c *Client) BusStationKeyword(req *busStationKeyword.StationKeywordRequest) (*busStationKeyword.StationKeywordResponse, error) {
	// 校验必填参数
	if req.Keywords == "" {
		return nil, amapErr.NewInvalidConfigError("公交站关键字查询：keywords参数不能为空")
	}
	if req.City == "" {
		return nil, amapErr.NewInvalidConfigError("公交站关键字查询：city参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp busStationKeyword.StationKeywordResponse
	if err := c.DoRequest("bus/station/search", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BusLineID 公交路线ID查询API调用方法
// 根据公交线路ID查询该线路的详细信息
func (c *Client) BusLineID(req *busLineID.LineIDRequest) (*busLineID.LineIDResponse, error) {
	// 校验必填参数
	if req.ID == "" {
		return nil, amapErr.NewInvalidConfigError("公交路线ID查询：id参数不能为空")
	}
	if req.City == "" {
		return nil, amapErr.NewInvalidConfigError("公交路线ID查询：city参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp busLineID.LineIDResponse
	if err := c.DoRequest("bus/lineid", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// BusLineKeyword 公交路线关键字查询API调用方法
// 根据公交线路名称关键字查询公交线路详细信息
func (c *Client) BusLineKeyword(req *busLineKeyword.LineKeywordRequest) (*busLineKeyword.LineKeywordResponse, error) {
	// 校验必填参数
	if req.Keywords == "" {
		return nil, amapErr.NewInvalidConfigError("公交路线关键字查询：keywords参数不能为空")
	}
	if req.City == "" {
		return nil, amapErr.NewInvalidConfigError("公交路线关键字查询：city参数不能为空")
	}

	// 转换请求参数为map
	params := req.ToParams()

	// 调用核心请求方法
	var resp busLineKeyword.LineKeywordResponse
	if err := c.DoRequest("bus/line/search", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
