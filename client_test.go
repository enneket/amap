package amap

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	busLineID "github.com/enneket/amap/api/bus/line_id"
	busLineKeyword "github.com/enneket/amap/api/bus/line_keyword"
	busStationID "github.com/enneket/amap/api/bus/station_id"
	busStationKeyword "github.com/enneket/amap/api/bus/station_keyword"
	convert "github.com/enneket/amap/api/convert"
	bicycling "github.com/enneket/amap/api/direction/v1/bicycling"
	driving "github.com/enneket/amap/api/direction/v1/driving"
	walking "github.com/enneket/amap/api/direction/v1/walking"
	bicyclingV2 "github.com/enneket/amap/api/direction/v2/bicycling"
	busV2 "github.com/enneket/amap/api/direction/v2/bus"
	drivingV2 "github.com/enneket/amap/api/direction/v2/driving"
	electricV2 "github.com/enneket/amap/api/direction/v2/electric"
	walkingV2 "github.com/enneket/amap/api/direction/v2/walking"
	distance "github.com/enneket/amap/api/distance"
	district "github.com/enneket/amap/api/district"
	etdDrivingV4 "github.com/enneket/amap/api/etd/v4/driving"
	geoCode "github.com/enneket/amap/api/geo_code"
	grasproad "github.com/enneket/amap/api/grasproad"
	inputtips "github.com/enneket/amap/api/input_tips"
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
	trafficIncident "github.com/enneket/amap/api/traffic_incident"
	circle "github.com/enneket/amap/api/traffic_situation/circle"
	line "github.com/enneket/amap/api/traffic_situation/line"
	rectangle "github.com/enneket/amap/api/traffic_situation/rectangle"
	"github.com/enneket/amap/api/weatherinfo"
	amapErr "github.com/enneket/amap/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockResponse 用于构建测试响应
func mockResponse(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
}

// 测试业务响应结构体
type TestResponse struct {
	Result string `json:"result"`
}

// TestDoRequest_Success 测试DoRequest方法正常请求成功
func TestDoRequest_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"result": "test success"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest(mockServer.URL+"/test/path", map[string]string{"param1": "value1"}, &resp)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "test success", resp.Result)
}

// TestDoRequest_APIError 测试API返回错误（status != "1"）
func TestDoRequest_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest(mockServer.URL+"/test/path", nil, &resp)

	// 4. 验证结果
	assert.Error(t, err)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestPlaceV3ID_Success 测试PlaceV3ID方法成功
func TestPlaceV3ID_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"name": "测试POI",
		"type": "010000",
		"location": "116.397428,39.90923"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev3id.IDRequest{ID: "B000A83M61"}
	resp, err := client.PlaceV3ID(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10000", resp.InfoCode)
}

// TestPlaceV3Text_Success 测试PlaceV3Text方法成功
func TestPlaceV3Text_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev3text.TextSearchRequest{}
	resp, err := client.PlaceV3Text(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestPlaceV3Around_Success 测试PlaceV3Around方法成功
func TestPlaceV3Around_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev3around.AroundSearchRequest{Location: "116.397428,39.90923"}
	resp, err := client.PlaceV3Around(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestPlaceV3Polygon_Success 测试PlaceV3Polygon方法成功
func TestPlaceV3Polygon_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev3polygon.PolygonSearchRequest{}
	resp, err := client.PlaceV3Polygon(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestPlaceV3AOI_Success 测试PlaceV3AOI方法成功
func TestPlaceV3AOI_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev3aoi.AOISearchRequest{}
	resp, err := client.PlaceV3AOI(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestPlaceV5ID_Success 测试PlaceV5ID方法成功
func TestPlaceV5ID_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev5id.IDRequest{ID: "B000A83M61"}
	resp, err := client.PlaceV5ID(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "10000", resp.InfoCode)
}

// TestPlaceV5Text_Success 测试PlaceV5Text方法成功
func TestPlaceV5Text_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev5text.TextSearchRequest{Keyword: "测试"}
	resp, err := client.PlaceV5Text(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestPlaceV5Around_Success 测试PlaceV5Around方法成功
func TestPlaceV5Around_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev5around.AroundSearchRequest{Location: "116.397428,39.90923"}
	resp, err := client.PlaceV5Around(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestPlaceV5Polygon_Success 测试PlaceV5Polygon方法成功
func TestPlaceV5Polygon_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev5polygon.PolygonSearchRequest{Polygon: "116.397428,39.90923;116.407428,39.90923"}
	resp, err := client.PlaceV5Polygon(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestPlaceV5AOI_Success 测试PlaceV5AOI方法成功
func TestPlaceV5AOI_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"pois": [{
			"id": "B000A83M61",
			"name": "测试POI",
			"type": "010000",
			"location": "116.397428,39.90923"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &placev5aoi.AOISearchRequest{ID: "B000A83M61"}
	resp, err := client.PlaceV5AOI(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Count)
}

// TestDoRequest_ParseError 测试响应解析错误
func TestDoRequest_ParseError(t *testing.T) {
	// 1. 创建mock服务器，返回无效的JSON
	mockServer := mockResponse(http.StatusOK, "invalid json")
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest(mockServer.URL+"/test/path", nil, &resp)

	// 4. 验证结果
	assert.Error(t, err)
	assert.IsType(t, amapErr.ParseError(""), err)
}

// TestDoRequest_NetworkError 测试网络错误
func TestDoRequest_NetworkError(t *testing.T) {
	// 1. 创建Client实例，使用不存在的地址
	config := NewConfig("test_key")

	config.Timeout = 100 * time.Millisecond
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 执行请求
	var resp TestResponse
	err = client.DoRequest("test/path", nil, &resp)

	// 3. 验证结果
	assert.Error(t, err)
	assert.IsType(t, amapErr.NetworkError(""), err)
}

// TestDoRequest_WithSignature 测试带签名的请求
func TestDoRequest_WithSignature(t *testing.T) {
	// 1. 创建mock服务器，验证签名参数
	var receivedParams url.Values
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		receivedParams = r.URL.Query()

		// 返回成功响应
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": "1",
			"info": "OK",
			"infocode": "10000",
			"result": "test with signature"
		}`))
	}))
	defer mockServer.Close()

	// 2. 创建Client实例，配置SecurityKey
	config := NewConfig("test_key")

	config.SecurityKey = "test_security_key"
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest(mockServer.URL+"/test/path", map[string]string{"param1": "value1"}, &resp)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "test with signature", resp.Result)

	// 5. 验证请求参数包含签名
	assert.Equal(t, "test_key", receivedParams.Get("key"))
	assert.Equal(t, "value1", receivedParams.Get("param1"))
	assert.Equal(t, "JSON", receivedParams.Get("output"))
	assert.NotEmpty(t, receivedParams.Get("timestamp"))
	assert.NotEmpty(t, receivedParams.Get("sig")) // 验证包含签名
}

// TestDoRequest_PublicParams 测试公共参数的构建
func TestDoRequest_PublicParams(t *testing.T) {
	// 1. 创建mock服务器，验证公共参数
	var receivedParams url.Values
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		receivedParams = r.URL.Query()

		// 返回成功响应
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": "1",
			"info": "OK",
			"infocode": "10000",
			"result": "test public params"
		}`))
	}))
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求，传入部分参数
	var resp TestResponse
	customUserAgent := "test-agent/1.0"
	config.UserAgent = customUserAgent
	err = client.DoRequest(mockServer.URL+"/test/path", map[string]string{"param1": "value1", "param2": ""}, &resp)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "test public params", resp.Result)

	// 5. 验证公共参数
	assert.Equal(t, "test_key", receivedParams.Get("key"))
	assert.Equal(t, "value1", receivedParams.Get("param1"))
	assert.Equal(t, "", receivedParams.Get("param2")) // 空参数应该被过滤掉
	assert.Equal(t, "JSON", receivedParams.Get("output"))
	assert.NotEmpty(t, receivedParams.Get("timestamp"))
}

// TestDoRequest_InvalidJSON 测试响应JSON无法解析到目标结构体
func TestDoRequest_InvalidJSON(t *testing.T) {
	// 1. 创建mock服务器，返回的JSON结构与目标结构体不匹配
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"result": {"nested": "value"}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求，使用错误的结构体类型接收
	var resp struct {
		Result string `json:"result"` // 期望string类型，但实际是object
	}
	err = client.DoRequest("test/path", nil, &resp)

	// 4. 验证结果
	assert.Error(t, err) // 应该返回JSON解析错误
}

