package utils

import (
	"math"
	"testing"
)

// 测试签名方法（可替换为自己的测试用例）
func TestSign(t *testing.T) {
	params := map[string]string{
		"key":       "testkey",
		"address":   "北京市朝阳区",
		"timestamp": "1716234567",
		"empty":     "", // 空值应被过滤
	}
	securityKey := "testsecuritykey"
	expected := "E8B75885687145819A9D88C99A8E9F87" // 示例值，需替换为实际计算结果

	sig := Sign(params, securityKey)
	if sig != expected {
		t.Errorf("签名错误，期望：%s，实际：%s", expected, sig)
	}
}

// 测试坐标转换
func TestCoordinateTransform(t *testing.T) {
	// 北京天安门WGS84坐标（示例）
	wgs := Coordinate{Lng: 116.39748, Lat: 39.908823}
	// 转换为GCJ02
	gcj := WGS84ToGCJ02(wgs)
	// 转换回WGS84
	wgs2 := GCJ02ToWGS84(gcj)

	// 误差应小于1米（经纬度误差约0.00001）
	if math.Abs(wgs.Lng-wgs2.Lng) > 0.00001 || math.Abs(wgs.Lat-wgs2.Lat) > 0.00001 {
		t.Errorf("坐标转换误差过大，原坐标：%+v，转换后：%+v", wgs, wgs2)
	}
}

// 测试参数编码
func TestEncodeParams(t *testing.T) {
	params := map[string]string{
		"address": "北京市朝阳区",
		"city":    "北京",
		"empty":   "",
	}
	encoded := EncodeParamsDefault(params)
	expected := "address=%E5%8C%97%E4%BA%AC%E5%B8%82%E6%9C%9D%E9%98%B3%E5%8C%BA&city=%E5%8C%97%E4%BA%AC"
	if encoded != expected {
		t.Errorf("参数编码错误，期望：%s，实际：%s", expected, encoded)
	}
}
