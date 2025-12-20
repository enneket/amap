package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/traffic_situation/circle"
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

	// 圆形区域内交通态势查询示例
	req := &circle.CircleTrafficRequest{
		Center: "116.481028,39.989643",
		Radius: "5000",
		Level:  "2",
	}

	resp, err := client.CircleTrafficStatus(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 圆形区域内交通态势查询结果 ===")
	fmt.Printf("状态: %s\n", resp.Status)
}
