package amap

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

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
