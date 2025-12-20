package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/etd/v4/driving"
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

	// 未来驾车路径规划V4示例
	req := &driving.ETDDrivingRequestV4{
		Origin:      "116.481028,39.989643",
		Destination: "116.514203,39.905409",
		DepartureTime: fmt.Sprintf("%d", time.Now().Add(time.Hour * 2).Unix()),
	}

	resp, err := client.ETDDrivingV4(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 未来驾车路径规划V4结果 ===")
	for i, path := range resp.Route.Paths {
		fmt.Printf("路径 %d:\n", i+1)
		fmt.Printf("距离: %d 米\n", path.Distance)
		fmt.Printf("耗时: %d 秒\n", path.Duration)
		fmt.Printf("路段数量: %d\n", len(path.Steps))
		fmt.Println()
	}
}
