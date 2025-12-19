package amap

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	bicycling "github.com/enneket/amap/api/direction/v1/bicycling"
	driving "github.com/enneket/amap/api/direction/v1/driving"
	walking "github.com/enneket/amap/api/direction/v1/walking"
	bicyclingV2 "github.com/enneket/amap/api/direction/v2/bicycling"
	busV2 "github.com/enneket/amap/api/direction/v2/bus"
	drivingV2 "github.com/enneket/amap/api/direction/v2/driving"
	electricV2 "github.com/enneket/amap/api/direction/v2/electric"
	walkingV2 "github.com/enneket/amap/api/direction/v2/walking"

	// v2版本API

	distance "github.com/enneket/amap/api/distance"
	geoCode "github.com/enneket/amap/api/geo_code"
	reGeoCode "github.com/enneket/amap/api/re_geo_code"
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
	config.BaseURL = mockServer.URL + "/"
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest("test/path", map[string]string{"param1": "value1"}, &resp)

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
	config.BaseURL = mockServer.URL + "/"
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest("test/path", nil, &resp)

	// 4. 验证结果
	assert.Error(t, err)
	assert.IsType(t, &amapErr.APIError{}, err)
	apiErr := err.(*amapErr.APIError)
	assert.Equal(t, "10001", apiErr.Code)
	assert.Equal(t, "无效的Key", apiErr.Info)
}

// TestDoRequest_ParseError 测试响应解析错误
func TestDoRequest_ParseError(t *testing.T) {
	// 1. 创建mock服务器，返回无效的JSON
	mockServer := mockResponse(http.StatusOK, "invalid json")
	defer mockServer.Close()

	// 2. 创建Client实例
	config := NewConfig("test_key")
	config.BaseURL = mockServer.URL + "/"
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest("test/path", nil, &resp)

	// 4. 验证结果
	assert.Error(t, err)
	assert.IsType(t, amapErr.ParseError(""), err)
}

// TestDoRequest_NetworkError 测试网络错误
func TestDoRequest_NetworkError(t *testing.T) {
	// 1. 创建Client实例，使用不存在的地址
	config := NewConfig("test_key")
	config.BaseURL = "http://non_existent_domain_12345.com/"
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
	config.BaseURL = mockServer.URL + "/"
	config.SecurityKey = "test_security_key"
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest("test/path", map[string]string{"param1": "value1"}, &resp)

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
	config.BaseURL = mockServer.URL + "/"
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求，传入部分参数
	var resp TestResponse
	customUserAgent := "test-agent/1.0"
	config.UserAgent = customUserAgent
	err = client.DoRequest("test/path", map[string]string{"param1": "value1", "param2": ""}, &resp)

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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
	client, err := NewClient(config)
	require.NoError(t, err)

	// 3. 执行请求
	var resp TestResponse
	err = client.DoRequest("test/path", nil, &resp)
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
	config.BaseURL = mockServer.URL + "/"
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
