package amap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	bicycling "github.com/enneket/amap/api/direction/bicycling"
	driving "github.com/enneket/amap/api/direction/driving"
	walking "github.com/enneket/amap/api/direction/walking"
	distance "github.com/enneket/amap/api/distance"
	geoCode "github.com/enneket/amap/api/geo_code"
	reGeoCode "github.com/enneket/amap/api/re_geo_code"
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

// Walking 步行路径规划API调用方法
func (c *Client) Walking(req *walking.WalkingRequest) (*walking.WalkingResponse, error) {
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
	var resp walking.WalkingResponse
	if err := c.DoRequest("direction/walking", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Driving 驾车路径规划API调用方法
func (c *Client) Driving(req *driving.DrivingRequest) (*driving.DrivingResponse, error) {
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
	var resp driving.DrivingResponse
	if err := c.DoRequest("direction/driving", params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Bicycling 骑行路径规划API调用方法
func (c *Client) Bicycling(req *bicycling.BicyclingRequest) (*bicycling.BicyclingResponse, error) {
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
	var resp bicycling.BicyclingResponse
	if err := c.DoRequest("direction/bicycling", params, &resp); err != nil {
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