// TestBuildPublicParams 测试buildPublicParams方法（虽然是私有方法，但可以通过DoRequest间接测试）
func TestBuildPublicParams(t *testing.T) {
	// 1. 创建mock服务器，记录接收到的参数
	var receivedTimestamp string
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedTimestamp = r.URL.Query().Get("timestamp")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": "1",
			"info": "OK",
			"infocode": "10000",
			"result": "ok"
		}`))
	}))
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest(mockServer.URL+"/test/path", nil, &resp)
	require.NoError(t, err)

	// 4. 验证时间戳格式正确
	_, err = strconv.ParseInt(receivedTimestamp, 10, 64)
	assert.NoError(t, err, "timestamp should be a valid integer")

	// 5. 验证时间戳接近当前时间（10秒内）
	timestamp := time.Now().Unix()
	receivedTs, _ := strconv.ParseInt(receivedTimestamp, 10, 64)
	assert.InDelta(t, timestamp, receivedTs, 10, "timestamp should be close to current time")
}

// TestGeoCode_Success 测试GeoCode方法正常请求成功
func TestGeoCode_Success(t *testing.T) {
	// 1. 创建mock服务器，返回地理编码成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"geocodes": [{
			"formatted_address": "北京市朝阳区望京SOHO",
			"country": "中国",
			"province": "北京市",
			"city": "北京市",
			"citycode": "110000",
			"district": "朝阳区",
			"adcode": "110105",
			"street": "望京街",
			"number": "8号",
			"location": "116.48649,39.99947",
			"level": "门牌号"
		}]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &geoCode.GeocodeRequest{
		Address: "北京市朝阳区望京SOHO",
		City:    "北京",
	}

	// 4. 执行地理编码请求
	resp, err := client.GeoCode(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "1", resp.Count)
	assert.Len(t, resp.Geocodes, 1)
	assert.Equal(t, "北京市朝阳区望京SOHO", resp.Geocodes[0].FormattedAddress)
	assert.Equal(t, "中国", resp.Geocodes[0].Country)
	assert.Equal(t, "北京市", resp.Geocodes[0].Province)
	assert.Equal(t, "北京市", resp.Geocodes[0].City)
	assert.Equal(t, "110000", resp.Geocodes[0].Citycode)
	assert.Equal(t, "朝阳区", resp.Geocodes[0].District)
	assert.Equal(t, "110105", resp.Geocodes[0].Adcode)
	assert.Equal(t, "望京街", resp.Geocodes[0].Street)
	assert.Equal(t, "8号", resp.Geocodes[0].Number)
	assert.Equal(t, "116.48649,39.99947", resp.Geocodes[0].Location)
	assert.Equal(t, "门牌号", resp.Geocodes[0].Level)
}

// TestGeoCode_MissingAddress 测试GeoCode方法缺少必填参数Address
func TestGeoCode_MissingAddress(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Address的请求参数
	req := &geoCode.GeocodeRequest{
		City: "北京", // 只有City，没有Address
	}

	// 3. 执行地理编码请求
	resp, err := client.GeoCode(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "address参数不能为空")
}

// TestGeoCode_APIError 测试GeoCode方法API返回错误
func TestGeoCode_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &geoCode.GeocodeRequest{
		Address: "北京市朝阳区望京SOHO",
	}

	// 4. 执行地理编码请求
	resp, err := client.GeoCode(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestReGeocode_Success 测试ReGeocode方法正常请求成功
func TestReGeocode_Success(t *testing.T) {
	// 1. 创建mock服务器，返回逆地理编码成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"regeocode": {
			"formatted_address": "北京市朝阳区望京街道望京SOHO T1",
			"addressComponent": {
				"country": "中国",
				"province": "北京市",
				"city": "北京市",
				"citycode": "110000",
				"district": "朝阳区",
				"adcode": "110105",
				"township": "望京街道",
				"towncode": "110105028",
				"street": "望京街",
				"number": "8号"
			},
			"township": "望京街道"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &reGeoCode.ReGeocodeRequest{
		Location:   "116.48649,39.99947",
		Radius:     500,
		Extensions: "base",
	}

	// 4. 执行逆地理编码请求
	resp, err := client.ReGeocode(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.NotNil(t, resp.ReGeocode)
	assert.Equal(t, "北京市朝阳区望京街道望京SOHO T1", resp.ReGeocode.FormattedAddress)
	assert.Equal(t, "中国", resp.ReGeocode.AddressComponent.Country)
	assert.Equal(t, "北京市", resp.ReGeocode.AddressComponent.Province)
	assert.Equal(t, "北京市", resp.ReGeocode.AddressComponent.City)
	assert.Equal(t, "110000", resp.ReGeocode.AddressComponent.Citycode)
	assert.Equal(t, "朝阳区", resp.ReGeocode.AddressComponent.District)
	assert.Equal(t, "110105", resp.ReGeocode.AddressComponent.Adcode)
	assert.Equal(t, "望京街道", resp.ReGeocode.AddressComponent.Township)
	assert.Equal(t, "110105028", resp.ReGeocode.AddressComponent.Towncode)
	assert.Equal(t, "望京街", resp.ReGeocode.AddressComponent.Street)
	assert.Equal(t, "8号", resp.ReGeocode.AddressComponent.Number)
	assert.Equal(t, "望京街道", resp.ReGeocode.Township)
}

// TestReGeocode_MissingLocation 测试ReGeocode方法缺少必填参数Location
func TestReGeocode_MissingLocation(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Location的请求参数
	req := &reGeoCode.ReGeocodeRequest{
		Radius: 500,
	}

	// 3. 执行逆地理编码请求
	resp, err := client.ReGeocode(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "逆地理编码：location参数不能为空")
}

// TestReGeocode_InvalidLocationFormat 测试ReGeocode方法Location格式错误
func TestReGeocode_InvalidLocationFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的Location参数
	req := &reGeoCode.ReGeocodeRequest{
		Location: "116.48649 39.99947", // 使用空格分隔而不是逗号
	}

	// 3. 执行逆地理编码请求
	resp, err := client.ReGeocode(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "逆地理编码：location格式错误，应为\"经度,纬度\"")
}

// TestReGeocode_APIError 测试ReGeocode方法API返回错误
func TestReGeocode_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &reGeoCode.ReGeocodeRequest{
		Location: "116.48649,39.99947",
	}

	// 4. 执行逆地理编码请求
	resp, err := client.ReGeocode(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestReGeocode_DefaultValues 测试ReGeocode方法默认值设置
func TestReGeocode_DefaultValues(t *testing.T) {
	// 1. 记录接收到的参数
	var receivedRadius, receivedExtensions string
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRadius = r.URL.Query().Get("radius")
		receivedExtensions = r.URL.Query().Get("extensions")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": "1",
			"info": "OK",
			"infocode": "10000",
			"regeocode": {
				"formatted_address": "北京市朝阳区望京街道",
				"addressComponent": {
					"country": "中国",
					"province": "北京市",
					"city": "北京市",
					"citycode": "110000",
					"district": "朝阳区",
					"adcode": "110105",
					"township": "望京街道"
				},
				"township": "望京街道"
			}
		}`))
	}))
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数，不设置Radius和Extensions
	req := &reGeoCode.ReGeocodeRequest{
		Location: "116.48649,39.99947",
		// 不设置Radius和Extensions，测试默认值
	}

	// 4. 执行逆地理编码请求
	resp, err := client.ReGeocode(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// 验证默认值
	assert.Equal(t, "1000", receivedRadius, "默认radius应该为1000")
	assert.Equal(t, "base", receivedExtensions, "默认extensions应该为base")
}

