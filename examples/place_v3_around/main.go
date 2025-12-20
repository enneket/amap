package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/place/v3/around"
)

func main() {
	// 创建配置
	config := &amap.Config{
		Key:     "your_amap_api_key",
		Timeout: 30 * time.Second,
	}

	// 初始化客户端
	client, err := amap.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// POI周边搜索V3示例
	req := &around.AroundSearchRequest{
		Location: "116.481028,39.989643",
		Radius:   1000,
		Offset:   10,
		Page:     1,
	}

	resp, err := client.PlaceV3Around(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== POI周边搜索V3结果 ===")
	fmt.Printf("总结果数: %d\n", resp.Count)
	for i, poi := range resp.Pois {
		fmt.Printf("POI %d:\n", i+1)
		fmt.Printf("名称: %s\n", poi.Name)
		fmt.Printf("地址: %s\n", poi.Address)
		fmt.Printf("经纬度: %s\n", poi.Location)
		fmt.Printf("电话: %s\n", poi.Tel)
		fmt.Printf("类别: %s\n", poi.Type)
		fmt.Printf("距离中心点: %d 米\n", poi.Distance)
		fmt.Println()
	}
}
