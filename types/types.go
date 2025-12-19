package amap

import (
	"encoding/json"
	"io"

	amapErr "github.com/enneket/amap/errors"
)

// BaseResponse 高德API所有接口的通用基础响应
// 参考：https://lbs.amap.com/api/web-service/guide/common/errorcode
type BaseResponse struct {
	Status   string          `json:"status"`   // 响应状态：1=成功，0=失败
	Info     string          `json:"info"`     // 错误信息（失败时返回）
	InfoCode string          `json:"infocode"` // 错误码（失败时返回，如10001=Key无效）
	RawJSON  json.RawMessage `json:"-"`        // 原始响应JSON（用于后续解析业务数据）
}

// ReadBaseResponse 从HTTP响应体读取并解析BaseResponse
// 同时保留原始JSON数据，用于后续解析业务响应
func ReadBaseResponse(r io.Reader) (BaseResponse, []byte, error) {
	// 1. 先读取完整的响应体（避免body被消费后无法复用）
	rawBody, err := io.ReadAll(r)
	if err != nil {
		return BaseResponse{}, nil, amapErr.NewParseError("读取响应体失败: " + err.Error())
	}

	// 2. 解析基础响应
	var baseResp BaseResponse
	if err := json.Unmarshal(rawBody, &baseResp); err != nil {
		return BaseResponse{}, nil, amapErr.NewParseError("解析基础响应失败: " + err.Error())
	}

	// 3. 保存原始JSON（供业务响应解析）
	baseResp.RawJSON = rawBody

	return baseResp, rawBody, nil
}

// -------------------------- 其他通用类型 --------------------------

// CoordinateType 坐标系类型（高德API支持的坐标系）
type CoordinateType string

const (
	CoordinateType_GCJ02  CoordinateType = "gcj02"  // 高德火星坐标系（默认）
	CoordinateType_WGS84  CoordinateType = "wgs84"  // GPS坐标系（需申请权限）
	CoordinateType_BD09LL CoordinateType = "bd09ll" // 百度坐标系（仅部分接口支持）
)

// OutputType 响应格式类型
type OutputType string

const (
	OutputTypeJSON OutputType = "JSON" // 默认
	OutputTypeXML  OutputType = "XML"
)

// LanguageType 响应语言类型
type LanguageType string

const (
	LanguageTypeZH LanguageType = "zh_cn" // 中文（默认）
	LanguageTypeEN LanguageType = "en"    // 英文
)
