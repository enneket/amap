package re_geo_code

import (
	amapType "github.com/enneket/amap/types"
)

// ReGeocodeResponse 逆地理编码响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo#t5
type ReGeocodeResponse struct {
	amapType.BaseResponse               // 继承基础响应（Status/Info/InfoCode）
	ReGeocode             ReGeocodeData `json:"regeocode"` // 逆地理编码核心数据
}

// ReGeocodeData 逆地理编码核心数据
type ReGeocodeData struct {
	FormattedAddress  string              `json:"formatted_address"` // 格式化地址（省+市+区+乡镇+街道+门牌号）
	AddressComponent  AddressComponent    `json:"addressComponent"`  // 地址元素列表
	Roads             []RoadItem          `json:"roads"`             // 道路列表（仅extensions=all时返回）
	RoadIntersections []RoadIntersectItem `json:"roadinters"`        // 道路交叉口列表（仅extensions=all时返回）
	Pois              []POIItem           `json:"pois"`              // POI列表（仅extensions=all时返回）
	Aois              []AoiItem           `json:"aois"`              // AOI列表（仅extensions=all时返回）
}

// AddressComponent 地址组件详情
type AddressComponent struct {
	Country      string           `json:"country"`      // 国家（默认中国）
	Province     string           `json:"province"`     // 省份（如"北京市"）
	City         []string         `json:"city"`         // 城市（如"北京市"）
	Citycode     string           `json:"citycode"`     // 城市编码（如"110000"）
	District     string           `json:"district"`     // 区县（如"朝阳区"）
	Adcode       string           `json:"adcode"`       // 行政区划编码（如"110105"）
	Township     string           `json:"township"`     // 乡镇（如"望京街道"）
	Towncode     string           `json:"towncode"`     // 乡镇编码（如"110105028"）
	Neighborhood NeighborhoodItem `json:"neighborhood"` // 周边小区信息（仅extensions=all时返回）
	Building     BuildingItem     `json:"building"`     // 建筑物信息（仅extensions=all时返回）
	StreetNumber StreetNumberItem `json:"streetNumber"` // 门牌信息列表
	SeaArea      string           `json:"seaArea"`      // 所属海域信息"渤海"）
	BusinessArea BusinessAreaItem `json:"businessArea"` // 商圈（如"望京"）
}

// POIItem POI信息（仅extensions=all时返回）
type POIItem struct {
	ID           string `json:"id"`           // POI唯一标识
	Name         string `json:"name"`         // POI名称（如"望京SOHO"）
	Type         string `json:"type"`         // POI类型（如"写字楼|商务办公"）
	Tel          any    `json:"tel"`          // 联系电话
	Distance     string `json:"distance"`     // 与请求坐标的距离（单位：米）
	Direction    string `json:"direction"`    // 相对于请求坐标的方向（如"东北"）
	Address      string `json:"address"`      // POI地址
	Location     string `json:"location"`     // POI经纬度（"经度,纬度"）
	BusinessArea string `json:"businessarea"` // POI所在商圈名称
}

// RoadItem 道路信息（仅extensions=all时返回）
type RoadItem struct {
	ID        string `json:"id"`        // 道路唯一标识
	Name      string `json:"name"`      // 道路名称（如"京密路"）
	Distance  string `json:"distance"`  // 与请求坐标的距离（米）
	Direction string `json:"direction"` // 方位
	Location  string `json:"location"`  // 坐标点
}

// RoadIntersectItem 道路交叉口信息（仅extensions=all时返回）
type RoadIntersectItem struct {
	Distance   string `json:"distance"`   // 交叉路口到请求坐标的距离
	Direction  string `json:"direction"`  // 方位
	Location   string `json:"location"`   // 路口经纬度
	FirstID    string `json:"firstID"`    // 第一条道路ID
	FirstName  string `json:"firstName"`  // 第一条道路名称
	SecondID   string `json:"secondID"`   // 第二条道路ID
	SecondName string `json:"secondName"` // 第二条道路名称
}

// BuildingItem 建筑物信息（仅extensions=all时返回）
type BuildingItem struct {
	Name     string `json:"name"`     // 建筑物名称（如"望京SOHO"）
	Type     string `json:"type"`     // 建筑物类型（如"写字楼"）
	Location string `json:"location"` // 建筑物中心点经纬度
}

// NeighborhoodItem 周边小区信息（仅extensions=all时返回）
type NeighborhoodItem struct {
	Name string `json:"name"` // 小区名称
	Type string `json:"type"` // 小区类型（如"住宅区"）
}

// StreetNumberItem 门牌信息
type StreetNumberItem struct {
	Street    string `json:"street"`    // 街道（如"望京街"）
	Number    string `json:"number"`    // 门牌号（如"8号"）
	Location  string `json:"location"`  // 坐标点
	Direction string `json:"direction"` // 方向
	Distance  string `json:"distance"`  // 门牌地址到请求坐标的距离
}

type BusinessAreaItem struct {
	BusinessArea string `json:"businessArea"` // 商圈（如"望京"）
	Location     string `json:"location"`     // 商圈中心点经纬度
	Name         string `json:"name"`         // 商圈名称（如"望京"）
	ID           string `json:"id"`           // 商圈所在区域的 adcode
}

type AoiItem struct {
	ID       string `json:"id"`       // AOI唯一标识
	Name     string `json:"name"`     // AOI名称（如"望京SOHO"）
	AdCode   string `json:"adcode"`   // AOI所在区域的 adcode
	Location string `json:"location"` // AOI经纬度（"经度,纬度"）
	Area     string `json:"area"`     // AOI面积（单位：平方米）
	Distance string `json:"distance"` // 与请求坐标的距离（单位：米）
	Type     string `json:"type"`     // AOI类型（如"写字楼|商务办公"）
}
