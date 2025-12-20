package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/traffic_situation/line"
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

	// 指定线路交通态势查询示例
	req := &line.LineTrafficRequest{
		Path:   "116.481028,39.989643;116.514203,39.905409",
		Level:  "5",
	}

	resp, err := client.LineTrafficStatus(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 指定线路交通态势查询结果 ===")
	fmt.Printf("状态: %s\n", resp.Status)
}
