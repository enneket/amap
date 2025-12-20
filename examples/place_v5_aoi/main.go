package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/place/v5/aoi"
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

	// POI AOI查询V5示例
	req := &aoi.AOISearchRequest{
		ID:       "B000A837X5",
		Offset:   "10",
		Page:     "1",
	}

	resp, err := client.PlaceV5AOI(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== POI AOI查询V5结果 ===")
	fmt.Printf("总结果数: %d\n", resp.Count)
	for i, poi := range resp.Pois {
		fmt.Printf("POI %d:\n", i+1)
		fmt.Printf("名称: %s\n", poi.Name)
		fmt.Printf("地址: %s\n", poi.Address)
		fmt.Printf("经纬度: %s\n", poi.Location)
		fmt.Printf("类别: %s\n", poi.Type)
		fmt.Println()
	}
}
