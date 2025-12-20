package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/district"
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

	// 行政区划查询示例
	req := &district.DistrictRequest{
		Keywords:    "北京",
		Subdistrict: "1",
		Extensions:  "base",
	}

	resp, err := client.District(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 行政区划查询结果 ===")
	for i, dist := range resp.Districts {
		fmt.Printf("结果 %d:\n", i+1)
		fmt.Printf("名称: %s\n", dist.Name)
		fmt.Printf("级别: %s\n", dist.Level)
		fmt.Printf("城市代码: %s\n", dist.Citycode)
		fmt.Printf("行政区划代码: %s\n", dist.Adcode)
		fmt.Printf("中心点: %s\n", dist.Center)
		fmt.Printf("子区域数量: %d\n", len(dist.Districts))
		if len(dist.Districts) > 0 {
			fmt.Println("子区域:")
			for j, subDist := range dist.Districts {
				fmt.Printf("  %d. %s (%s)\n", j+1, subDist.Name, subDist.Adcode)
			}
		}
		fmt.Println()
	}
}
