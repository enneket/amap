package re_geo_code

import "github.com/enneket/amap/types"

// ReGeocodeResponse 逆地理编码响应
// 文档：https://lbs.amap.com/api/webservice/guide/api/georegeo/#regeo
type ReGeocodeResponse struct {
	types.BaseResponse               // 继承基础响应（Status/Info/InfoCode）
	ReGeocode             ReGeocodeData `json:"regeocode"` // 逆地理编码核心数据
}

// ReGeocodeData 逆地理编码核心数据
type ReGeocodeData struct {
	FormattedAddress  string              `json:"formatted_address"` // 格式化地址（省+市+区+乡镇+街道+门牌号）
	AddressComponent  AddressComponent    `json:"addressComponent"`  // 地址组件详情
	Pois              []POIItem           `json:"pois"`              // POI列表（仅extensions=all时返回）
	Roads             []RoadItem          `json:"roads"`             // 道路列表（仅extensions=all时返回）
	RoadIntersections []RoadIntersectItem `json:"roadintersections"` // 道路交叉口列表（仅extensions=all时返回）
	Building          BuildingItem        `json:"building"`          // 建筑物信息（仅extensions=all时返回）
	Neighborhood      NeighborhoodItem    `json:"neighborhood"`      // 周边小区信息（仅extensions=all时返回）
	Township          string              `json:"township"`          // 乡镇名称
}

// AddressComponent 地址组件详情
type AddressComponent struct {
	Country      string `json:"country"`      // 国家（默认中国）
	Province     string `json:"province"`     // 省份（如"北京市"）
	City         string `json:"city"`         // 城市（如"北京市"）
	Citycode     string `json:"citycode"`     // 城市编码（如"110000"）
	District     string `json:"district"`     // 区县（如"朝阳区"）
	Adcode       string `json:"adcode"`       // 行政区划编码（如"110105"）
	Township     string `json:"township"`     // 乡镇（如"望京街道"）
	Towncode     string `json:"towncode"`     // 乡镇编码（如"110105028"）
	Street       string `json:"street"`       // 街道（如"望京街"）
	Number       string `json:"number"`       // 门牌号（如"8号"）
	BusinessArea string `json:"businessArea"` // 商圈（如"望京"）
	Floor        string `json:"floor"`        // 楼层（如"F1"，仅部分地址有）
	OfficeArea   string `json:"officeArea"`   // 办公区（如"SOHO塔1"，仅部分地址有）
	CountryCode  string `json:"countryCode"`  // 国家编码（如"CN"）
}

// POIItem POI信息（仅extensions=all时返回）
type POIItem struct {
	ID        string `json:"id"`        // POI唯一标识
	Name      string `json:"name"`      // POI名称（如"望京SOHO"）
	Type      string `json:"type"`      // POI类型（如"写字楼|商务办公"）
	Location  string `json:"location"`  // POI经纬度（"经度,纬度"）
	Address   string `json:"address"`   // POI地址
	Distance  string `json:"distance"`  // 与请求坐标的距离（单位：米）
	Direction string `json:"direction"` // 相对于请求坐标的方向（如"东北"）
	Tel       string `json:"tel"`       // 联系电话
	Website   string `json:"website"`   // 官网地址
	Email     string `json:"email"`     // 邮箱
	Pcode     string `json:"pcode"`     // 省份编码
	Pname     string `json:"pname"`     // 省份名称
	Citycode  string `json:"citycode"`  // 城市编码
	Cityname  string `json:"cityname"`  // 城市名称
	Adcode    string `json:"adcode"`    // 行政区划编码
	District  string `json:"district"`  // 区县名称
	Towncode  string `json:"towncode"`  // 乡镇编码
	Townname  string `json:"townname"`  // 乡镇名称
}

// RoadItem 道路信息（仅extensions=all时返回）
type RoadItem struct {
	Name      string `json:"name"`      // 道路名称（如"京密路"）
	Location  string `json:"location"`  // 道路中心点经纬度
	Distance  string `json:"distance"`  // 与请求坐标的距离（米）
	Direction string `json:"direction"` // 相对于请求坐标的方向
	Adcode    string `json:"adcode"`    // 行政区划编码
	Citycode  string `json:"citycode"`  // 城市编码
}

// RoadIntersectItem 道路交叉口信息（仅extensions=all时返回）
type RoadIntersectItem struct {
	Location  string   `json:"location"`  // 交叉口经纬度
	Distance  string   `json:"distance"`  // 与请求坐标的距离（米）
	Direction string   `json:"direction"` // 相对于请求坐标的方向
	Roads     []string `json:"roads"`     // 相交道路名称列表（如["京密路","望京街"]）
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
