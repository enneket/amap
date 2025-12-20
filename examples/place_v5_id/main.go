package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/place/v5/id"
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

	// POI ID查询V5示例
	req := &id.IDRequest{
		ID: "B0FFH2L32K",
	}

	resp, err := client.PlaceV5ID(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== POI ID查询V5结果 ===")
	for i, poi := range resp.Pois {
		fmt.Printf("POI %d:\n", i+1)
		fmt.Printf("名称: %s\n", poi.Name)
		fmt.Printf("地址: %s\n", poi.Address)
		fmt.Printf("经纬度: %s\n", poi.Location)
		fmt.Printf("电话: %s\n", poi.Tel)
		fmt.Printf("类别: %s\n", poi.Type)
		fmt.Printf("评分: %s\n", poi.BizExt.Rating)
		fmt.Println()
	}
}
