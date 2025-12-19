package utils

import (
	"math"
)

// 坐标系转换常量（火星坐标系偏移参数）
const (
	pi  = 3.1415926535897932384626
	a   = 6378245.0              // 长半轴
	ee  = 0.00669342162296594323 // 扁率
	xPi = pi * 3000.0 / 180.0
)

// Coordinate  坐标结构体（经度、纬度）
type Coordinate struct {
	Lng float64 // 经度
	Lat float64 // 纬度
}

// outOfChina 判断坐标是否在中国境内（境外坐标无需转换）
func outOfChina(c Coordinate) bool {
	if c.Lng < 72.004 || c.Lng > 137.8347 {
		return true
	}
	if c.Lat < 0.8293 || c.Lat > 55.8271 {
		return true
	}
	return false
}

// transformLat 纬度转换核心算法
func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*pi) + 40.0*math.Sin(y/3.0*pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*pi) + 320*math.Sin(y*pi/30.0)) * 2.0 / 3.0
	return ret
}

// transformLng 经度转换核心算法
func transformLng(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*pi) + 20.0*math.Sin(2.0*x*pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*pi) + 40.0*math.Sin(x/3.0*pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*pi) + 300.0*math.Sin(x/30.0*pi)) * 2.0 / 3.0
	return ret
}

// WGS84ToGCJ02 WGS84坐标系转GCJ02（火星坐标系/高德坐标系）
func WGS84ToGCJ02(wgs Coordinate) Coordinate {
	// 境外坐标直接返回
	if outOfChina(wgs) {
		return wgs
	}

	dLat := transformLat(wgs.Lng-105.0, wgs.Lat-35.0)
	dLng := transformLng(wgs.Lng-105.0, wgs.Lat-35.0)
	radLat := wgs.Lat / 180.0 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * pi)
	dLng = (dLng * 180.0) / (a / sqrtMagic * math.Cos(radLat) * pi)
	gcjLat := wgs.Lat + dLat
	gcjLng := wgs.Lng + dLng

	return Coordinate{Lng: gcjLng, Lat: gcjLat}
}

// GCJ02ToWGS84 GCJ02（火星坐标系/高德坐标系）转WGS84
// 采用迭代法提高精度（误差<1m）
func GCJ02ToWGS84(gcj Coordinate) Coordinate {
	// 境外坐标直接返回
	if outOfChina(gcj) {
		return gcj
	}

	// 初始迭代值
	wgs := Coordinate{
		Lng: gcj.Lng - 0.0065,
		Lat: gcj.Lat - 0.006,
	}
	// 迭代计算偏移量
	gcj2 := WGS84ToGCJ02(wgs)
	dLng := gcj2.Lng - gcj.Lng
	dLat := gcj2.Lat - gcj.Lat

	// 迭代3次足够满足精度
	for i := 0; i < 3; i++ {
		wgs.Lng -= dLng
		wgs.Lat -= dLat
		gcj2 = WGS84ToGCJ02(wgs)
		dLng = gcj2.Lng - gcj.Lng
		dLat = gcj2.Lat - gcj.Lat
	}

	return wgs
}
