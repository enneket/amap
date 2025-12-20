package main

import (
	"fmt"
	"log"
	"time"

	"github.com/enneket/amap"
	"github.com/enneket/amap/api/traffic_situation/rectangle"
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

	// 矩形区域内交通态势查询示例
	req := &rectangle.RectangleTrafficRequest{
		Rectangle: "116.351147,39.966309,116.357136,39.968722",
		Level:     "2",
	}

	resp, err := client.RectangleTrafficStatus(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== 矩形区域内交通态势查询结果 ===")
	fmt.Printf("状态: %s\n", resp.Status)
}
