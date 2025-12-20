package polygon

import (
	amapType "github.com/enneket/amap/types"
)

// PolygonSearchResponse POI搜索2.0多边形搜索响应
// 文档：https://lbs.amap.com/api/webservice/guide/api-advanced/newpoisearch
// 返回搜索结果，包括POI列表、分页信息等

type PolygonSearchResponse struct {
	amapType.BaseResponse               // 继承基础响应（Status/Info/InfoCode）
	Count                 string        `json:"count"`                   // 匹配的POI数量
	Pois                  []PoiItem     `json:"pois"`                    // POI列表
	Suggestion            *Suggestion   `json:"suggestion,omitempty"`    // 建议词列表（可选）
	UserLocation          *UserLocation `json:"user_location,omitempty"` // 用户位置信息（可选）
}

// PoiItem POI信息
// 包含POI的基本信息、地址、经纬度、类型等

type PoiItem struct {
	ID               string      `json:"id"`                         // POI唯一标识
	Name             string      `json:"name"`                       // POI名称
	Type             string      `json:"type"`                       // POI类型
	TypeCode         string      `json:"typecode"`                   // POI类型编码
	Address          string      `json:"address,omitempty"`          // 地址信息
	Location         string      `json:"location"`                   // POI坐标（经度,纬度）
	Tel              string      `json:"tel,omitempty"`              // 电话
	Postcode         string      `json:"postcode,omitempty"`         // 邮政编码
	Website          string      `json:"website,omitempty"`          // 网址
	Email            string      `json:"email,omitempty"`            // 邮箱
	Pcode            string      `json:"pcode,omitempty"`            // 省份编码
	Pname            string      `json:"pname,omitempty"`            // 省份名称
	Citycode         string      `json:"citycode,omitempty"`         // 城市编码
	Cityname         string      `json:"cityname,omitempty"`         // 城市名称
	Adcode           string      `json:"adcode,omitempty"`           // 行政区划编码
	Adname           string      `json:"adname,omitempty"`           // 行政区划名称
	BusinessArea     string      `json:"business_area,omitempty"`    // 商圈
	ShopID           string      `json:"shopid,omitempty"`           // 店铺ID
	ShopInfo         int         `json:"shopinfo,omitempty"`         // 是否有店铺信息（0/1）
	NaviPoiid        string      `json:"navipoiid,omitempty"`        // 导航POI ID
	EntranceLocation string      `json:"entrancelocation,omitempty"` // 入口坐标
	ExitLocation     string      `json:"exitlocation,omitempty"`     // 出口坐标
	Photos           []Photo     `json:"photos,omitempty"`           // 图片列表
	Children         []PoiItem   `json:"children,omitempty"`         // 子POI列表
	Rating           string      `json:"rating,omitempty"`           // 评分
	Cost             string      `json:"cost,omitempty"`             // 人均消费
	OpenTime         string      `json:"opentime,omitempty"`         // 营业时间
	Tags             string      `json:"tags,omitempty"`             // 标签
	IndoorMap        string      `json:"indoor_map,omitempty"`       // 是否有室内地图（0/1）
	IndoorData       *IndoorData `json:"indoor_data,omitempty"`      // 室内地图数据
	Distance         string      `json:"distance,omitempty"`         // 距离（仅周边搜索时返回）
	Direction        string      `json:"direction,omitempty"`        // 方向（仅周边搜索时返回）
	Floor            string      `json:"floor,omitempty"`            // 楼层
	ShopType         string      `json:"shop_type,omitempty"`        // 店铺类型
	GridCode         string      `json:"gridcode,omitempty"`         // 网格编码
	DistanceSort     string      `json:"distance_sort,omitempty"`    // 距离排序
	BizExt           *BizExt     `json:"biz_ext,omitempty"`          // 业务扩展信息
	Event            *Event      `json:"event,omitempty"`            // 活动信息
	Polyline         string      `json:"polyline,omitempty"`         // 边界坐标（仅AOI查询时返回）
}

// Photo POI图片信息
// 包含图片URL和标题

type Photo struct {
	Title string `json:"title"` // 图片标题
	URL   string `json:"url"`   // 图片URL
}

// IndoorData 室内地图数据
// 包含室内POI信息

type IndoorData struct {
	Floor     string    `json:"floor"`          // 楼层
	TrueFloor string    `json:"truefloor"`      // 真实楼层
	Cpid      string    `json:"cpid"`           // 建筑ID
	Pois      []PoiItem `json:"pois,omitempty"` // 室内POI列表
}

// BizExt 业务扩展信息
// 包含POI的业务相关信息

type BizExt struct {
	Cost        string `json:"cost,omitempty"`        // 人均消费
	Rating      string `json:"rating,omitempty"`      // 评分
	OpenTime    string `json:"opentime,omitempty"`    // 营业时间
	Charge      string `json:"charge,omitempty"`      // 是否收费（0/1）
	MCTags      string `json:"mctags,omitempty"`      // 商户标签
	SpecialTags string `json:"specialtags,omitempty"` // 特色标签
	FoodType    string `json:"foodtype,omitempty"`    // 餐饮类型
}

// Event 活动信息
// 包含POI相关的活动信息

type Event struct {
	StartTime string `json:"start_time,omitempty"` // 活动开始时间
	EndTime   string `json:"end_time,omitempty"`   // 活动结束时间
	Name      string `json:"name,omitempty"`       // 活动名称
	Type      string `json:"type,omitempty"`       // 活动类型
	Desc      string `json:"desc,omitempty"`       // 活动描述
}

// Suggestion 建议词列表
// 包含搜索建议和城市建议

type Suggestion struct {
	Keywords []string `json:"keywords"` // 搜索建议词列表
	Cities   []string `json:"cities"`   // 城市建议列表
}

// UserLocation 用户位置信息
// 包含用户的经纬度坐标

type UserLocation struct {
	Location string `json:"location"` // 用户坐标（经度,纬度）
}
