package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/distance"
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

	// 距离测量示例
	req := &distance.DistanceRequest{
		Origins:      "116.481028,39.989643|116.455087,39.990464",
		Destination:  "116.514203,39.905409",
		Type:         0,
	}

	resp, err := client.Distance(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 距离测量结果 ===")
	for i, element := range resp.Results {
		fmt.Printf("结果 %d:\n", i+1)
		fmt.Printf("距离: %d 米\n", element.Distance)
		fmt.Println()
	}
}
