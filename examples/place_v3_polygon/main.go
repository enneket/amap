package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/place/v3/polygon"
)

func main() {
	// 创建配置
	config := amap.NewConfig("your_amap_api_key")
	config.Timeout = 30 * time.Second

	// 初始化客户端
	client, err := amap.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// POI多边形搜索V3示例
	req := &polygon.PolygonSearchRequest{
		Polygon:  "116.405467,39.907761;116.475098,39.907761;116.475098,39.940931;116.405467,39.940931;116.405467,39.907761",
		Offset:   10,
		Page:     1,
	}

	resp, err := client.PlaceV3Polygon(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== POI多边形搜索V3结果 ===")
	fmt.Printf("总结果数: %d\n", resp.Count)
	for i, poi := range resp.Pois {
		fmt.Printf("POI %d:\n", i+1)
		fmt.Printf("名称: %s\n", poi.Name)
		fmt.Printf("地址: %s\n", poi.Address)
		fmt.Printf("经纬度: %s\n", poi.Location)
		fmt.Printf("电话: %s\n", poi.Tel)
		fmt.Printf("类别: %s\n", poi.Type)
		fmt.Println()
	}
}