// TestDistance_Success 测试Distance方法正常请求成功
func TestDistance_Success(t *testing.T) {
	// 1. 创建mock服务器，返回距离测量成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"results": [
			{
				"origin_id": "",
				"dest_id": "",
				"distance": "3237",
				"duration": "324",
				"info": "OK",
				"status": "1"
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &distance.DistanceRequest{
		Origins:     "116.351147,39.936871",
		Destination: "116.410001,39.910113",
		Type:        1,
	}

	// 4. 执行距离测量请求
	resp, err := client.Distance(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Len(t, resp.Results, 1)
	assert.Equal(t, "3237", resp.Results[0].Distance)
	assert.Equal(t, "324", resp.Results[0].Duration)
	assert.Equal(t, "OK", resp.Results[0].Info)
	assert.Equal(t, "1", resp.Results[0].Status)
}

// TestDistance_MissingOrigins 测试Distance方法缺少必填参数Origins
func TestDistance_MissingOrigins(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origins的请求参数
	req := &distance.DistanceRequest{
		Destination: "116.410001,39.910113",
		Type:        1,
	}

	// 3. 执行距离测量请求
	resp, err := client.Distance(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origins参数不能为空")
}

// TestDistance_MissingDestination 测试Distance方法缺少必填参数Destination
func TestDistance_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &distance.DistanceRequest{
		Origins: "116.351147,39.936871",
		Type:    1,
	}

	// 3. 执行距离测量请求
	resp, err := client.Distance(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestDistance_InvalidCoordinateFormat 测试Distance方法坐标格式错误
func TestDistance_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &distance.DistanceRequest{
		Origins:     "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
		Type:        1,
	}

	// 3. 执行距离测量请求
	resp, err := client.Distance(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestDistance_APIError 测试Distance方法API返回错误
func TestDistance_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &distance.DistanceRequest{
		Origins:     "116.351147,39.936871",
		Destination: "116.410001,39.910113",
		Type:        1,
	}

	// 4. 执行距离测量请求
	resp, err := client.Distance(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestDistance_MultiOriginsDestinations 测试Distance方法多起点多终点情况
func TestDistance_MultiOriginsDestinations(t *testing.T) {
	// 1. 创建mock服务器，返回多结果响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"results": [
			{
				"origin_id": "",
				"dest_id": "",
				"distance": "3237",
				"duration": "324",
				"info": "OK",
				"status": "1"
			},
			{
				"origin_id": "",
				"dest_id": "",
				"distance": "8502",
				"duration": "648",
				"info": "OK",
				"status": "1"
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数，使用多起点
	req := &distance.DistanceRequest{
		Origins:     "116.351147,39.936871|116.481247,39.996746",
		Destination: "116.410001,39.910113",
		Type:        2,
	}

	// 4. 执行距离测量请求
	resp, err := client.Distance(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 2)
	assert.Equal(t, "3237", resp.Results[0].Distance)
	assert.Equal(t, "8502", resp.Results[1].Distance)
}

// -------------------------- 步行路径规划测试 --------------------------

// TestWalking_Success 测试Walking方法正常请求成功
func TestWalking_Success(t *testing.T) {
	// 1. 创建mock服务器，返回步行路径规划成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "3237",
					"duration": "324",
					"steps": [
						{
							"instruction": "步行100米，右转进入建国路",
							"orientation": "东",
							"road": "建国路",
							"distance": "100",
							"duration": "60",
							"polyline": "116.351147,39.936871;116.351247,39.936771",
							"action": "右转",
							"assistant_action": ""
						}
					],
					"polyline": "116.351147,39.936871;116.351247,39.936771;116.410001,39.910113",
					"tolls": "0"
				}
			],
			"distance": "3237",
			"duration": "324",
			"tolls": "0"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &walking.WalkingRequest{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行步行路径规划请求
	resp, err := client.Walking(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "3237", resp.Route.Distance)
	assert.Equal(t, "324", resp.Route.Duration)
	assert.Equal(t, "0", resp.Route.Tolls)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Equal(t, "3237", resp.Route.Paths[0].Distance)
	assert.Equal(t, "324", resp.Route.Paths[0].Duration)
	assert.Len(t, resp.Route.Paths[0].Steps, 1)
	assert.Equal(t, "步行100米，右转进入建国路", resp.Route.Paths[0].Steps[0].Instruction)
}

// TestWalking_MissingOrigin 测试Walking方法缺少必填参数Origin
func TestWalking_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &walking.WalkingRequest{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行步行路径规划请求
	resp, err := client.Walking(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestWalking_MissingDestination 测试Walking方法缺少必填参数Destination
func TestWalking_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &walking.WalkingRequest{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行步行路径规划请求
	resp, err := client.Walking(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestWalking_InvalidCoordinateFormat 测试Walking方法坐标格式错误
func TestWalking_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &walking.WalkingRequest{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行步行路径规划请求
	resp, err := client.Walking(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestLineTrafficStatus_Success 测试线路交通态势查询成功
func TestLineTrafficStatus_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"trafficinfo": {
			"description": "整体路况良好",
			"evaluation": {
				"expedite": 1,
				"congested": 0,
				"blocking": 0,
				"unknown": 0,
				"status": "expedite",
				"description": "畅通"
			},
			"roads": [{
				"name": "测试道路",
				"status": 1,
				"direction": "东向西",
				"lcodes": ["123456"],
				"polyline": "116.481028,39.989643;116.489028,39.999643",
				"speed": 60,
				"jams": []
			}]
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &line.LineTrafficRequest{
		Path: "116.481028,39.989643;116.489028,39.999643",
	}
	resp, err := client.LineTrafficStatus(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.Infocode)
	assert.NotNil(t, resp.Trafficinfo)
	assert.Equal(t, "整体路况良好", resp.Trafficinfo.Description)
	assert.Equal(t, 1, resp.Trafficinfo.Evaluation.Expedite)
	assert.Equal(t, 0, resp.Trafficinfo.Evaluation.Congested)
	assert.Equal(t, "expedite", resp.Trafficinfo.Evaluation.Status)
	assert.Len(t, resp.Trafficinfo.Roads, 1)
	assert.Equal(t, "测试道路", resp.Trafficinfo.Roads[0].Name)
	assert.Equal(t, 1, resp.Trafficinfo.Roads[0].Status)
	assert.Equal(t, 60.0, resp.Trafficinfo.Roads[0].Speed)
}

// TestCircleTrafficStatus_Success 测试圆形区域交通态势查询成功
func TestCircleTrafficStatus_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"trafficinfo": {
			"description": "整体路况良好",
			"evaluation": {
				"expedite": 2,
				"congested": 0,
				"blocking": 0,
				"unknown": 0,
				"status": "expedite",
				"description": "畅通"
			},
			"roads": [{
				"name": "测试道路1",
				"status": 1,
				"direction": "东向西",
				"lcodes": ["123456"],
				"polyline": "116.481028,39.989643;116.489028,39.999643",
				"speed": 60,
				"jams": []
			}, {
				"name": "测试道路2",
				"status": 1,
				"direction": "西向东",
				"lcodes": ["789012"],
				"polyline": "116.489028,39.999643;116.481028,39.989643",
				"speed": 55,
				"jams": []
			}]
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &circle.CircleTrafficRequest{
		Center: "116.481028,39.989643",
		Radius: "1000",
	}
	resp, err := client.CircleTrafficStatus(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.Infocode)
	assert.NotNil(t, resp.Trafficinfo)
	assert.Equal(t, "整体路况良好", resp.Trafficinfo.Description)
	assert.Equal(t, 2, resp.Trafficinfo.Evaluation.Expedite)
	assert.Equal(t, "expedite", resp.Trafficinfo.Evaluation.Status)
	assert.Len(t, resp.Trafficinfo.Roads, 2)
	assert.Equal(t, "测试道路1", resp.Trafficinfo.Roads[0].Name)
	assert.Equal(t, 1, resp.Trafficinfo.Roads[0].Status)
	assert.Equal(t, 60.0, resp.Trafficinfo.Roads[0].Speed)
	assert.Equal(t, "测试道路2", resp.Trafficinfo.Roads[1].Name)
	assert.Equal(t, 55.0, resp.Trafficinfo.Roads[1].Speed)
}

// TestRectangleTrafficStatus_Success 测试矩形区域交通态势查询成功
func TestRectangleTrafficStatus_Success(t *testing.T) {
	// 1. 创建mock服务器
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"trafficinfo": {
			"description": "整体路况良好",
			"evaluation": {
				"expedite": 3,
				"congested": 0,
				"blocking": 0,
				"unknown": 0,
				"status": "expedite",
				"description": "畅通"
			},
			"roads": [{
				"name": "测试道路1",
				"status": 1,
				"direction": "东向西",
				"lcodes": ["123456"],
				"polyline": "116.481028,39.989643;116.489028,39.999643",
				"speed": 60,
				"jams": []
			}, {
				"name": "测试道路2",
				"status": 1,
				"direction": "西向东",
				"lcodes": ["789012"],
				"polyline": "116.489028,39.999643;116.481028,39.989643",
				"speed": 55,
				"jams": []
			}, {
				"name": "测试道路3",
				"status": 1,
				"direction": "南向北",
				"lcodes": ["345678"],
				"polyline": "116.485028,39.989643;116.485028,39.999643",
				"speed": 50,
				"jams": []
			}]
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &rectangle.RectangleTrafficRequest{
		Rectangle: "116.481028,39.989643;116.489028,39.999643",
	}
	resp, err := client.RectangleTrafficStatus(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.Infocode)
	assert.NotNil(t, resp.Trafficinfo)
	assert.Equal(t, "整体路况良好", resp.Trafficinfo.Description)
	assert.Equal(t, 3, resp.Trafficinfo.Evaluation.Expedite)
	assert.Equal(t, "expedite", resp.Trafficinfo.Evaluation.Status)
	assert.Len(t, resp.Trafficinfo.Roads, 3)
	assert.Equal(t, "测试道路1", resp.Trafficinfo.Roads[0].Name)
	assert.Equal(t, 1, resp.Trafficinfo.Roads[0].Status)
	assert.Equal(t, 60.0, resp.Trafficinfo.Roads[0].Speed)
	assert.Equal(t, "测试道路2", resp.Trafficinfo.Roads[1].Name)
	assert.Equal(t, 55.0, resp.Trafficinfo.Roads[1].Speed)
	assert.Equal(t, "测试道路3", resp.Trafficinfo.Roads[2].Name)
	assert.Equal(t, 50.0, resp.Trafficinfo.Roads[2].Speed)
}

// TestWalking_APIError 测试Walking方法API返回错误
func TestWalking_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &walking.WalkingRequest{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行步行路径规划请求
	resp, err := client.Walking(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 驾车路径规划测试 --------------------------

// TestDriving_Success 测试Driving方法正常请求成功
func TestDriving_Success(t *testing.T) {
	// 1. 创建mock服务器，返回驾车路径规划成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "5237",
					"duration": "240",
					"steps": [
						{
							"instruction": "沿建国路向东行驶1.2公里",
							"orientation": "东",
							"road": "建国路",
							"distance": "1200",
							"duration": "60",
							"polyline": "116.351147,39.936871;116.352247,39.936771"
						}
					],
					"polyline": "116.351147,39.936871;116.352247,39.936771;116.410001,39.910113",
					"tolls": "5"
				}
			],
			"distance": "5237",
			"duration": "240",
			"tolls": "5"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &driving.DrivingRequest{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行驾车路径规划请求
	resp, err := client.Driving(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "5237", resp.Route.Distance)
	assert.Equal(t, "240", resp.Route.Duration)
	assert.Equal(t, "5", resp.Route.Tolls)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Equal(t, "5237", resp.Route.Paths[0].Distance)
	assert.Equal(t, "240", resp.Route.Paths[0].Duration)
	assert.Len(t, resp.Route.Paths[0].Steps, 1)
	assert.Equal(t, "沿建国路向东行驶1.2公里", resp.Route.Paths[0].Steps[0].Instruction)
}

// TestDriving_MissingOrigin 测试Driving方法缺少必填参数Origin
func TestDriving_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &driving.DrivingRequest{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行驾车路径规划请求
	resp, err := client.Driving(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestDriving_MissingDestination 测试Driving方法缺少必填参数Destination
func TestDriving_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &driving.DrivingRequest{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行驾车路径规划请求
	resp, err := client.Driving(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestDriving_InvalidCoordinateFormat 测试Driving方法坐标格式错误
func TestDriving_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &driving.DrivingRequest{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行驾车路径规划请求
	resp, err := client.Driving(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestDriving_APIError 测试Driving方法API返回错误
func TestDriving_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &driving.DrivingRequest{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行驾车路径规划请求
	resp, err := client.Driving(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 骑行路径规划测试 --------------------------

// TestBicycling_Success 测试Bicycling方法正常请求成功
func TestBicycling_Success(t *testing.T) {
	// 1. 创建mock服务器，返回骑行路径规划成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "4237",
					"duration": "280",
					"steps": [
						{
							"instruction": "沿骑行专用道向东行驶800米",
							"orientation": "东",
							"road": "建国路骑行道",
							"distance": "800",
							"duration": "50",
							"polyline": "116.351147,39.936871;116.351947,39.936771"
						}
					],
					"polyline": "116.351147,39.936871;116.351947,39.936771;116.410001,39.910113"
				}
			],
			"distance": "4237",
			"duration": "280"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &bicycling.BicyclingRequest{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行骑行路径规划请求
	resp, err := client.Bicycling(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "4237", resp.Route.Distance)
	assert.Equal(t, "280", resp.Route.Duration)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Equal(t, "4237", resp.Route.Paths[0].Distance)
	assert.Equal(t, "280", resp.Route.Paths[0].Duration)
	assert.Len(t, resp.Route.Paths[0].Steps, 1)
	assert.Equal(t, "沿骑行专用道向东行驶800米", resp.Route.Paths[0].Steps[0].Instruction)
}

// TestBicycling_MissingOrigin 测试Bicycling方法缺少必填参数Origin
func TestBicycling_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &bicycling.BicyclingRequest{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行骑行路径规划请求
	resp, err := client.Bicycling(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestBicycling_MissingDestination 测试Bicycling方法缺少必填参数Destination
func TestBicycling_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &bicycling.BicyclingRequest{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行骑行路径规划请求
	resp, err := client.Bicycling(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestBicycling_InvalidCoordinateFormat 测试Bicycling方法坐标格式错误
func TestBicycling_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &bicycling.BicyclingRequest{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行骑行路径规划请求
	resp, err := client.Bicycling(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestBicycling_APIError 测试Bicycling方法API返回错误
func TestBicycling_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &bicycling.BicyclingRequest{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行骑行路径规划请求
	resp, err := client.Bicycling(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- v2版本路径规划测试 --------------------------

// -------------------------- 步行路径规划v2测试 --------------------------

// TestWalkingV2_Success 测试WalkingV2方法正常请求成功
func TestWalkingV2_Success(t *testing.T) {
	// 1. 创建mock服务器，返回步行路径规划v2成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "3237",
					"duration": "324",
					"steps": [
						{
							"instruction": "步行100米，右转进入建国路",
							"orientation": "东",
							"road": "建国路",
							"distance": "100",
							"duration": "60",
							"polyline": "116.351147,39.936871;116.351247,39.936771",
							"action": "右转",
							"assistant_action": "",
							"walk_type": "1"
						}
					],
					"polyline": "116.351147,39.936871;116.351247,39.936771;116.410001,39.910113"
				}
			],
			"distance": "3237",
			"duration": "324"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &walkingV2.WalkingRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
		Strategy:    "0",
	}

	// 4. 执行步行路径规划v2请求
	resp, err := client.WalkingV2(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "3237", resp.Route.Distance)
	assert.Equal(t, "324", resp.Route.Duration)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Equal(t, "3237", resp.Route.Paths[0].Distance)
	assert.Equal(t, "324", resp.Route.Paths[0].Duration)
	assert.Len(t, resp.Route.Paths[0].Steps, 1)
	assert.Equal(t, "步行100米，右转进入建国路", resp.Route.Paths[0].Steps[0].Instruction)
	assert.Equal(t, "1", resp.Route.Paths[0].Steps[0].WalkType)
}

// TestWalkingV2_MissingOrigin 测试WalkingV2方法缺少必填参数Origin
func TestWalkingV2_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &walkingV2.WalkingRequestV2{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行步行路径规划v2请求
	resp, err := client.WalkingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestWalkingV2_MissingDestination 测试WalkingV2方法缺少必填参数Destination
func TestWalkingV2_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &walkingV2.WalkingRequestV2{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行步行路径规划v2请求
	resp, err := client.WalkingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestWalkingV2_InvalidCoordinateFormat 测试WalkingV2方法坐标格式错误
func TestWalkingV2_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &walkingV2.WalkingRequestV2{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行步行路径规划v2请求
	resp, err := client.WalkingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestWalkingV2_APIError 测试WalkingV2方法API返回错误
func TestWalkingV2_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &walkingV2.WalkingRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行步行路径规划v2请求
	resp, err := client.WalkingV2(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 驾车路径规划v2测试 --------------------------

// TestDrivingV2_Success 测试DrivingV2方法正常请求成功
func TestDrivingV2_Success(t *testing.T) {
	// 1. 创建mock服务器，返回驾车路径规划v2成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "5237",
					"duration": "240",
					"steps": [
						{
							"instruction": "沿建国路向东行驶1.2公里",
							"orientation": "东",
							"road": "建国路",
							"distance": "1200",
							"duration": "60",
							"polyline": "116.351147,39.936871;116.352247,39.936771"
						}
					],
					"polyline": "116.351147,39.936871;116.352247,39.936771;116.410001,39.910113",
					"tolls": "5"
				}
			],
			"distance": "5237",
			"duration": "240",
			"tolls": "5"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &drivingV2.DrivingRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行驾车路径规划v2请求
	resp, err := client.DrivingV2(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "5237", resp.Route.Distance)
	assert.Equal(t, "240", resp.Route.Duration)
	assert.Equal(t, "5", resp.Route.Tolls)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Equal(t, "5237", resp.Route.Paths[0].Distance)
	assert.Equal(t, "240", resp.Route.Paths[0].Duration)
	assert.Equal(t, "5", resp.Route.Paths[0].Tolls)
	assert.Len(t, resp.Route.Paths[0].Steps, 1)
	assert.Equal(t, "沿建国路向东行驶1.2公里", resp.Route.Paths[0].Steps[0].Instruction)
}

// TestDrivingV2_MissingOrigin 测试DrivingV2方法缺少必填参数Origin
func TestDrivingV2_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &drivingV2.DrivingRequestV2{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行驾车路径规划v2请求
	resp, err := client.DrivingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestDrivingV2_MissingDestination 测试DrivingV2方法缺少必填参数Destination
func TestDrivingV2_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &drivingV2.DrivingRequestV2{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行驾车路径规划v2请求
	resp, err := client.DrivingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestDrivingV2_InvalidCoordinateFormat 测试DrivingV2方法坐标格式错误
func TestDrivingV2_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &drivingV2.DrivingRequestV2{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行驾车路径规划v2请求
	resp, err := client.DrivingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestDrivingV2_APIError 测试DrivingV2方法API返回错误
func TestDrivingV2_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &drivingV2.DrivingRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行驾车路径规划v2请求
	resp, err := client.DrivingV2(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 未来驾车路径规划v4测试 --------------------------

// TestETDDrivingV4_Success 测试ETDDrivingV4方法正常请求成功
func TestETDDrivingV4_Success(t *testing.T) {
	// 1. 创建mock服务器，返回未来驾车路径规划v4成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "5237",
					"duration": "240",
					"steps": [
						{
							"instruction": "沿建国路向东行驶1.2公里",
							"orientation": "东",
							"road": "建国路",
							"distance": "1200",
							"duration": "60",
							"polyline": "116.351147,39.936871;116.352247,39.936771"
						}
					],
					"polyline": "116.351147,39.936871;116.352247,39.936771;116.410001,39.910113",
					"tolls": "5"
				}
			],
			"distance": "5237",
			"duration": "240",
			"tolls": "5",
			"etd": "2025-01-01 08:00",
			"eta": "2025-01-01 08:30"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &etdDrivingV4.ETDDrivingRequestV4{
		Origin:        "116.351147,39.936871",
		Destination:   "116.410001,39.910113",
		DepartureTime: "2025-01-01 08:00",
	}

	// 4. 执行未来驾车路径规划v4请求
	resp, err := client.ETDDrivingV4(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "5237", resp.Route.Distance)
	assert.Equal(t, "240", resp.Route.Duration)
	assert.Equal(t, "5", resp.Route.Tolls)
	assert.Equal(t, "2025-01-01 08:00", resp.Route.ETD)
	assert.Equal(t, "2025-01-01 08:30", resp.Route.ETA)
	assert.Len(t, resp.Route.Paths, 1)
}

// TestETDDrivingV4_MissingOrigin 测试ETDDrivingV4方法缺少必填参数Origin
func TestETDDrivingV4_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &etdDrivingV4.ETDDrivingRequestV4{
		Destination:   "116.410001,39.910113",
		DepartureTime: "2025-01-01 08:00",
	}

	// 3. 执行未来驾车路径规划v4请求
	resp, err := client.ETDDrivingV4(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestETDDrivingV4_MissingDestination 测试ETDDrivingV4方法缺少必填参数Destination
func TestETDDrivingV4_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &etdDrivingV4.ETDDrivingRequestV4{
		Origin:        "116.351147,39.936871",
		DepartureTime: "2025-01-01 08:00",
	}

	// 3. 执行未来驾车路径规划v4请求
	resp, err := client.ETDDrivingV4(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestETDDrivingV4_MissingDepartureTime 测试ETDDrivingV4方法缺少必填参数DepartureTime
func TestETDDrivingV4_MissingDepartureTime(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少DepartureTime的请求参数
	req := &etdDrivingV4.ETDDrivingRequestV4{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 3. 执行未来驾车路径规划v4请求
	resp, err := client.ETDDrivingV4(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "departure_time参数不能为空")
}

// -------------------------- 骑行路径规划v2测试 --------------------------

// TestBicyclingV2_Success 测试BicyclingV2方法正常请求成功
func TestBicyclingV2_Success(t *testing.T) {
	// 1. 创建mock服务器，返回骑行路径规划v2成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "4237",
					"duration": "280",
					"steps": [
						{
							"instruction": "沿骑行专用道向东行驶800米",
							"orientation": "东",
							"road": "建国路骑行道",
							"distance": "800",
							"duration": "50",
							"polyline": "116.351147,39.936871;116.351947,39.936771"
						}
					],
					"polyline": "116.351147,39.936871;116.351947,39.936771;116.410001,39.910113"
				}
			],
			"distance": "4237",
			"duration": "280"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &bicyclingV2.BicyclingRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行骑行路径规划v2请求
	resp, err := client.BicyclingV2(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "4237", resp.Route.Distance)
	assert.Equal(t, "280", resp.Route.Duration)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Equal(t, "4237", resp.Route.Paths[0].Distance)
	assert.Equal(t, "280", resp.Route.Paths[0].Duration)
	assert.Len(t, resp.Route.Paths[0].Steps, 1)
	assert.Equal(t, "沿骑行专用道向东行驶800米", resp.Route.Paths[0].Steps[0].Instruction)
}

// TestBicyclingV2_MissingOrigin 测试BicyclingV2方法缺少必填参数Origin
func TestBicyclingV2_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &bicyclingV2.BicyclingRequestV2{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行骑行路径规划v2请求
	resp, err := client.BicyclingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestBicyclingV2_MissingDestination 测试BicyclingV2方法缺少必填参数Destination
func TestBicyclingV2_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &bicyclingV2.BicyclingRequestV2{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行骑行路径规划v2请求
	resp, err := client.BicyclingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestBicyclingV2_InvalidCoordinateFormat 测试BicyclingV2方法坐标格式错误
func TestBicyclingV2_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &bicyclingV2.BicyclingRequestV2{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行骑行路径规划v2请求
	resp, err := client.BicyclingV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestBicyclingV2_APIError 测试BicyclingV2方法API返回错误
func TestBicyclingV2_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &bicyclingV2.BicyclingRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行骑行路径规划v2请求
	resp, err := client.BicyclingV2(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 公交路径规划v2测试 --------------------------

// TestBusV2_Success 测试BusV2方法正常请求成功
func TestBusV2_Success(t *testing.T) {
	// 1. 创建mock服务器，返回公交路径规划v2成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "6237",
					"duration": "480",
					"transits": [
						{
							"distance": "6237",
							"duration": "480",
							"cost": "2",
							"segments": [
								{
									"walking": {
										"distance": "300",
										"duration": "30",
										"instruction": "步行300米至建国门站"
									}
								}
							]
						}
					]
				}
			]
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &busV2.BusRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行公交路径规划v2请求
	resp, err := client.BusV2(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.NotNil(t, resp.Route)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Len(t, resp.Route.Paths[0].Transits, 1)
	assert.Equal(t, "6237", resp.Route.Paths[0].Distance)
	assert.Equal(t, "480", resp.Route.Paths[0].Duration)
}

// TestBusV2_MissingOrigin 测试BusV2方法缺少必填参数Origin
func TestBusV2_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &busV2.BusRequestV2{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行公交路径规划v2请求
	resp, err := client.BusV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestBusV2_MissingDestination 测试BusV2方法缺少必填参数Destination
func TestBusV2_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &busV2.BusRequestV2{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行公交路径规划v2请求
	resp, err := client.BusV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestBusV2_InvalidCoordinateFormat 测试BusV2方法坐标格式错误
func TestBusV2_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &busV2.BusRequestV2{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行公交路径规划v2请求
	resp, err := client.BusV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestBusV2_APIError 测试BusV2方法API返回错误
func TestBusV2_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &busV2.BusRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行公交路径规划v2请求
	resp, err := client.BusV2(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 电动车路径规划v2测试 --------------------------

// TestElectricV2_Success 测试ElectricV2方法正常请求成功
func TestElectricV2_Success(t *testing.T) {
	// 1. 创建mock服务器，返回电动车路径规划v2成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"route": {
			"origin": "116.351147,39.936871",
			"destination": "116.410001,39.910113",
			"paths": [
				{
					"distance": "5537",
					"duration": "320",
					"steps": [
						{
							"instruction": "沿电动车专用道向东行驶1.5公里",
							"orientation": "东",
							"road": "建国路电动车道",
							"distance": "1500",
							"duration": "75",
							"polyline": "116.351147,39.936871;116.352647,39.936771"
						}
					],
					"polyline": "116.351147,39.936871;116.352647,39.936771;116.410001,39.910113",
					"tolls": "0",
					"charge_info": {
						"battery_usage": "2.5",
						"charge_stations": [
							{
								"name": "建国门充电站",
								"location": "116.380000,39.920000",
								"distance": "2500"
							}
						],
						"total_charge_fee": "15"
					}
				}
			],
			"distance": "5537",
			"duration": "320",
			"tolls": "0"
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &electricV2.ElectricRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行电动车路径规划v2请求
	resp, err := client.ElectricV2(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "116.351147,39.936871", resp.Route.Origin)
	assert.Equal(t, "116.410001,39.910113", resp.Route.Destination)
	assert.Equal(t, "5537", resp.Route.Distance)
	assert.Equal(t, "320", resp.Route.Duration)
	assert.Equal(t, "0", resp.Route.Tolls)
	assert.Len(t, resp.Route.Paths, 1)
	assert.Equal(t, "5537", resp.Route.Paths[0].Distance)
	assert.Equal(t, "320", resp.Route.Paths[0].Duration)
	assert.Equal(t, "0", resp.Route.Paths[0].Tolls)
	assert.Len(t, resp.Route.Paths[0].Steps, 1)
	assert.Equal(t, "沿电动车专用道向东行驶1.5公里", resp.Route.Paths[0].Steps[0].Instruction)
	assert.Len(t, resp.Route.Paths[0].ChargeInfo.ChargeStations, 1)
	assert.Equal(t, "建国门充电站", resp.Route.Paths[0].ChargeInfo.ChargeStations[0].Name)
}

// TestElectricV2_MissingOrigin 测试ElectricV2方法缺少必填参数Origin
func TestElectricV2_MissingOrigin(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Origin的请求参数
	req := &electricV2.ElectricRequestV2{
		Destination: "116.410001,39.910113",
	}

	// 3. 执行电动车路径规划v2请求
	resp, err := client.ElectricV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "origin参数不能为空")
}

// TestElectricV2_MissingDestination 测试ElectricV2方法缺少必填参数Destination
func TestElectricV2_MissingDestination(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Destination的请求参数
	req := &electricV2.ElectricRequestV2{
		Origin: "116.351147,39.936871",
	}

	// 3. 执行电动车路径规划v2请求
	resp, err := client.ElectricV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "destination参数不能为空")
}

// TestElectricV2_InvalidCoordinateFormat 测试ElectricV2方法坐标格式错误
func TestElectricV2_InvalidCoordinateFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的坐标参数
	req := &electricV2.ElectricRequestV2{
		Origin:      "116.351147 39.936871", // 使用空格分隔而不是逗号
		Destination: "116.410001,39.910113",
	}

	// 3. 执行电动车路径规划v2请求
	resp, err := client.ElectricV2(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "坐标格式错误")
}

// TestElectricV2_APIError 测试ElectricV2方法API返回错误
func TestElectricV2_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &electricV2.ElectricRequestV2{
		Origin:      "116.351147,39.936871",
		Destination: "116.410001,39.910113",
	}

	// 4. 执行电动车路径规划v2请求
	resp, err := client.ElectricV2(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 行政区查询测试 --------------------------

// TestDistrict_Success 测试District方法正常请求成功
func TestDistrict_Success(t *testing.T) {
	// 1. 创建mock服务器，返回行政区查询成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"districts": [
			{
				"name": "北京市",
				"level": "province",
				"adcode": "110000",
				"citycode": "110000",
				"center": "116.405285,39.904989",
				"districts": [
					{
						"name": "东城区",
						"level": "district",
						"adcode": "110101",
						"citycode": "110000",
						"center": "116.410708,39.915224",
						"parent_city": ["北京市"]
					},
					{
						"name": "西城区",
						"level": "district",
						"adcode": "110102",
						"citycode": "110000",
						"center": "116.363593,39.913362",
						"parent_city": ["北京市"]
					}
				]
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &district.DistrictRequest{
		Keywords:    "北京市",
		Subdistrict: "1", // 返回下一级行政区
	}

	// 4. 执行行政区查询请求
	resp, err := client.District(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "1", resp.Count)
	assert.Len(t, resp.Districts, 1)
	assert.Equal(t, "北京市", resp.Districts[0].Name)
	assert.Equal(t, "province", resp.Districts[0].Level)
	assert.Equal(t, "110000", resp.Districts[0].Adcode)
	assert.Equal(t, "116.405285,39.904989", resp.Districts[0].Center)
	assert.Len(t, resp.Districts[0].Districts, 2) // 应该返回2个区县
	assert.Equal(t, "东城区", resp.Districts[0].Districts[0].Name)
	assert.Equal(t, "西城区", resp.Districts[0].Districts[1].Name)
}

// TestDistrict_MissingKeywords 测试District方法缺少必填参数Keywords
func TestDistrict_MissingKeywords(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Keywords的请求参数
	req := &district.DistrictRequest{
		Subdistrict: "1", // 只有Subdistrict，没有Keywords
	}

	// 3. 执行行政区查询请求
	resp, err := client.District(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "keywords参数不能为空")
}

// TestDistrict_APIError 测试District方法API返回错误
func TestDistrict_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &district.DistrictRequest{
		Keywords: "北京市",
	}

	// 4. 执行行政区查询请求
	resp, err := client.District(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestDistrict_WithAdcode 测试通过adcode查询行政区
func TestDistrict_WithAdcode(t *testing.T) {
	// 1. 创建mock服务器，返回行政区查询成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"districts": [
			{
				"name": "深圳市",
				"level": "city",
				"adcode": "440300",
				"citycode": "0755",
				"center": "114.057868,22.543099"
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数，使用adcode查询
	req := &district.DistrictRequest{
		Keywords:    "440300", // 深圳的adcode
		Subdistrict: "0",      // 不返回子级行政区
	}

	// 4. 执行行政区查询请求
	resp, err := client.District(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "1", resp.Count)
	assert.Len(t, resp.Districts, 1)
	assert.Equal(t, "深圳市", resp.Districts[0].Name)
	assert.Equal(t, "city", resp.Districts[0].Level)
	assert.Equal(t, "440300", resp.Districts[0].Adcode)
	assert.Equal(t, "0755", resp.Districts[0].Citycode)
	assert.Len(t, resp.Districts[0].Districts, 0) // 不返回子级行政区
}

// TestDistrict_WithFilter 测试使用筛选条件查询行政区
func TestDistrict_WithFilter(t *testing.T) {
	// 1. 创建mock服务器，返回行政区查询成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"districts": [
			{
				"name": "广州市",
				"level": "city",
				"adcode": "440100",
				"citycode": "020",
				"center": "113.280637,23.125178"
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数，使用筛选条件
	req := &district.DistrictRequest{
		Keywords: "广州",
		Filter:   "citycode:020", // 筛选citycode为020的城市
	}

	// 4. 执行行政区查询请求
	resp, err := client.District(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "1", resp.Count)
	assert.Len(t, resp.Districts, 1)
	assert.Equal(t, "广州市", resp.Districts[0].Name)
	assert.Equal(t, "020", resp.Districts[0].Citycode)
}

// -------------------------- 交通事件查询测试 --------------------------

// TestTrafficIncident_Success 测试TrafficIncident方法正常请求成功
func TestTrafficIncident_Success(t *testing.T) {
	// 1. 创建mock服务器，返回交通事件查询成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"trafficincidents": [
			{
				"id": "12345",
				"location": "116.351147,39.904989",
				"type": "1",
				"type_des": "道路施工",
				"level": "2",
				"level_des": "一般",
				"description": "建国路与长安街交叉口道路施工，影响车辆通行",
				"polyline": "116.351147,39.904989;116.352147,39.905989",
				"road": "建国路",
				"start_time": "1600000000",
				"end_time": "1600003600",
				"direction": "东向西",
				"status": "1",
				"impact_level": "2",
				"affect_road_length": "500",
				"first_report_time": "1600000000",
				"last_report_time": "1600000100"
			},
			{
				"id": "12346",
				"location": "116.361147,39.914989",
				"type": "3",
				"type_des": "交通事故",
				"level": "3",
				"level_des": "严重",
				"description": "长安街东单路口发生交通事故，占用2条车道",
				"polyline": "116.361147,39.914989;116.362147,39.915989",
				"road": "长安街",
				"start_time": "1600000000",
				"end_time": "1600002400",
				"direction": "西向东",
				"status": "1",
				"impact_level": "3",
				"affect_road_length": "800",
				"first_report_time": "1600000000",
				"last_report_time": "1600000200"
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &trafficIncident.TrafficIncidentRequest{
		Level:     "1",                                         // 所有级别
		Type:      "1|3",                                       // 道路施工和交通事故
		Rectangle: "116.300000,39.900000,116.400000,39.950000", // 北京核心区域
	}

	// 4. 执行交通事件查询请求
	resp, err := client.TrafficIncident(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Len(t, resp.Trafficincidents, 2)

	// 验证第一个事件
	assert.Equal(t, "12345", resp.Trafficincidents[0].Id)
	assert.Equal(t, "1", resp.Trafficincidents[0].Type)
	assert.Equal(t, "道路施工", resp.Trafficincidents[0].TypeDes)
	assert.Equal(t, "2", resp.Trafficincidents[0].Level)
	assert.Equal(t, "一般", resp.Trafficincidents[0].LevelDes)
	assert.Equal(t, "建国路", resp.Trafficincidents[0].Road)
	assert.Equal(t, "1", resp.Trafficincidents[0].Status)

	// 验证第二个事件
	assert.Equal(t, "12346", resp.Trafficincidents[1].Id)
	assert.Equal(t, "3", resp.Trafficincidents[1].Type)
	assert.Equal(t, "交通事故", resp.Trafficincidents[1].TypeDes)
	assert.Equal(t, "3", resp.Trafficincidents[1].Level)
	assert.Equal(t, "严重", resp.Trafficincidents[1].LevelDes)
	assert.Equal(t, "长安街", resp.Trafficincidents[1].Road)
}

// TestTrafficIncident_MissingLevel 测试TrafficIncident方法缺少必填参数Level
func TestTrafficIncident_MissingLevel(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Level的请求参数
	req := &trafficIncident.TrafficIncidentRequest{
		Type:      "1",
		Rectangle: "116.300000,39.900000,116.400000,39.950000",
	}

	// 3. 执行交通事件查询请求
	resp, err := client.TrafficIncident(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "level参数不能为空")
}

// TestTrafficIncident_MissingType 测试TrafficIncident方法缺少必填参数Type
func TestTrafficIncident_MissingType(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Type的请求参数
	req := &trafficIncident.TrafficIncidentRequest{
		Level:     "1",
		Rectangle: "116.300000,39.900000,116.400000,39.950000",
	}

	// 3. 执行交通事件查询请求
	resp, err := client.TrafficIncident(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "type参数不能为空")
}

// TestTrafficIncident_MissingRectangle 测试TrafficIncident方法缺少必填参数Rectangle
func TestTrafficIncident_MissingRectangle(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少Rectangle的请求参数
	req := &trafficIncident.TrafficIncidentRequest{
		Level: "1",
		Type:  "1",
	}

	// 3. 执行交通事件查询请求
	resp, err := client.TrafficIncident(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "rectangle参数不能为空")
}

// TestTrafficIncident_InvalidRectangleFormat 测试TrafficIncident方法Rectangle格式错误
func TestTrafficIncident_InvalidRectangleFormat(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建格式错误的Rectangle参数
	req := &trafficIncident.TrafficIncidentRequest{
		Level:     "1",
		Type:      "1",
		Rectangle: "116.300000,39.900000,116.400000", // 缺少一个坐标点
	}

	// 3. 执行交通事件查询请求
	resp, err := client.TrafficIncident(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "rectangle格式错误")
}

// TestTrafficIncident_APIError 测试TrafficIncident方法API返回错误
func TestTrafficIncident_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &trafficIncident.TrafficIncidentRequest{
		Level:     "1",
		Type:      "1",
		Rectangle: "116.300000,39.900000,116.400000,39.950000",
	}

	// 4. 执行交通事件查询请求
	resp, err := client.TrafficIncident(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestTrafficIncident_WithExtensionsAll 测试TrafficIncident方法使用extensions=all
func TestTrafficIncident_WithExtensionsAll(t *testing.T) {
	// 1. 创建mock服务器，返回带详细信息的交通事件响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"trafficincidents": [
			{
				"id": "12347",
				"location": "116.371147,39.924989",
				"type": "5",
				"type_des": "交通拥堵",
				"level": "3",
				"level_des": "严重",
				"description": "东三环中路交通拥堵，车辆行驶缓慢",
				"polyline": "116.371147,39.924989;116.372147,39.925989",
				"road": "东三环中路",
				"start_time": "1600000000",
				"end_time": "1600003600",
				"direction": "南向北",
				"status": "1",
				"impact_level": "3",
				"affect_road_length": "1000",
				"jams": [
					{
						"location": "116.371147,39.924989",
						"direction": "南向北",
						"length": "1000",
						"level": "4",
						"status": "1",
						"speed": "10",
						"polyline": "116.371147,39.924989;116.372147,39.925989",
						"start_time": "1600000000",
						"end_time": "1600003600"
					}
				],
				"first_report_time": "1600000000",
				"last_report_time": "1600000100"
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数，使用extensions=all
	req := &trafficIncident.TrafficIncidentRequest{
		Level:      "1",
		Type:       "5",
		Rectangle:  "116.300000,39.900000,116.400000,39.950000",
		Extensions: "all", // 返回详细信息
	}

	// 4. 执行交通事件查询请求
	resp, err := client.TrafficIncident(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Trafficincidents, 1)
	assert.Len(t, resp.Trafficincidents[0].Jams, 1) // 包含拥堵信息
	assert.Equal(t, "1000", resp.Trafficincidents[0].Jams[0].Length)
	assert.Equal(t, "4", resp.Trafficincidents[0].Jams[0].Level)
	assert.Equal(t, "10", resp.Trafficincidents[0].Jams[0].Speed)
}

// -------------------------- IP定位测试 --------------------------

// TestIPConfig_Success 测试IPConfig方法正常请求成功
func TestIPConfig_Success(t *testing.T) {
	// 1. 创建mock服务器，返回IP定位成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"ip": "114.114.114.114",
		"country": "中国",
		"province": "江苏省",
		"city": "南京市",
		"district": "秦淮区",
		"adcode": "320104",
		"center": "118.796877,32.048458",
		"isp": "江苏省电信"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &ipV3.IPConfigRequest{
		IP: "114.114.114.114",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPConfig(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "114.114.114.114", resp.IP)
	assert.Equal(t, "中国", resp.Country)
	assert.Equal(t, "江苏省", resp.Province)
	assert.Equal(t, "南京市", resp.City)
	assert.Equal(t, "秦淮区", resp.District)
	assert.Equal(t, "320104", resp.Adcode)
	assert.Equal(t, "118.796877,32.048458", resp.Center)
	assert.Equal(t, "江苏省电信", resp.ISP)
}

// TestIPConfig_MissingIP 测试IPConfig方法缺少必填参数IP
func TestIPConfig_MissingIP(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少IP的请求参数
	req := &ipV3.IPConfigRequest{
		// IP参数为空
	}

	// 3. 执行IP定位请求
	resp, err := client.IPConfig(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "ip参数不能为空")
}

// TestIPConfig_APIError 测试IPConfig方法API返回错误
func TestIPConfig_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &ipV3.IPConfigRequest{
		IP: "114.114.114.114",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPConfig(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestIPConfig_IPv6 测试IPConfig方法使用IPv6地址
func TestIPConfig_IPv6(t *testing.T) {
	// 1. 创建mock服务器，返回IPv6地址定位响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"ip": "2001:4860:4860::8888",
		"country": "美国",
		"province": "加利福尼亚州",
		"city": "山景城",
		"district": "圣克拉拉县",
		"adcode": "84006085",
		"center": "-122.083851,37.422258",
		"isp": "Google LLC"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数，使用IPv6地址
	req := &ipV3.IPConfigRequest{
		IP: "2001:4860:4860::8888",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPConfig(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "2001:4860:4860::8888", resp.IP)
	assert.Equal(t, "美国", resp.Country)
	assert.Equal(t, "加利福尼亚州", resp.Province)
	assert.Equal(t, "山景城", resp.City)
	assert.Equal(t, "圣克拉拉县", resp.District)
	assert.Equal(t, "84006085", resp.Adcode)
	assert.Equal(t, "-122.083851,37.422258", resp.Center)
	assert.Equal(t, "Google LLC", resp.ISP)
}

// TestIPConfig_WithLocationInfo 测试IPConfig方法返回详细位置信息
func TestIPConfig_WithLocationInfo(t *testing.T) {
	// 1. 创建mock服务器，返回带详细位置信息的IP定位响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"ip": "8.8.8.8",
		"country": "美国",
		"province": "弗吉尼亚州",
		"city": "阿什本",
		"district": "劳登县",
		"adcode": "84051107",
		"center": "-77.438217,39.032222",
		"isp": "Google LLC",
		"location": {
			"lat": "39.032222",
			"lon": "-77.438217",
			"address": "美国弗吉尼亚州阿什本劳登县",
			"city_code": "840511070000",
			"province_code": "840510000000",
			"district_code": "840511070000",
			"isp_info": {
				"name": "Google LLC",
				"type": "ISP",
				"mcc": "310",
				"mnc": "260"
			}
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &ipV3.IPConfigRequest{
		IP: "8.8.8.8",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPConfig(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "8.8.8.8", resp.IP)
	assert.NotNil(t, resp.Location)
	assert.Equal(t, "39.032222", resp.Location.Lat)
	assert.Equal(t, "-77.438217", resp.Location.Lon)
	assert.Equal(t, "美国弗吉尼亚州阿什本劳登县", resp.Location.Address)
	assert.NotNil(t, resp.Location.ISPInfo)
	assert.Equal(t, "Google LLC", resp.Location.ISPInfo.Name)
	assert.Equal(t, "ISP", resp.Location.ISPInfo.Type)
	assert.Equal(t, "310", resp.Location.ISPInfo.MCC)
	assert.Equal(t, "260", resp.Location.ISPInfo.MNC)
}

// -------------------------- IP v5 测试 --------------------------

// TestIPV5Config_Success 测试IPV5Config方法正常请求成功
func TestIPV5Config_Success(t *testing.T) {
	// 1. 创建mock服务器，返回正常IP定位响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"ip": "114.114.114.114",
		"country": "中国",
		"province": "江苏省",
		"city": "南京市",
		"district": "秦淮区",
		"adcode": "320104",
		"center": "118.796877,32.048458",
		"isp": "江苏省电信"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &ipV5.IPConfigRequest{
		IP: "114.114.114.114",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPV5Config(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "114.114.114.114", resp.IP)
	assert.Equal(t, "中国", resp.Country)
	assert.Equal(t, "江苏省", resp.Province)
	assert.Equal(t, "南京市", resp.City)
	assert.Equal(t, "秦淮区", resp.District)
	assert.Equal(t, "320104", resp.Adcode)
	assert.Equal(t, "118.796877,32.048458", resp.Center)
	assert.Equal(t, "江苏省电信", resp.ISP)
}

// TestIPV5Config_MissingIP 测试IPV5Config方法缺少必填参数IP
func TestIPV5Config_MissingIP(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少IP的请求参数
	req := &ipV5.IPConfigRequest{
		// IP参数为空
	}

	// 3. 执行IP定位请求
	resp, err := client.IPV5Config(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "ip参数不能为空")
}

// TestIPV5Config_APIError 测试IPV5Config方法API返回错误
func TestIPV5Config_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &ipV5.IPConfigRequest{
		IP: "114.114.114.114",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPV5Config(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestIPV5Config_IPv6 测试IPV5Config方法使用IPv6地址
func TestIPV5Config_IPv6(t *testing.T) {
	// 1. 创建mock服务器，返回IPv6地址定位响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"ip": "2001:4860:4860::8888",
		"country": "美国",
		"province": "加利福尼亚州",
		"city": "山景城",
		"district": "圣克拉拉县",
		"adcode": "84006085",
		"center": "-122.083851,37.422258",
		"isp": "Google LLC"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数，使用IPv6地址
	req := &ipV5.IPConfigRequest{
		IP: "2001:4860:4860::8888",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPV5Config(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "2001:4860:4860::8888", resp.IP)
	assert.Equal(t, "美国", resp.Country)
	assert.Equal(t, "加利福尼亚州", resp.Province)
	assert.Equal(t, "山景城", resp.City)
	assert.Equal(t, "圣克拉拉县", resp.District)
	assert.Equal(t, "84006085", resp.Adcode)
	assert.Equal(t, "-122.083851,37.422258", resp.Center)
	assert.Equal(t, "Google LLC", resp.ISP)
}

// TestIPV5Config_WithLocationInfo 测试IPV5Config方法返回详细位置信息
func TestIPV5Config_WithLocationInfo(t *testing.T) {
	// 1. 创建mock服务器，返回带详细位置信息的IP定位响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"ip": "8.8.8.8",
		"country": "美国",
		"province": "弗吉尼亚州",
		"city": "阿什本",
		"district": "劳登县",
		"adcode": "84051107",
		"center": "-77.438217,39.032222",
		"isp": "Google LLC",
		"location": {
			"lat": "39.032222",
			"lon": "-77.438217",
			"address": "美国弗吉尼亚州阿什本劳登县",
			"city_code": "ASB",
			"province_code": "VA",
			"district_code": "LDN",
			"isp_info": {
				"name": "Google LLC",
				"type": "ISP",
				"mcc": "310",
				"mnc": "260"
			}
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建请求参数
	req := &ipV5.IPConfigRequest{
		IP: "8.8.8.8",
	}

	// 4. 执行IP定位请求
	resp, err := client.IPV5Config(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "8.8.8.8", resp.IP)
	assert.Equal(t, "美国", resp.Country)
	assert.Equal(t, "弗吉尼亚州", resp.Province)
	assert.Equal(t, "阿什本", resp.City)
	assert.Equal(t, "劳登县", resp.District)
	assert.Equal(t, "84051107", resp.Adcode)
	assert.Equal(t, "-77.438217,39.032222", resp.Center)
	assert.Equal(t, "Google LLC", resp.ISP)
	assert.NotNil(t, resp.Location)
	assert.Equal(t, "39.032222", resp.Location.Lat)
	assert.Equal(t, "-77.438217", resp.Location.Lon)
	assert.Equal(t, "美国弗吉尼亚州阿什本劳登县", resp.Location.Address)
	assert.Equal(t, "ASB", resp.Location.CityCode)
	assert.Equal(t, "VA", resp.Location.ProvinceCode)
	assert.Equal(t, "LDN", resp.Location.DistrictCode)
	assert.NotNil(t, resp.Location.ISPInfo)
	assert.Equal(t, "Google LLC", resp.Location.ISPInfo.Name)
	assert.Equal(t, "ISP", resp.Location.ISPInfo.Type)
	assert.Equal(t, "310", resp.Location.ISPInfo.MCC)
	assert.Equal(t, "260", resp.Location.ISPInfo.MNC)
}

// -------------------------- 坐标转换测试 --------------------------

// TestConvert_Success 测试Convert方法正常转换成功
func TestConvert_Success(t *testing.T) {
	// 1. 创建mock服务器，返回正常转换结果
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"locations": "116.480656,39.989610;116.30815,39.95965"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求，转换GPS坐标到高德坐标
	req := &convert.ConvertRequest{
		Locations: "116.480656,39.989610;116.30815,39.95965",
		CoordSys:  "gps",
	}
	resp, err := client.Convert(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "116.480656,39.989610;116.30815,39.95965", resp.Locations)
}

// TestConvert_MissingLocations 测试Convert方法缺少必填参数locations
func TestConvert_MissingLocations(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 执行请求，缺少locations参数
	req := &convert.ConvertRequest{
		Locations: "", // 缺少必填参数
		CoordSys:  "gps",
	}
	resp, err := client.Convert(req)

	// 3. 验证结果
	assert.Error(t, err)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Nil(t, resp)
}

// TestConvert_MissingCoordSys 测试Convert方法缺少必填参数coordsys
func TestConvert_MissingCoordSys(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 执行请求，缺少coordsys参数
	req := &convert.ConvertRequest{
		Locations: "116.480656,39.989610",
		CoordSys:  "", // 缺少必填参数
	}
	resp, err := client.Convert(req)

	// 3. 验证结果
	assert.Error(t, err)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Nil(t, resp)
}

// TestConvert_APIError 测试Convert方法API返回错误
func TestConvert_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的coordsys参数",
		"infocode": "10003"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求，使用无效的coordsys
	req := &convert.ConvertRequest{
		Locations: "116.480656,39.989610",
		CoordSys:  "invalid_coordsys", // 无效的coordsys值
	}
	resp, err := client.Convert(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10003", apiErr.Code)
	assert.Equal(t, "无效的coordsys参数", apiErr.Info)
}

// -------------------------- 轨迹纠偏测试 --------------------------

// TestGraspRoad_Success 测试GraspRoad方法正常纠偏成功
func TestGraspRoad_Success(t *testing.T) {
	// 1. 创建mock服务器，返回正常纠偏结果
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"sid": "test_sid",
		"paths": [
			{
				"points": [
					{
						"location": "116.480656,39.989610",
						"time": 1600000000,
						"speed": 30.5
					},
					{
						"location": "116.481656,39.990610",
						"time": 1600000100,
						"speed": 32.0
					}
				]
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求，纠偏轨迹点
	req := &grasproad.GraspRoadRequest{
		SID:            "test_sid",
		Points:         "116.480656,39.989610,1600000000,30.5;116.481656,39.990610,1600000100,32.0",
		CoordTypeInput: "gps",
	}
	resp, err := client.GraspRoad(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test_sid", resp.SID)
	assert.Len(t, resp.Paths, 1)
	assert.Len(t, resp.Paths[0].Points, 2)
	assert.Equal(t, "116.480656,39.989610", resp.Paths[0].Points[0].Location)
	assert.Equal(t, int64(1600000000), resp.Paths[0].Points[0].Time)
	assert.Equal(t, 30.5, resp.Paths[0].Points[0].Speed)
}

// TestGraspRoad_MissingSID 测试GraspRoad方法缺少必填参数sid
func TestGraspRoad_MissingSID(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 执行请求，缺少sid参数
	req := &grasproad.GraspRoadRequest{
		Points:         "116.480656,39.989610,1600000000,30.5",
		CoordTypeInput: "gps",
	}
	resp, err := client.GraspRoad(req)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "sid参数不能为空")
}

// TestGraspRoad_MissingPoints 测试GraspRoad方法缺少必填参数points
func TestGraspRoad_MissingPoints(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 执行请求，缺少points参数
	req := &grasproad.GraspRoadRequest{
		SID:            "test_sid",
		CoordTypeInput: "gps",
	}
	resp, err := client.GraspRoad(req)

	// 3. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "points参数不能为空")
}

// TestGraspRoad_APIError 测试GraspRoad方法API返回错误
func TestGraspRoad_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	req := &grasproad.GraspRoadRequest{
		SID:            "test_sid",
		Points:         "116.480656,39.989610,1600000000,30.5",
		CoordTypeInput: "gps",
	}
	resp, err := client.GraspRoad(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestGraspRoad_WithExtensionsAll 测试GraspRoad方法使用extensions=all
func TestGraspRoad_WithExtensionsAll(t *testing.T) {
	// 1. 创建mock服务器，返回带详细信息的纠偏结果
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"sid": "test_sid_all",
		"paths": [
			{
				"distance": 200,
				"time": 20,
				"points": [
					{
						"location": "116.480656,39.989610",
						"time": 1600000000,
						"speed": 30.5,
						"direction": 90,
						"road_id": "road123",
						"road_name": "建国路",
						"match_type": 1,
						"status": 0
					},
					{
						"location": "116.481656,39.990610",
						"time": 1600000100,
						"speed": 32.0,
						"direction": 95,
						"road_id": "road123",
						"road_name": "建国路",
						"match_type": 1,
						"status": 0
					}
				],
				"steps": [
					{
						"start_index": 0,
						"end_index": 1,
						"road": {
							"id": "road123",
							"name": "建国路",
							"type": 4,
							"level": 3,
							"width": 30.5,
							"lanes": 4,
							"max_speed": 60
						}
					}
				]
			}
		]
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求，使用extensions=all
	req := &grasproad.GraspRoadRequest{
		SID:            "test_sid_all",
		Points:         "116.480656,39.989610,1600000000,30.5;116.481656,39.990610,1600000100,32.0",
		CoordTypeInput: "gps",
		Extensions:     "all",
	}
	resp, err := client.GraspRoad(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test_sid_all", resp.SID)
	assert.Len(t, resp.Paths, 1)
	assert.Equal(t, 200, resp.Paths[0].Distance)
	assert.Equal(t, 20, resp.Paths[0].Time)
	assert.Len(t, resp.Paths[0].Points, 2)
	assert.Equal(t, 90, resp.Paths[0].Points[0].Direction)
	assert.Equal(t, "road123", resp.Paths[0].Points[0].RoadID)
	assert.Equal(t, "建国路", resp.Paths[0].Points[0].RoadName)
	assert.Equal(t, 1, resp.Paths[0].Points[0].MatchType)
	assert.Equal(t, 0, resp.Paths[0].Points[0].Status)
	assert.Len(t, resp.Paths[0].Steps, 1)
	assert.Equal(t, "建国路", resp.Paths[0].Steps[0].Road.Name)
	assert.Equal(t, 60, resp.Paths[0].Steps[0].Road.MaxSpeed)
}

// -------------------------- 输入提示测试 --------------------------

// TestInputtips_Success 测试Inputtips方法正常请求成功
func TestInputtips_Success(t *testing.T) {
	// 1. 创建mock服务器，返回输入提示成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "2",
		"tips": [
			{
				"id": "B000A83M61",
				"name": "北京市朝阳区望京SOHO",
				"district": "朝阳区",
				"adcode": "110105",
				"location": "116.48693,39.99936",
				"address": "望京街8号",
				"type": "商务住宅;楼宇;商务写字楼",
				"typecode": "120201",
				"weight": "90",
				"city": "北京市",
				"citycode": "010",
				"districtadcode": "110105",
				"province": "北京市",
				"business_area": "望京"
			},
			{
				"id": "B000A83M62",
				"name": "北京市朝阳区望京公园",
				"district": "朝阳区",
				"adcode": "110105",
				"location": "116.47693,39.98936",
				"address": "望京西路",
				"type": "风景名胜;公园广场;公园",
				"typecode": "060301",
				"weight": "85",
				"city": "北京市",
				"citycode": "010",
				"districtadcode": "110105",
				"province": "北京市",
				"business_area": "望京"
			}
		],
		"suggestion": {
			"keywords": ["望京SOHO", "望京公园", "望京医院"],
			"cities": ["北京市", "上海市"]
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建输入提示请求
	req := &inputtips.InputtipsRequest{
		Keywords: "望京",
		City:     "北京",
		Datatype: "all",
	}

	// 4. 执行输入提示请求
	resp, err := client.Inputtips(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "2", resp.Count)
	assert.Len(t, resp.Tips, 2)
	assert.Equal(t, "B000A83M61", resp.Tips[0].ID)
	assert.Equal(t, "北京市朝阳区望京SOHO", resp.Tips[0].Name)
	assert.Equal(t, "朝阳区", resp.Tips[0].District)
	assert.Equal(t, "110105", resp.Tips[0].Adcode)
	assert.Equal(t, "116.48693,39.99936", resp.Tips[0].Location)
	assert.Equal(t, "商务住宅;楼宇;商务写字楼", resp.Tips[0].Type)
	assert.Equal(t, "望京", resp.Tips[0].BusinessArea)
	assert.Len(t, resp.Suggestion.Keywords, 3)
	assert.Contains(t, resp.Suggestion.Keywords, "望京SOHO")
	assert.Len(t, resp.Suggestion.Cities, 2)
	assert.Contains(t, resp.Suggestion.Cities, "北京市")
}

// TestInputtips_MissingKeywords 测试Inputtips方法缺少必填参数keywords
func TestInputtips_MissingKeywords(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少关键字的请求
	req := &inputtips.InputtipsRequest{
		City: "北京", // 缺少keywords
	}

	// 3. 执行输入提示请求
	resp, err := client.Inputtips(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "keywords参数不能为空")
}

// TestInputtips_APIError 测试Inputtips方法API返回错误
func TestInputtips_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建输入提示请求
	req := &inputtips.InputtipsRequest{
		Keywords: "望京",
		City:     "北京",
	}

	// 4. 执行输入提示请求
	resp, err := client.Inputtips(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// -------------------------- 天气信息测试 --------------------------

// TestWeatherinfo_Success 测试Weatherinfo方法正常请求成功
func TestWeatherinfo_Success(t *testing.T) {
	// 1. 创建mock服务器，返回天气信息成功响应
	mockServer := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"weatherinfo": {
			"city": "北京市",
			"cityid": "101010100",
			"temp": "22",
			"WD": "东南风",
			"WS": "1级",
			"SD": "40%",
			"AP": "1013hPa",
			"NJD": "10km",
			"WSE": "1",
			"time": "10:30",
			"isRadar": "1",
			"Radar": "JC_RADAR_AZ9010_JB",
			"weather": "晴",
			"temperature": "10~22℃",
			"winddirection": "东南",
			"windpower": "1-2级",
			"humidity": "40%"
		},
		"forecasts": [
			{
				"city": "北京市",
				"adcode": "110000",
				"province": "北京",
				"reporttime": "2023-05-20 10:30:00",
				"castype": "1",
				"forecast": [
					{
						"date": "2023-05-20",
						"week": "六",
						"dayweather": "晴",
						"nightweather": "晴",
						"daytemp": "22",
						"nighttemp": "10",
						"daywind": "东南风",
						"nightwind": "东南风",
						"daypower": "1级",
						"nightpower": "1级",
						"daytemp_float": 22.0,
						"nighttemp_float": 10.0
					}
				]
			}
		],
		"suggestion": {
			"comf": {
				"brf": "舒适",
				"txt": "白天温度适宜，风力不大，相信您在这样的天气条件下，应会感到比较清爽和舒适。",
				"type": "comf"
			},
			"cw": {
				"brf": "较适宜",
				"txt": "较适宜洗车，未来一天无雨，风力较小，擦洗一新的汽车至少能保持一天。",
				"type": "cw"
			}
		}
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建天气信息请求
	req := &weatherinfo.WeatherinfoRequest{
		City:       "北京市",
		Extensions: "all",
	}

	// 4. 执行天气信息请求
	resp, err := client.Weatherinfo(req)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "北京市", resp.Weatherinfo.City)
	assert.Equal(t, "101010100", resp.Weatherinfo.CityID)
	assert.Equal(t, "22", resp.Weatherinfo.Temp)
	assert.Equal(t, "晴", resp.Weatherinfo.Weather)
	assert.Equal(t, "10~22℃", resp.Weatherinfo.Temperature)
	assert.Len(t, resp.Forecasts, 1)
	assert.Len(t, resp.Forecasts[0].Forecast, 1)
	assert.Equal(t, "2023-05-20", resp.Forecasts[0].Forecast[0].Date)
	assert.Equal(t, "晴", resp.Forecasts[0].Forecast[0].Dayweather)
	assert.Equal(t, "舒适", resp.Suggestion.Comf.Brf)
	assert.Equal(t, "较适宜", resp.Suggestion.Cw.Brf)
}

// TestWeatherinfo_MissingCity 测试Weatherinfo方法缺少必填参数city
func TestWeatherinfo_MissingCity(t *testing.T) {
	// 1. 创建Client实例
	config := NewConfig("test_key")
	client, err := NewClient(config)
	require.NoError(t, err)

	// 2. 创建缺少城市的请求
	req := &weatherinfo.WeatherinfoRequest{
		Extensions: "base", // 缺少city
	}

	// 3. 执行天气信息请求
	resp, err := client.Weatherinfo(req)

	// 4. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, amapErr.InvalidConfigError(""), err)
	assert.Contains(t, err.Error(), "city参数不能为空")
}

// TestWeatherinfo_APIError 测试Weatherinfo方法API返回错误
func TestWeatherinfo_APIError(t *testing.T) {
	// 1. 创建mock服务器，返回API错误
	mockServer := mockResponse(http.StatusOK, `{
		"status": "0",
		"info": "无效的Key",
		"infocode": "10001"
	}`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 创建天气信息请求
	req := &weatherinfo.WeatherinfoRequest{
		City: "北京市",
	}

	// 4. 执行天气信息请求
	resp, err := client.Weatherinfo(req)

	// 5. 验证结果
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestHardwarePosition_Success 测试硬件定位API调用成功（v1）
func TestHardwarePosition_Success(t *testing.T) {
	// 1. 创建mock服务器，返回硬件定位API v1响应
	mockServer := mockResponse(http.StatusOK, `{
                "status": "1",
                "info": "OK",
                "infocode": "10000",
                "deviceid": "test_device",
                "latitude": 39.908722,
                "longitude": 116.397496,
                "accuracy": 5.0,
                "speed": 0.0,
                "direction": 0.0,
                "altitude": 43.5,
                "floor": 1,
                "timestamp": "2023-10-10T12:00:00Z",
                "location_type": "hybrid",
                "address": "北京市东城区东华门街道天安门广场",
                "poi": [{"name": "天安门", "distance": 100, "latitude": 39.9087, "longitude": 116.3975, "type": "风景名胜"}],
                "ad_info": {"province": "北京市", "city": "北京市", "district": "东城区", "adcode": "110101", "citycode": "010", "provincecode": "110000"}
        }`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行硬件定位API v1请求
	req := &positionV1.HardwarePositionRequest{
		DeviceID:    "test_device",
		GPS:         "39.908722,116.397496,0,0,1696900800,5",
		WiFi:        "00:11:22:33:44:55,-70,test_ssid|66:77:88:99:AA:BB,-60,test_ssid2",
		BaseStation: "460,0,12345,67890,-85",
		Output:      "JSON",
	}
	resp, err := client.HardwarePosition(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "test_device", resp.DeviceID)
	assert.Equal(t, 39.908722, resp.Latitude)
	assert.Equal(t, 116.397496, resp.Longitude)
	assert.Equal(t, 5.0, resp.Accuracy)
	assert.Equal(t, "hybrid", resp.LocationType)
	assert.Len(t, resp.POI, 1)
	assert.Equal(t, "天安门", resp.POI[0].Name)
}

// TestHardwarePositionV5_Success 测试硬件定位API调用成功（v5）
func TestHardwarePositionV5_Success(t *testing.T) {
	// 1. 创建mock服务器，返回硬件定位API v5响应
	mockServer := mockResponse(http.StatusOK, `{
                "status": "1",
                "info": "OK",
                "infocode": "10000",
                "deviceid": "test_device",
                "latitude": 39.908722,
                "longitude": 116.397496,
                "accuracy": 3.0,
                "speed": 0.0,
                "direction": 0.0,
                "altitude": 43.5,
                "floor": 1,
                "timestamp": "2023-10-10T12:00:00Z",
                "location_type": "hybrid",
                "address": "北京市东城区东华门街道天安门广场",
                "poi": [{"name": "天安门", "distance": 100, "latitude": 39.9087, "longitude": 116.3975, "type": "风景名胜", "address": "北京市东城区", "phone": "010-12345678"}],
                "ad_info": {"province": "北京市", "city": "北京市", "district": "东城区", "adcode": "110101", "citycode": "010", "provincecode": "110000", "township": "东华门街道", "village": "天安门社区"},
                "sensor_info": {"gps_valid": true, "wifi_valid": true, "basestation_valid": true, "bluetooth_valid": false, "used_sensor_types": ["gps", "wifi", "basestation"]},
                "confidence": 98.5,
                "indoor": false,
                "map_match": {"road_name": "长安街", "road_type": "主干道", "offset": 5.0, "direction": 90.0},
                "trace_id": "test_trace_id_123456"
        }`)
	defer mockServer.Close()

	// 2. 创建Client实例，使用mock服务器地址
	config := NewConfig("test_key")

	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行硬件定位API v5请求
	req := &positionV5.HardwarePositionRequest{
		DeviceID:      "test_device",
		GPS:           "39.908722,116.397496,0,0,1696900800,3",
		WiFi:          "00:11:22:33:44:55,-70,test_ssid,1|66:77:88:99:AA:BB,-60,test_ssid2,2",
		BaseStation:   "460,0,12345,67890,-85,2G",
		Barometer:     "1013.25",
		Accelerometer: "0,0,9.8",
		Gyroscope:     "0,0,0",
		Magnetometer:  "0,0,0",
		Orientation:   "0,0,0",
		PositionMode:  "1",
		Output:        "JSON",
	}
	resp, err := client.HardwarePositionV5(req)

	// 4. 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "test_device", resp.DeviceID)
	assert.Equal(t, 39.908722, resp.Latitude)
	assert.Equal(t, 116.397496, resp.Longitude)
	assert.Equal(t, 3.0, resp.Accuracy)
	assert.Equal(t, "hybrid", resp.LocationType)
	assert.Len(t, resp.POI, 1)
	assert.Equal(t, "天安门", resp.POI[0].Name)
	assert.Equal(t, "北京市东城区", resp.POI[0].Address)
	assert.Equal(t, "010-12345678", resp.POI[0].Phone)
	assert.True(t, resp.SensorInfo.GPSValid)
	assert.True(t, resp.SensorInfo.WiFiValid)
	assert.True(t, resp.SensorInfo.BaseStationValid)
	assert.Equal(t, 98.5, resp.Confidence)
	assert.False(t, resp.Indoor)
	assert.Equal(t, "长安街", resp.MapMatch.RoadName)
	assert.Equal(t, "test_trace_id_123456", resp.TraceID)
}

// TestBusStationID 测试公交站ID查询接口
func TestBusStationID(t *testing.T) {
	// 创建模拟服务器
	server := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"stationid": "123456",
		"name": "测试站点",
		"location": "116.397428,39.90923",
		"lines": [{
			"lineid": "line123",
			"name": "测试线路1",
			"start_time": "06:00",
			"end_time": "22:00",
			"distance": "10000",
			"stations": [{
				"id": "s1",
				"name": "站点1",
				"location": "116.397428,39.90923"
			}, {
				"id": "s2",
				"name": "站点2",
				"location": "116.398428,39.91023"
			}]
		}]}
`)
	defer server.Close()

	// 创建客户端
	client, _ := NewClient(&Config{
		Key: "test_key",
	})

	// 发送请求
	resp, err := client.BusStationID(&busStationID.StationIDRequest{
		ID:   "123456",
		City: "北京",
	})

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "123456", resp.StationID)
	assert.Equal(t, "测试站点", resp.Name)
	assert.Len(t, resp.Lines, 1)
}

// TestBusStationKeyword 测试公交站关键字查询接口
func TestBusStationKeyword(t *testing.T) {
	// 创建模拟服务器
	server := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"suggestion": {
			"keywords": ["测试站点"],
			"cities": ["北京"]
		},
		"stations": [{
			"id": "123456",
			"name": "测试站点",
			"location": "116.397428,39.90923",
			"cityid": "110000",
			"cityname": "北京",
			"address": "测试地址"
		}]}
`)
	defer server.Close()

	// 创建客户端
	client, _ := NewClient(&Config{
		Key: "test_key",
	})

	// 发送请求
	resp, err := client.BusStationKeyword(&busStationKeyword.StationKeywordRequest{
		Keywords: "测试站点",
		City:     "北京",
		Page:     "1",
		Offset:   "20",
	})

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "1", resp.Count)
	assert.Len(t, resp.Stations, 1)
}

// TestBusLineID 测试公交路线ID查询接口
func TestBusLineID(t *testing.T) {
	// 创建模拟服务器
	server := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"lineid": "line123",
		"name": "测试线路1",
		"type": "公交车",
		"start_time": "06:00",
		"end_time": "22:00",
		"distance": "10000",
		"polyline": "116.397428,39.90923;116.398428,39.91023",
		"stations": [{
			"id": "s1",
			"name": "站点1",
			"location": "116.397428,39.90923"
		}, {
			"id": "s2",
			"name": "站点2",
			"location": "116.398428,39.91023"
		}]}
`)
	defer server.Close()

	// 创建客户端
	client, _ := NewClient(&Config{
		Key: "test_key",
	})

	// 发送请求
	resp, err := client.BusLineID(&busLineID.LineIDRequest{
		ID:   "line123",
		City: "北京",
	})

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "line123", resp.LineID)
	assert.Equal(t, "测试线路1", resp.Name)
	assert.Len(t, resp.Stations, 2)
}

// TestBusLineKeyword 测试公交路线关键字查询接口
func TestBusLineKeyword(t *testing.T) {
	// 创建模拟服务器
	server := mockResponse(http.StatusOK, `{
		"status": "1",
		"info": "OK",
		"infocode": "10000",
		"count": "1",
		"suggestion": {
			"keywords": ["测试线路"],
			"cities": ["北京"]
		},
		"lines": [{
			"lineid": "line123",
			"name": "测试线路1",
			"type": "公交车",
			"start_time": "06:00",
			"end_time": "22:00",
			"distance": "10000",
			"from_stop": "站点1",
			"to_stop": "站点10"
		}]}
`)
	defer server.Close()

	// 创建客户端
	client, _ := NewClient(&Config{
		Key: "test_key",
	})

	// 发送请求
	resp, err := client.BusLineKeyword(&busLineKeyword.LineKeywordRequest{
		Keywords: "测试线路",
		City:     "北京",
		Page:     "1",
		Offset:   "20",
	})

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "1", resp.Status)
	assert.Equal(t, "OK", resp.Info)
	assert.Equal(t, "10000", resp.InfoCode)
	assert.Equal(t, "1", resp.Count)
	assert.Len(t, resp.Lines, 1)
}
